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
	InstagramHandle *string `json:"instagramHandle,omitempty" db:"instagram_handle"`
	YoutubeChannel *string `json:"youtubeChannel,omitempty" db:"youtube_channel"`
	FavoriteGames *string `json:"favoriteGames,omitempty" db:"favorite_games"`
	ConnectedGames []string `json:"connectedGames,omitempty" db:"connected_games"`
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

	// Add retry logic for database connection
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", dbURI)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second)
	}

	if db == nil {
		log.Fatal("Failed to connect to database after multiple attempts")
	}

	// Create database if it doesn't exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Fatalf("Error creating database: %v", err)
		}
		log.Printf("Note: %v", err)
	}

	// Connect to the specific database
	dbURI = fmt.Sprintf("%s dbname=%s", dbURI, dbName)
	db, err = sqlx.Connect("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	}

	// Drop and recreate tables to ensure correct schema
	_, err = db.Exec(`
		DROP TABLE IF EXISTS users;
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			twitch_username TEXT,
			discord_username TEXT,
			instagram_handle TEXT,
			youtube_channel TEXT,
			favorite_games TEXT,
			connected_games TEXT[] DEFAULT '{}'
		)
	`)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	log.Println("Successfully connected to the database")

	if db == nil {
		log.Fatal("Database connection is nil after successful connection")
	}

	jwtKey = []byte(os.Getenv("JWT_SECRET"))

	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/api/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/login", loginHandler).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.HandleFunc("/profile", authMiddleware(profileHandler)).Methods("GET")
	protected.HandleFunc("/connect/twitch", authMiddleware(connectTwitchHandler)).Methods("POST")
	protected.HandleFunc("/disconnect/twitch", authMiddleware(disconnectTwitchHandler)).Methods("POST")
	protected.HandleFunc("/connect/discord", authMiddleware(connectDiscordHandler)).Methods("POST")
	protected.HandleFunc("/disconnect/discord", authMiddleware(disconnectDiscordHandler)).Methods("POST")
	protected.HandleFunc("/connect/instagram", authMiddleware(connectInstagramHandler)).Methods("POST")
	protected.HandleFunc("/connect/youtube", authMiddleware(connectYoutubeHandler)).Methods("POST")
	protected.HandleFunc("/games/search", authMiddleware(searchGamesHandler)).Methods("GET")
	protected.HandleFunc("/connect/game", authMiddleware(connectGameHandler)).Methods("POST")

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Printf("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
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
		log.Printf("Error decoding request: %v", err)
		http.Error(w, `{"error":"Invalid request format"}`, http.StatusBadRequest)
		return
	}

	// Trim spaces from username and password
	username := strings.TrimSpace(loginRequest.Username)
	password := strings.TrimSpace(loginRequest.Password)

	log.Printf("Login attempt for username: %s", username)

	if username == "" || password == "" {
		log.Printf("Empty username or password")
		http.Error(w, `{"error":"Username and password are required"}`, http.StatusBadRequest)
		return
	}

	var user User
	err := db.Get(&user, `
		SELECT id, username, password
		FROM users 
		WHERE username = $1`,
		username)

	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password mismatch: %v", err)
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"token":    tokenString,
		"username": user.Username,
		"message":  "Login successful",
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		return
	}

	log.Printf("Login successful for user: %s", username)
	log.Printf("=== Login Handler End ===")
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	log.Printf("=== Profile Handler Start ===")
	log.Printf("Fetching profile for user: %s", claims.Username)

	var user User
	err := db.Get(&user, `
		SELECT id, username, 
			   COALESCE(twitch_username, '') as twitch_username,
			   COALESCE(discord_username, '') as discord_username,
			   COALESCE(instagram_handle, '') as instagram_handle,
			   COALESCE(youtube_channel, '') as youtube_channel,
			   COALESCE(favorite_games, '') as favorite_games
		FROM users WHERE username = $1`, claims.Username)
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		http.Error(w, "Error fetching user profile", http.StatusInternalServerError)
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

func connectInstagramHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	var requestBody struct {
		InstagramHandle string `json:"instagramHandle"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET instagram_handle = $1 WHERE username = $2", 
		requestBody.InstagramHandle, claims.Username)
	if err != nil {
		http.Error(w, "Error connecting Instagram account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Instagram account connected successfully"})
}

func connectYoutubeHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	var requestBody struct {
		YoutubeChannel string `json:"youtubeChannel"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET youtube_channel = $1 WHERE username = $2", 
		requestBody.YoutubeChannel, claims.Username)
	if err != nil {
		http.Error(w, "Error connecting YouTube account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "YouTube channel connected successfully"})
}

func searchGamesHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Popular games among Gen Z
	games := []string{
		"Fortnite",
		"Valorant",
		"Minecraft",
		"Among Us",
		"Roblox",
		"Call of Duty: Warzone",
		"Apex Legends",
		"Genshin Impact",
		"League of Legends",
		"Fall Guys",
		"Rocket League",
		"PUBG",
		"GTA V",
		"Overwatch 2",
		"FIFA 24",
		"NBA 2K24",
		"Lethal Company",
		"Palworld",
		"Helldivers 2",
		"Counter-Strike 2",
	}

	// Filter games based on search query
	var filteredGames []string
	for _, game := range games {
		if strings.Contains(strings.ToLower(game), query) {
			filteredGames = append(filteredGames, game)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredGames)
}

func disconnectTwitchHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	
	_, err := db.Exec("UPDATE users SET twitch_username = NULL WHERE username = $1", claims.Username)
	if err != nil {
		log.Printf("Error disconnecting Twitch account: %v", err)
		http.Error(w, "Error disconnecting Twitch account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Twitch account disconnected successfully"})
}

func disconnectDiscordHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	
	_, err := db.Exec("UPDATE users SET discord_username = NULL WHERE username = $1", claims.Username)
	if err != nil {
		log.Printf("Error disconnecting Discord account: %v", err)
		http.Error(w, "Error disconnecting Discord account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Discord account disconnected successfully"})
}

func connectGameHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	
	var requestBody struct {
		GameName string `json:"gameName"`
		GameID   string `json:"gameId"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the user's connected games
	_, err := db.Exec(`
		UPDATE users 
		SET connected_games = array_append(connected_games, $1) 
		WHERE username = $2 AND NOT $1 = ANY(connected_games)`,
		requestBody.GameName, claims.Username)
	
	if err != nil {
		log.Printf("Error connecting game: %v", err)
		http.Error(w, "Error connecting game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Successfully connected to %s", requestBody.GameName),
	})
}

// NEED TO WORK ON AUTH