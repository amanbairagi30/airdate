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

	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	"database/sql/driver"
)

var db *sqlx.DB
var jwtKey []byte

type User struct {
	ID              int      `json:"id" db:"id"`
	Username        string   `json:"username" db:"username"`
	Password        string   `json:"-" db:"password"`
	TwitchUsername  *string  `json:"twitchUsername,omitempty" db:"twitch_username"`
	DiscordUsername *string  `json:"discordUsername,omitempty" db:"discord_username"`
	InstagramHandle *string  `json:"instagramHandle,omitempty" db:"instagram_handle"`
	YoutubeChannel  *string  `json:"youtubeChannel,omitempty" db:"youtube_channel"`
	FavoriteGames   *string  `json:"favoriteGames,omitempty" db:"favorite_games"`
	ConnectedGames  StringArray `json:"connectedGames" db:"connected_games"`
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

type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		if string(v) == "{}" {
			*a = StringArray{}
			return nil
		}
		// Remove the curly braces and split
		trimmed := string(v)[1 : len(string(v))-1]
		if len(trimmed) > 0 {
			*a = strings.Split(trimmed, ",")
		} else {
			*a = StringArray{}
		}
		return nil
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type *StringArray", value)
	}
}

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	return fmt.Sprintf("{%s}", strings.Join(a, ",")), nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	jwtKey = []byte(os.Getenv("JWT_SECRET"))

	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	// First connect without database to create it if needed
	var err error
	db, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Create database if it doesn't exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", os.Getenv("DB_NAME")))
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Fatalf("Error creating database: %v", err)
		}
	}

	// Close initial connection
	db.Close()

	// Connect to the specific database
	dbURL = fmt.Sprintf("%s dbname=%s", dbURL, os.Getenv("DB_NAME"))
	db, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			twitch_username TEXT,
			discord_username TEXT,
			instagram_handle TEXT,
			youtube_channel TEXT,
			favorite_games TEXT,
			connected_games TEXT[] DEFAULT '{}'::TEXT[]
		)
	`)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		Debug:           true,
	})

	r := mux.NewRouter()

	// Configure CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		Debug:           true,
	})

	// Apply CORS middleware to all routes
	r.Use(corsMiddleware.Handler)

	// Public routes
	r.HandleFunc("/api/register", registerHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/login", loginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/users", getAllUsersHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/health", healthCheckHandler).Methods("GET", "OPTIONS")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(mux.MiddlewareFunc(authMiddleware))
	protected.HandleFunc("/profile", profileHandler).Methods("GET")
	protected.HandleFunc("/connect/twitch", connectTwitchHandler).Methods("POST")
	protected.HandleFunc("/connect/discord", connectDiscordHandler).Methods("POST")
	protected.HandleFunc("/connect/instagram", connectInstagramHandler).Methods("POST")
	protected.HandleFunc("/connect/youtube", connectYoutubeHandler).Methods("POST")
	protected.HandleFunc("/connect/game", connectGameHandler).Methods("POST")
	protected.HandleFunc("/disconnect/twitch", disconnectTwitchHandler).Methods("POST")
	protected.HandleFunc("/disconnect/discord", disconnectDiscordHandler).Methods("POST")

	handler := corsHandler.Handler(r)

	log.Printf("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Register Handler Start ===")
	
	// Set CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	
	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check database connection
	if err := db.Ping(); err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, `{"error":"Database connection error"}`, http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, `{"error":"Error reading request body"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Raw request body: %s", string(body))

	var req RegisterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("JSON parse error: %v", err)
		http.Error(w, `{"error":"Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)

	if username == "" || password == "" {
		http.Error(w, `{"error":"Username and password are required"}`, http.StatusBadRequest)
		return
	}

	// Check if username already exists
	var exists bool
	err = db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username)
	if err != nil {
		log.Printf("Error checking username existence: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, `{"error":"Username already exists"}`, http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Insert user
	result, err := db.Exec(`
		INSERT INTO users (username, password, connected_games) 
		VALUES ($1, $2, '{}'::text[])`,
		username, string(hashedPassword))
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	log.Printf("User created successfully. Rows affected: %d", rows)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
		"username": username,
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
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
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
	claims := r.Context().Value(userClaimsKey).(*Claims)
	log.Printf("=== Profile Handler Start ===")
	log.Printf("Fetching profile for user: %s", claims.Username)

	var user User
	err := db.Get(&user, `
		SELECT id, username, twitch_username, discord_username, 
			   instagram_handle, youtube_channel, favorite_games, connected_games
		FROM users WHERE username = $1`, claims.Username)
	
	if err != nil {
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

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Get All Users Handler Start ===")

	// Enable CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Check if database connection is alive
	if err := db.Ping(); err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, `{"error":"Database connection error"}`, http.StatusInternalServerError)
		return
	}

	var users []User
	err := db.Select(&users, `
		SELECT 
			id, 
			username, 
			twitch_username, 
			discord_username, 
			instagram_handle, 
			youtube_channel, 
			favorite_games, 
			COALESCE(connected_games, '{}'::text[]) as connected_games
		FROM users
		ORDER BY id DESC
		LIMIT 50
	`)
	
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Error fetching users: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// If no users found, return empty array instead of null
	if users == nil {
		users = []User{}
	}

	log.Printf("Found %d users", len(users))

	// Pretty print the first user for debugging
	if len(users) > 0 {
		userJSON, _ := json.MarshalIndent(users[0], "", "  ")
		log.Printf("First user data: %s", string(userJSON))
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("Error encoding users response: %v", err)
		http.Error(w, `{"error":"Error encoding response"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("=== Get All Users Handler End ===")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
