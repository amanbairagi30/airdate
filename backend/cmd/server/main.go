package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
)

var db *sqlx.DB
var jwtKey []byte

type User struct {
	ID              int     `json:"id,omitempty" db:"id"`
	Username        string  `json:"username,omitempty" db:"username"`
	Password        string  `json:"-" db:"password"`  // Never send password in JSON
	TwitchUsername  *string `json:"twitchUsername,omitempty" db:"twitch_username"`  // Use pointer for nullable
	DiscordUsername *string `json:"discordUsername,omitempty" db:"discord_username"` // Use pointer for nullable
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type contextKey string

const userClaimsKey contextKey = "claims"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Error loading .env files:", err)
		}
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "pixel_and_chill"
	}

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword)

	var err error
	db, err = sqlx.Connect("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Fatalf("Error creating database: %v", err)
		}
		log.Printf("Note: %v", err)
	}

	dbURI = fmt.Sprintf("%s dbname=%s", dbURI, dbName)
	db, err = sqlx.Connect("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to the database")

	if db == nil {
		log.Fatal("Database connection is nil after successful connection")
	}

	jwtKey = []byte(os.Getenv("JWT_SECRET"))

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			twitch_username TEXT,
			discord_username TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},  
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})


	r.HandleFunc("/api/health", func(w http.ResponseWriter, _ *http.Request) {
		response := map[string]string{"status": "ok"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")
	r.HandleFunc("/api/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/login", loginHandler).Methods("POST")


	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(authMiddleware)
	protected.HandleFunc("/profile", profileHandler).Methods("GET")
	protected.HandleFunc("/connect/twitch", connectTwitchHandler).Methods("POST")
	protected.HandleFunc("/connect/discord", connectDiscordHandler).Methods("POST")
	protected.HandleFunc("/games/search", searchGamesHandler).Methods("GET")


	handler := c.Handler(r)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Register Handler Start ===")
	w.Header().Set("Content-Type", "application/json")
	

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, `{"error":"Error reading request body"}`, http.StatusBadRequest)
		return
	}
	

	log.Printf("Raw request body: %s", string(body))

	//  to parse JSON
	var req RegisterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("JSON parse error: %v", err)
		http.Error(w, `{"error":"Invalid JSON format"}`, http.StatusBadRequest)
		return
	}
	log.Printf("Username after parse: '%v'", req.Username)
	log.Printf("Password after parse: '%v'", req.Password)
	log.Printf("Username length: %d", len(req.Username))
	log.Printf("Password length: %d", len(req.Password))
	log.Printf("Username trimmed empty: %v", strings.TrimSpace(req.Username) == "")
	log.Printf("Password trimmed empty: %v", strings.TrimSpace(req.Password) == "")


	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		log.Printf("Validation failed - Username empty: %v, Password empty: %v", 
			strings.TrimSpace(req.Username) == "", 
			strings.TrimSpace(req.Password) == "")
		http.Error(w, `{"error":"Username and password are required"}`, http.StatusBadRequest)
		return
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	result, err := db.Exec(
		"INSERT INTO users (username, password) VALUES ($1, $2)",
		strings.TrimSpace(req.Username),
		string(hashedPassword),
	)
	if err != nil {
		log.Printf("Database error: %v", err)
		if strings.Contains(err.Error(), "unique constraint") {
			http.Error(w, `{"error":"Username already exists"}`, http.StatusConflict)
		} else {
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	rows, _ := result.RowsAffected()
	log.Printf("User created successfully. Rows affected: %d", rows)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
		"username": req.Username,
	})
	
	log.Printf("=== Register Handler End ===")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Login Handler Start ===")
	w.Header().Set("Content-Type", "application/json")


	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Printf("Error decoding login request: %v", err)
		http.Error(w, `{"error":"Invalid request format"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt for username: %s", loginRequest.Username)


	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE username = $1", loginRequest.Username)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}


	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		log.Printf("Password mismatch for user %s", loginRequest.Username)
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}


	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
		"username": user.Username,
	})
	log.Printf("=== Login Handler End ===")
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("=== Auth Middleware Start ===")
		
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("No Authorization header found")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}


		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}


		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			log.Printf("Error parsing token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Printf("Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("Successfully authenticated user: %s", claims.Username)
		log.Printf("=== Auth Middleware End ===")

		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Profile Handler Start ===")
	claims := r.Context().Value(userClaimsKey).(*Claims)
	
	log.Printf("Fetching profile for user: %s", claims.Username)
	
	var user User
	err := db.Get(&user, `
		SELECT id, username, twitch_username, discord_username 
			FROM users WHERE username = $1
	`, claims.Username)
	
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding user profile: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully fetched profile for user: %s", claims.Username)
	log.Printf("=== Profile Handler End ===")
}

func connectTwitchHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	var requestBody struct {
		TwitchUsername string `json:"twitchUsername"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET twitch_username = $1 WHERE username = $2", requestBody.TwitchUsername, claims.Username)
	if err != nil {
		log.Printf("Error connecting Twitch account: %v", err)
		http.Error(w, "Error connecting Twitch account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Twitch account connected successfully"})
}

func connectDiscordHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	var requestBody struct {
		DiscordUsername string `json:"discordUsername"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET discord_username = $1 WHERE username = $2", requestBody.DiscordUsername, claims.Username)
	if err != nil {
		log.Printf("Error connecting Discord account: %v", err)
		http.Error(w, "Error connecting Discord account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Discord account connected successfully"})
}

func searchGamesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	games := []string{"Game 1", "Game 2", "Game 3"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

// NEED TO WORK ON AUTH