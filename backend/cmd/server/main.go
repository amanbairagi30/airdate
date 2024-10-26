package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtKey []byte

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Initialize database connection
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	log.Println("Successfully connected to the database")

	// Set JWT key
	jwtKey = []byte(os.Getenv("JWT_SECRET"))

	// Create users table if not exists
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

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
	if err != nil {
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
	err = db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", user.Username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)
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
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(r.Header.Get("Authorization"), claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var user User
	err = db.QueryRow("SELECT id, username FROM users WHERE username = $1", claims.Username).Scan(&user.ID, &user.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func connectTwitchHandler(w http.ResponseWriter, r *http.Request) {
	// This is a simplified example. In a real-world scenario, you'd implement OAuth2 flow.
	var request struct {
		TwitchUsername string `json:"twitchUsername"`
	}
	json.NewDecoder(r.Body).Decode(&request)

	claims := r.Context().Value("claims").(*Claims)
	
	// Update user's Twitch username in the database
	_, err := db.Exec("UPDATE users SET twitch_username = $1 WHERE username = $2", request.TwitchUsername, claims.Username)
	if err != nil {
		http.Error(w, "Error updating Twitch username", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Twitch account connected successfully"})
}

func connectDiscordHandler(w http.ResponseWriter, r *http.Request) {
	// This is a simplified example. In a real-world scenario, you'd implement OAuth2 flow.
	var request struct {
		DiscordUsername string `json:"discordUsername"`
	}
	json.NewDecoder(r.Body).Decode(&request)

	claims := r.Context().Value("claims").(*Claims)
	
	// Update user's Discord username in the database
	_, err := db.Exec("UPDATE users SET discord_username = $1 WHERE username = $2", request.DiscordUsername, claims.Username)
	if err != nil {
		http.Error(w, "Error updating Discord username", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Discord account connected successfully"})
}

func searchGamesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	// This is a mock implementation. In a real-world scenario, you'd use the IGDB API.
	games := []map[string]string{
		{"id": "1", "name": "The Legend of Zelda: Breath of the Wild"},
		{"id": "2", "name": "Red Dead Redemption 2"},
		{"id": "3", "name": "God of War"},
	}

	json.NewEncoder(w).Encode(games)
}
