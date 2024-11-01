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
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"database/sql"
	"database/sql/driver"

	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

var db *sqlx.DB
var jwtKey []byte

type User struct {
	ID              int         `json:"id" db:"id"`
	Username        string      `json:"username" db:"username"`
	Password        string      `json:"-" db:"password"`
	TwitchUsername  *string     `json:"twitchUsername,omitempty" db:"twitch_username"`
	DiscordUsername *string     `json:"discordUsername,omitempty" db:"discord_username"`
	InstagramHandle *string     `json:"instagramHandle,omitempty" db:"instagram_handle"`
	YoutubeChannel  *string     `json:"youtubeChannel,omitempty" db:"youtube_channel"`
	FavoriteGames   *string     `json:"favoriteGames,omitempty" db:"favorite_games"`
	ConnectedGames  StringArray `json:"connectedGames" db:"connected_games"`
	IsPrivate       bool        `json:"isPrivate" db:"is_private"`
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

type UserProfile struct {
	User
	FollowersCount int  `json:"followersCount"`
	FollowingCount int  `json:"followingCount"`
	IsFollowing    bool `json:"isFollowing"`
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

	var err error
	db, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", os.Getenv("DB_NAME")))
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Fatalf("Error creating database: %v", err)
		}
	}

	db.Close()

	dbURL = fmt.Sprintf("%s dbname=%s", dbURL, os.Getenv("DB_NAME"))
	db, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

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
			connected_games TEXT[] DEFAULT '{}'::TEXT[],
			is_private BOOLEAN DEFAULT FALSE
		)
	`)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	_, err = db.Exec(`
		DO $$ 
		BEGIN 
			BEGIN
				ALTER TABLE users ADD COLUMN is_private BOOLEAN DEFAULT FALSE;
			EXCEPTION
				WHEN duplicate_column THEN 
					NULL;
			END;
		END $$;
	`)
	if err != nil {
		log.Printf("Error running migration: %v", err)
	}

	// Create followers table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS followers (
			follower_id INTEGER REFERENCES users(id),
			following_id INTEGER REFERENCES users(id),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (follower_id, following_id)
		)
	`)
	if err != nil {
		log.Fatal("Error creating followers table:", err)
	}

	r := mux.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		Debug:            true,
	})

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/register", registerHandler).Methods("POST")
	api.HandleFunc("/login", loginHandler).Methods("POST")
	api.HandleFunc("/users", getAllUsersHandler).Methods("GET")
	api.HandleFunc("/users/{username}", getUserProfileHandler).Methods("GET")
	api.HandleFunc("/health", healthCheckHandler).Methods("GET")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMiddleware)
	protected.HandleFunc("/profile", profileHandler).Methods("GET")
	protected.HandleFunc("/privacy", updatePrivacyHandler).Methods("POST")
	protected.HandleFunc("/connect/twitch", connectTwitchHandler).Methods("POST")
	protected.HandleFunc("/connect/discord", connectDiscordHandler).Methods("POST")
	protected.HandleFunc("/connect/instagram", connectInstagramHandler).Methods("POST")
	protected.HandleFunc("/connect/youtube", connectYoutubeHandler).Methods("POST")
	protected.HandleFunc("/connect/game", connectGameHandler).Methods("POST")
	protected.HandleFunc("/disconnect/instagram", disconnectInstagramHandler).Methods("POST")
	protected.HandleFunc("/disconnect/youtube", disconnectYoutubeHandler).Methods("POST")
	protected.HandleFunc("/disconnect/game", disconnectGameHandler).Methods("POST")
	protected.HandleFunc("/users/{username}/follow", followUserHandler).Methods("POST")
	protected.HandleFunc("/users/{username}/unfollow", unfollowUserHandler).Methods("POST")

	handler := corsHandler.Handler(r)

	log.Printf("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Register Handler Start ===")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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

	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)

	if username == "" || password == "" {
		http.Error(w, `{"error":"Username and password are required"}`, http.StatusBadRequest)
		return
	}

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

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
		"message":  "User created successfully",
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password mismatch: %v", err)
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

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
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Store claims in context
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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	query := strings.ToLower(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, `{"error":"Query parameter 'q' is required"}`, http.StatusBadRequest)
		return
	}

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

	var filteredGames []string
	for _, game := range games {
		if strings.Contains(strings.ToLower(game), query) {
			filteredGames = append(filteredGames, game)
		}
	}

	w.WriteHeader(http.StatusOK)
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
	var users []User
	err := db.Select(&users, `
		SELECT id, username, twitch_username, discord_username, 
			   instagram_handle, youtube_channel, favorite_games, 
			   connected_games, is_private 
		FROM users 
		WHERE is_private = false
		ORDER BY id DESC
	`)
	
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	// Get the requesting user's claims (if authenticated)
	var currentUsername string
	claims, ok := r.Context().Value(userClaimsKey).(*Claims)
	if ok && claims != nil {
		currentUsername = claims.Username
	}

	var profile User
	err := db.Get(&profile, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		}
		return
	}

	// Get follower and following counts
	var followersCount, followingCount int
	err = db.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM followers WHERE following_id = u.id) as followers_count,
			(SELECT COUNT(*) FROM followers WHERE follower_id = u.id) as following_count
		FROM users u WHERE u.username = $1
	`, username).Scan(&followersCount, &followingCount)
	if err != nil {
		log.Printf("Error getting follow counts: %v", err)
		// Continue without counts rather than failing
		followersCount = 0
		followingCount = 0
	}

	// Check if the current user is following this profile
	isFollowing := false
	if currentUsername != "" {
		err = db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM followers 
				WHERE follower_id = (SELECT id FROM users WHERE username = $1)
				AND following_id = (SELECT id FROM users WHERE username = $2)
			)
		`, currentUsername, username).Scan(&isFollowing)
		if err != nil {
			log.Printf("Error checking follow status: %v", err)
			// Continue without follow status rather than failing
		}
	}

	// Create response object
	response := struct {
		User
		FollowersCount int  `json:"followersCount"`
		FollowingCount int  `json:"followingCount"`
		IsFollowing    bool `json:"isFollowing"`
	}{
		User:           profile,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		IsFollowing:    isFollowing,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updatePrivacyHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	if claims == nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var requestBody struct {
		IsPrivate bool `json:"isPrivate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Updating privacy settings for user %s to %v", claims.Username, requestBody.IsPrivate)

	result, err := db.Exec(
		"UPDATE users SET is_private = $1 WHERE username = $2",
		requestBody.IsPrivate, claims.Username,
	)

	if err != nil {
		log.Printf("Error updating privacy settings: %v", err)
		http.Error(w, `{"error": "Error updating privacy settings"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, `{"error": "Error updating privacy settings"}`, http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"message":   "Privacy settings updated",
		"isPrivate": requestBody.IsPrivate,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
		return
	}
}

func disconnectInstagramHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)

	_, err := db.Exec("UPDATE users SET instagram_handle = NULL WHERE username = $1", claims.Username)
	if err != nil {
		log.Printf("Error disconnecting Instagram account: %v", err)
		http.Error(w, "Error disconnecting Instagram account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Instagram account disconnected successfully"})
}

func disconnectYoutubeHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)

	_, err := db.Exec("UPDATE users SET youtube_channel = NULL WHERE username = $1", claims.Username)
	if err != nil {
		log.Printf("Error disconnecting YouTube account: %v", err)
		http.Error(w, "Error disconnecting YouTube account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "YouTube channel disconnected successfully"})
}

func disconnectGameHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)

	var requestBody struct {
		GameName string `json:"gameName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//  connected games only show up in the user profile
	var user User
	err := db.Get(&user, "SELECT connected_games FROM users WHERE username = $1", claims.Username)
	if err != nil {
		http.Error(w, "Failed to get user games", http.StatusInternalServerError)
		return
	}

	var updatedGames StringArray
	for _, game := range user.ConnectedGames {
		if game != requestBody.GameName {
			updatedGames = append(updatedGames, game)
		}
	}

	_, err = db.Exec("UPDATE users SET connected_games = $1 WHERE username = $2",
		updatedGames, claims.Username)
	if err != nil {
		http.Error(w, "Failed to disconnect game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Game disconnected successfully",
	})
}

func followUserHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	vars := mux.Vars(r)
	targetUsername := vars["username"]

	// Get follower ID (current user)
	var followerID int
	err := db.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&followerID)
	if err != nil {
		http.Error(w, "Failed to get follower ID", http.StatusInternalServerError)
		return
	}

	// Get target user ID
	var followingID int
	err = db.QueryRow("SELECT id FROM users WHERE username = $1", targetUsername).Scan(&followingID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get following ID", http.StatusInternalServerError)
		}
		return
	}

	// Check if already following
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = $1 AND following_id = $2)",
		followerID, followingID).Scan(&exists)
	if err != nil {
		http.Error(w, "Failed to check follow status", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Already following this user", http.StatusBadRequest)
		return
	}

	// Create follow relationship
	_, err = db.Exec("INSERT INTO followers (follower_id, following_id) VALUES ($1, $2)",
		followerID, followingID)
	if err != nil {
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully followed user"})
}

func unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userClaimsKey).(*Claims)
	vars := mux.Vars(r)
	targetUsername := vars["username"]

	// Get IDs
	var followerID, followingID int
	err := db.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&followerID)
	if err != nil {
		http.Error(w, "Failed to get follower ID", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT id FROM users WHERE username = $1", targetUsername).Scan(&followingID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Delete follow relationship
	result, err := db.Exec("DELETE FROM followers WHERE follower_id = $1 AND following_id = $2",
		followerID, followingID)
	if err != nil {
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Not following this user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully unfollowed user"})
}
