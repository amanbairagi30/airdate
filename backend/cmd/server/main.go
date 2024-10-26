package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	ID              int    `json:"id" db:"id"`
	Username        string `json:"username" db:"username"`
	Password        string `json:"-" db:"password"`
	TwitchUsername  string `json:"twitchUsername" db:"twitch_username"`
	DiscordUsername string `json:"discordUsername" db:"discord_username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type contextKey string

const userClaimsKey contextKey = "claims"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
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
	r.HandleFunc("/api/health", healthCheckHandler).Methods("GET")
	r.HandleFunc("/api/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/profile", authMiddleware(profileHandler)).Methods("GET")
	r.HandleFunc("/api/connect/twitch", authMiddleware(connectTwitchHandler)).Methods("POST")
	r.HandleFunc("/api/connect/discord", authMiddleware(connectDiscordHandler)).Methods("POST")
	r.HandleFunc("/api/games/search", authMiddleware(searchGamesHandler)).Methods("GET")

	if db == nil {
		log.Fatal("Database connection is nil before starting the server")
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		log.Println("Database connection is nil")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dbUser User
	err = db.Get(&dbUser, "SELECT id, username, password FROM users WHERE username = $1", user.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: dbUser.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		bearerPrefix := "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					http.Error(w, "Token is malformed", http.StatusUnauthorized)
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					http.Error(w, "Token is expired or not yet valid", http.StatusUnauthorized)
				} else {
					http.Error(w, "Couldn't handle this token: "+err.Error(), http.StatusUnauthorized)
				}
			} else {
				http.Error(w, "Couldn't handle this token: "+err.Error(), http.StatusUnauthorized)
			}
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	var user User
	err := db.Get(&user, "SELECT id, username, twitch_username, discord_username FROM users WHERE username = $1", claims.Username)
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
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
