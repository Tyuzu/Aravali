package middleware

import (
	"log"
	"nae/utils"
	"net/http"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------
// Middlewares
// -----------------------------------------------------------------------------

// CORS Middleware Wrapper
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("[%s] %s -> %s", time.Now().Format("15:04:05"), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Authentication Middleware: Validates "Authorization: Bearer <token>" header
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("[Auth Failure] Missing or invalid Authorization header on %s", r.URL.Path)
			utils.SendError(w, http.StatusUnauthorized, "Authentication required. Missing Bearer token.")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			utils.SendError(w, http.StatusUnauthorized, "Invalid auth token.")
			return
		}

		// Token verified — proceed to handler
		next.ServeHTTP(w, r)
	}
}
