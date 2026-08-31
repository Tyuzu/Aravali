package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing JWTs
var jwtSecret = []byte("KhovaSuperSecretKey2026")

// Response Structs
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type AuthRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthData struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

// Custom JWT Claims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func generateToken(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func sendJSON(w http.ResponseWriter, status int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendJSON(w, http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "Authorization header missing",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendJSON(w, http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "Invalid authorization header format",
			})
			return
		}

		tokenStr := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			sendJSON(w, http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			return
		}

		next.ServeHTTP(w, r)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "Invalid request payload"})
		return
	}

	if req.Email == "" || req.Password == "" {
		sendJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "Email and password required"})
		return
	}

	token, err := generateToken(req.Email)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: "Token generation failed"})
		return
	}

	sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: AuthData{
			Token: token,
			User:  req.Email,
		},
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "Invalid request payload"})
		return
	}

	if req.Email == "" || req.Password == "" {
		sendJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "Missing email or password"})
		return
	}

	token, err := generateToken(req.Email)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: "Token generation failed"})
		return
	}

	sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: AuthData{
			Token: token,
			User:  req.Email,
		},
	})
}

func handleFeed(w http.ResponseWriter, r *http.Request) {
	items := []map[string]interface{}{
		{"id": 1, "title": "Welcome to Solar2D + Go", "body": "Stack setup complete with JWT authentication."},
		{"id": 2, "title": "Caching Active", "body": "Network client returns cached memory data on instant revisits."},
	}
	sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	profile := map[string]interface{}{
		"username": "DevUser",
		"role":     "Mobile Engineer",
		"joined":   "2026-01-15",
	}
	sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: profile})
}

func handleAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics := map[string]interface{}{
		"activeUsers": 1280,
		"requests":    45920,
		"uptime":      "99.98%",
	}
	sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: analytics})
}

func handleSettings(w http.ResponseWriter, r *http.Request) {
	settings := map[string]interface{}{
		"apiVersion":           "v1.4.0",
		"theme":                "dark",
		"notificationsEnabled": true,
		"darkMode":             true,
	}
	sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: settings})
}

func main() {
	mux := http.NewServeMux()

	// Auth Routes (Public)
	mux.HandleFunc("/login", handleLogin)
	mux.HandleFunc("/register", handleRegister)

	// App Content Routes (Protected by JWT Middleware)
	mux.HandleFunc("/feed", jwtMiddleware(handleFeed))
	mux.HandleFunc("/profile", jwtMiddleware(handleProfile))
	mux.HandleFunc("/analytics", jwtMiddleware(handleAnalytics))
	mux.HandleFunc("/settings", jwtMiddleware(handleSettings))

	port := ":8080"
	fmt.Printf("Server listening on http://localhost%s...\n", port)
	if err := http.ListenAndServe(port, corsMiddleware(mux)); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
