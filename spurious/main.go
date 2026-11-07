package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// Configurations
var jwtSecret = []byte("your-super-secure-production-secret-key")

type User struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type contextKey string

const userContextKey contextKey = "user"

var (
	ctx = context.Background()
	rdb *redis.Client
)

func main() {
	rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	mux := http.NewServeMux()

	// Public Auth Routes (Login wrapped with structural Rate Limiter)
	mux.HandleFunc("POST /auth/register{$}", registerHandler)
	mux.Handle("POST /auth/login{$}", rateLimitMiddleware(http.HandlerFunc(loginHandler)))
	mux.HandleFunc("POST /auth/refresh{$}", refreshHandler)
	mux.HandleFunc("POST /auth/logout{$}", logoutHandler)

	// Protected Application Profile Route
	mux.Handle("GET /user/profile{$}", authJWTMiddleware(http.HandlerFunc(profileHandler)))

	log.Println("JWT Auth, CSRF & Rate-Limited Server operational on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// ==========================================
// MIDDLEWARES (Rate Limiting & JWT Verification)
// ==========================================

// Fixed-window rate limiter preventing brute forcing on login (5 attempts per minute per IP)
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // Production note: Pull from "X-Forwarded-For" if behind reverse proxies
		limiterKey := fmt.Sprintf("ratelimit:login:%s", ip)

		// Increment request counter inside a 1-minute window block
		count, err := rdb.Incr(ctx, limiterKey).Result()
		if err != nil {
			http.Error(w, "Limiter internal failure", http.StatusInternalServerError)
			return
		}

		if count == 1 {
			rdb.Expire(ctx, limiterKey, 1*time.Minute)
		}

		if count > 5 {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Too many login attempts. Please try again in a minute.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Verifies Bearer JWT tokens passed in the standard Authorization Header
func authJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			http.Error(w, "Unauthorized: Missing Bearer Token", http.StatusUnauthorized)
			return
		}

		tokenStr := authHeader[7:]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid or expired access token", http.StatusUnauthorized)
			return
		}

		// Inject verified subject into native request runtime context
		ctxWithUser := context.WithValue(r.Context(), userContextKey, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}

// ==========================================
// CORE HANDLERS
// ==========================================

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var input User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	userKey := fmt.Sprintf("user:%s", input.Username)
	exists, _ := rdb.Exists(ctx, userKey).Result()
	if exists > 0 {
		http.Error(w, "Username taken", http.StatusConflict)
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashed)

	data, _ := json.Marshal(input)
	rdb.Set(ctx, userKey, data, 0)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Registration complete")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct{ Username, Password string }
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}

	userKey := fmt.Sprintf("user:%s", creds.Username)
	val, err := rdb.Get(ctx, userKey).Result()
	if errors.Is(err, redis.Nil) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	var storedUser User
	_ = json.Unmarshal([]byte(val), &storedUser)

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 1. Generate short-lived Access Token (15 Mins)
	accessTokenClaims := &Claims{
		Username: storedUser.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, _ := accessToken.SignedString(jwtSecret)

	// 2. Generate long-lived Refresh Token (7 Days)
	refreshTokenClaims := &Claims{
		Username: storedUser.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, _ := refreshToken.SignedString(jwtSecret)

	// 3. Track active Refresh Token inside Redis for stateful invalidation/rotation
	refreshKey := fmt.Sprintf("refresh_token:%s", storedUser.Username)
	rdb.Set(ctx, refreshKey, refreshTokenString, 7*24*time.Hour)

	// 4. Inject Refresh Token into an anti-CSRF, HttpOnly SameSite Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,                    // Blocks XSS execution extraction
		Secure:   false,                   // Set TRUE in production for mandatory HTTPS enforced links
		Path:     "/auth/refresh",         // Scopes cookie access strictly to token rejuvenation endpoint
		SameSite: http.SameSiteStrictMode, // Strict CSRF counter-mechanism policy validation
	})

	// 5. Send short-lived Access Token safely via raw JSON response body
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessTokenString,
		"expires_in":   "900",
	})
}

// Token Rotation & Re-issuance Strategy Endpoint
func refreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Missing refresh token cookie", http.StatusUnauthorized)
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token payload", http.StatusUnauthorized)
		return
	}

	// Double-check whitelist array registry in Redis to verify token isn't revoked
	refreshKey := fmt.Sprintf("refresh_token:%s", claims.Username)
	storedToken, err := rdb.Get(ctx, refreshKey).Result()
	if errors.Is(err, redis.Nil) || storedToken != cookie.Value {
		http.Error(w, "Token revoked or rotated out asynchronously", http.StatusUnauthorized)
		return
	}

	// ROTATION MECHANIC: Mint a brand new pair to neutralize reuse risk profiles
	newAccessTokenClaims := &Claims{
		Username:         claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute))},
	}
	newAccessToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessTokenClaims).SignedString(jwtSecret)

	newRefreshTokenClaims := &Claims{
		Username:         claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))},
	}
	newRefreshToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, newRefreshTokenClaims).SignedString(jwtSecret)

	// Commit updated whitelisted token state inside Redis storage context
	rdb.Set(ctx, refreshKey, newRefreshToken, 7*24*time.Hour)

	// Re-assign Cookie reference with Strict parameters
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		Path:     "/auth/refresh",
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		claims := &Claims{}
		_, _ = jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		// Evict whitelisted tokens explicitly from active cache clusters immediately
		refreshKey := fmt.Sprintf("refresh_token:%s", claims.Username)
		rdb.Del(ctx, refreshKey)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth/refresh",
		MaxAge:   -1,
		HttpOnly: true,
	})
	fmt.Fprintln(w, "Logged out cleanly")
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userContextKey).(string)
	userKey := fmt.Sprintf("user:%s", username)
	val, err := rdb.Get(ctx, userKey).Result()
	if err != nil {
		http.Error(w, "Profile missing", http.StatusNotFound)
		return
	}

	var u User
	_ = json.Unmarshal([]byte(val), &u)
	u.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}
