package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// ==========================================
// CONFIGURATIONS, DATA MODELS, & CORE TYPES
// ==========================================

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Config struct {
	HTTPPort       string
	AllowedOrigins []string
	JWTSecret      []byte
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
}

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Claims struct {
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}

type contextKey string

const (
	userContextKey contextKey = "user"
	roleContextKey contextKey = "role"
)

// SecurityConfig matches the production strictness flags provided
var SecurityConfig = struct {
	MaxRequestBodySize int64
	RequestTimeout     time.Duration
	RateLimitRequests  int
	RateLimitWindow    time.Duration
}{
	MaxRequestBodySize: 10 << 20, // 10 MB limit
	RequestTimeout:     30 * time.Second,
	RateLimitRequests:  5, // 5 requests per window for login brute-force prevention
	RateLimitWindow:    1 * time.Minute,
}

// Global runtime context and Redis Client instance pointer
var (
	ctx = context.Background()
	rdb *redis.Client
)

// ==========================================
// MAIN RUNTIME ORCHESTRATION
// ==========================================

func main() {
	// Initialize Structured Logging Output Format
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := Config{
		HTTPPort:       ":8080",
		AllowedOrigins: []string{"http://localhost:3000", "https://yourfrontend.com"},
		JWTSecret:      []byte("your-highly-secure-production-secret-signing-key"),
	}

	// 1. Establish Redis Connectivity
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("Redis instance cluster handshake failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Successfully linked volatile store pipeline directly to Redis")

	// 2. Pre-seed Master Admin Account for Testing
	seedAdminAccount()

	// 3. Instantiate the modern net/http multiplexer multiplex route table layout
	mux := http.NewServeMux()

	// 4. Define Functional Middleware Chain Topology Pipelines
	globalPipeline := Chain(RecoveryMiddleware(), SecurityLimitsMiddleware(), CORSMiddleware(cfg.AllowedOrigins))
	protectedPipeline := Chain(globalPipeline, AuthJWTMiddleware(cfg.JWTSecret))
	adminPipeline := Chain(protectedPipeline, RequireRoles(RoleAdmin))

	// ==========================================
	// REGISTRATION ROUTING MAP ENTRIES
	// ==========================================

	// Auth Endpoints
	mux.Handle("POST /auth/register{$}", globalPipeline(http.HandlerFunc(registerHandler)))
	mux.Handle("POST /auth/login{$}", Chain(globalPipeline, LoginRateLimitMiddleware())(http.HandlerFunc(loginHandler)))
	mux.Handle("POST /auth/refresh{$}", globalPipeline(http.HandlerFunc(refreshHandler)))
	mux.Handle("POST /auth/logout{$}", globalPipeline(http.HandlerFunc(logoutHandler)))

	// User Protected Endpoints
	mux.Handle("GET /user/profile{$}", protectedPipeline(http.HandlerFunc(profileHandler)))

	// Admin / Shared Product CRUD Endpoints
	mux.Handle("GET /products{$}", protectedPipeline(http.HandlerFunc(listProductsHandler)))
	mux.Handle("POST /products{$}", adminPipeline(http.HandlerFunc(createProductHandler)))
	mux.Handle("GET /products/{id}", protectedPipeline(http.HandlerFunc(getProductHandler)))
	mux.Handle("PUT /products/{id}", adminPipeline(http.HandlerFunc(updateProductHandler)))
	mux.Handle("DELETE /products/{id}", adminPipeline(http.HandlerFunc(deleteProductHandler)))

	// 5. Construct Structured Network Socket Listener Interface Binding Parameters
	server := &http.Server{
		Addr:              cfg.HTTPPort,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	serverErrors := make(chan error, 1)

	// Execute Server Async Thread Loop Loop Execution Block
	go func() {
		slog.Info("Production API Engine Core listening actively", "port", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("server breakdown intercept: %w", err)
		}
	}()

	// 6. Graceful Interruption Lifecycle Control
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		slog.Error("Fatal server loop runtime crash event captured", "error", err)
		os.Exit(1)

	case sig := <-shutdownSignal:
		slog.Info("System boundary halt invocation received", "signal", sig.String())

		// Establish 20-second connection drainage grace-period fence boundaries
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server forced termination step required during drainage", "error", err)
			_ = server.Close()
		}
		slog.Info("All active network processes completely synchronized down. Goodbye.")
	}
}

// ==========================================
// MIDDLEWARE DECLARATIVE ENGINE IMPLEMENTATION
// ==========================================

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

func RecoveryMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					slog.Error("Fatal application panic intercepted securely",
						"error", err,
						"path", r.URL.Path,
						"stack", string(debug.Stack()),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func SecurityLimitsMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, SecurityConfig.MaxRequestBodySize)
			routeCtx, cancel := context.WithTimeout(r.Context(), SecurityConfig.RequestTimeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(routeCtx))
		})
	}
}

func CORSMiddleware(allowedOrigins []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			isAllowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					isAllowed = true
					break
				}
			}

			if isAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			w.Header().Set("Access-Control-Expose-Headers", "Authorization, X-Refresh-Intent")

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Idempotency-Key, X-Requested-With, X-Refresh-Intent, Accept, Origin")
				w.Header().Set("Access-Control-Max-Age", "300")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func LoginRateLimitMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			limiterKey := fmt.Sprintf("ratelimit:login:%s", ip)

			count, err := rdb.Incr(ctx, limiterKey).Result()
			if err != nil {
				http.Error(w, "Limiter internal failure", http.StatusInternalServerError)
				return
			}

			if count == 1 {
				rdb.Expire(ctx, limiterKey, SecurityConfig.RateLimitWindow)
			}

			if count > int64(SecurityConfig.RateLimitRequests) {
				w.Header().Set("Retry-After", strconv.Itoa(int(SecurityConfig.RateLimitWindow.Seconds())))
				http.Error(w, "Too many login attempts. Suspended briefly.", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AuthJWTMiddleware(secret []byte) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) < 8 || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized: Token configuration missing", http.StatusUnauthorized)
				return
			}

			tokenStr := authHeader[7:]
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized: Signature invalid or claim windows past", http.StatusUnauthorized)
				return
			}

			routeCtx := context.WithValue(r.Context(), userContextKey, claims.Username)
			routeCtx = context.WithValue(routeCtx, roleContextKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(routeCtx))
		})
	}
}

func RequireRoles(allowedRoles ...Role) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(roleContextKey).(Role)
			if !ok {
				http.Error(w, "Forbidden: Security policy context structural absence", http.StatusForbidden)
				return
			}
			match := false
			for _, r := range allowedRoles {
				if userRole == r {
					match = true
					break
				}
			}
			if !match {
				http.Error(w, "Forbidden: Action bounds profile violation exceptions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ==========================================
// AUTHENTICATION LOGIC FLOW IMPLEMENTATIONS
// ==========================================
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var input User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}
	if input.Username == "" || input.Password == "" {
		http.Error(w, "Missing credentials parameters fields", http.StatusBadRequest)
		return
	}
	if input.Role == "" {
		input.Role = RoleUser
	}
	userKey := fmt.Sprintf("user:%s", input.Username)
	exists, _ := rdb.Exists(ctx, userKey).Result()
	if exists > 0 {
		http.Error(w, "Identity token space registration collision", http.StatusConflict)
		return
	}
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedBytes)
	data, _ := json.Marshal(input)
	rdb.Set(ctx, userKey, data, 0)
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintln(w, "Registration verified successfully")
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct{ Username, Password string }
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Bad JSON format", http.StatusBadRequest)
		return
	}
	userKey := fmt.Sprintf("user:%s", credentials.Username)
	val, err := rdb.Get(ctx, userKey).Result()
	if errors.Is(err, redis.Nil) {
		http.Error(w, "Invalid authentication pairs mapping context", http.StatusUnauthorized)
		return
	}
	var storedUser User
	_ = json.Unmarshal([]byte(val), &storedUser)
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid authentication pairs mapping context", http.StatusUnauthorized)
		return
	}
	// Generate Access and Refresh JWT Tokens
	accessTokenClaims := &Claims{
		Username: storedUser.Username,
		Role:     storedUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessTokenStr, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString([]byte("your-highly-secure-production-secret-signing-key"))
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	opaqueRefreshToken := hex.EncodeToString(b)
	// Persist token whitelist record state map on Cache Engine
	refreshKey := fmt.Sprintf("refresh_token:%s", storedUser.Username)
	rdb.Set(ctx, refreshKey, opaqueRefreshToken, 7*24*time.Hour)
	// Emit Strict Cookie Structure Configuration Target Matrix
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    opaqueRefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Production note: switch true with global Tls endpoints forced
		Path:     "/auth/refresh",
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessTokenStr,
		"expires_in":   "900",
	})
}
func refreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Unauthorized: Context boundary token absence", http.StatusUnauthorized)
		return
	}
	var targetUsername string
	keys, err := rdb.Keys(ctx, "refresh_token:*").Result()
	if err != nil {
		http.Error(w, "Unauthorized: token registry unavailable", http.StatusInternalServerError)
		return
	}
	for _, key := range keys {
		storedToken, err := rdb.Get(ctx, key).Result()
		if err == nil && storedToken == cookie.Value {
			targetUsername = strings.TrimPrefix(key, "refresh_token:")
			break
		}
	}
	if targetUsername == "" {
		http.Error(w, "Unauthorized: Revoked or rotated token mapping frame anomalies", http.StatusUnauthorized)
		return
	}
	// Dynamic User Object Identity extraction
	userKey := fmt.Sprintf("user:%s", targetUsername)
	val, err := rdb.Get(ctx, userKey).Result()
	if err != nil {
		http.Error(w, "Identity reference space structural execution dropped", http.StatusUnauthorized)
		return
	}
	var u User
	_ = json.Unmarshal([]byte(val), &u)
	// MINT ROTATED ALTERNATIVES
	newAccessTokenClaims := &Claims{
		Username:         u.Username,
		Role:             u.Role,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute))},
	}
	newAccessTokenStr, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessTokenClaims).SignedString([]byte("your-highly-secure-production-secret-signing-key"))
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	newOpaqueRefreshToken := hex.EncodeToString(b)
	refreshKey := fmt.Sprintf("refresh_token:%s", u.Username)
	rdb.Set(ctx, refreshKey, newOpaqueRefreshToken, 7*24*time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newOpaqueRefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		Path:     "/auth/refresh",
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessTokenStr,
	})
}
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		keys, scanErr := rdb.Keys(ctx, "refresh_token:*").Result()
		if scanErr == nil {
			for _, key := range keys {
				storedToken, _ := rdb.Get(ctx, key).Result()
				if storedToken == cookie.Value {
					rdb.Del(ctx, key)
					break
				}
			}
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth/refresh",
		MaxAge:   -1,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, "Session destroyed successfully")
}
func profileHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userContextKey).(string)
	userKey := fmt.Sprintf("user:%s", username)
	val, err := rdb.Get(ctx, userKey).Result()
	if err != nil {
		http.Error(w, "User details missing", http.StatusNotFound)
		return
	}
	var u User
	_ = json.Unmarshal([]byte(val), &u)
	u.Password = "" // Security Boundary Sanitization Check Guard
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(u)
}

// ==========================================
// BUSINESS DOMAIN CORE PRODUCT CRUD HANDLERS
// ==========================================
func listProductsHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := rdb.Keys(ctx, "product:").Result()
	if err != nil {
		http.Error(w, "Failed data parsing states processing engine actions", http.StatusInternalServerError)
		return
	}
	list := make([]Product, 0)
	if len(keys) > 0 {
		vals, err := rdb.MGet(ctx, keys...).Result()
		if err != nil {
			http.Error(w, "Failed atomic transactional multi-reads", http.StatusInternalServerError)
			return
		}
		for _, chunk := range vals {
			if str, ok := chunk.(string); ok {
				var p Product
				_ = json.Unmarshal([]byte(str), &p)
				list = append(list, p)
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Malformed request context frame mapping payload", http.StatusBadRequest)
		return
	}
	if p.ID == "" {
		http.Error(w, "Required unique identifier index space absent", http.StatusBadRequest)
		return
	}
	key := fmt.Sprintf("product:%s", p.ID)
	data, _ := json.Marshal(p)
	if err := rdb.Set(ctx, key, data, 0).Err(); err != nil {
		http.Error(w, "Database failure processing storage states write ops", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(p)
}
func getProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // Modern native path parameter extractor feature
	key := fmt.Sprintf("product:%s", id)
	val, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		http.Error(w, "Target record resource item absent inside indexes", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(val))
}
func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	key := fmt.Sprintf("product:%s", id)
	_, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		http.Error(w, "Target mutation resource target array entity index absent", http.StatusNotFound)
		return
	}
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid data serialization schemas mapping payload input structural error", http.StatusBadRequest)
		return
	}
	p.ID = id
	data, _ := json.Marshal(p)
	_ = rdb.Set(ctx, key, data, 0)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(p)
}
func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	key := fmt.Sprintf("product:%s", id)
	count, err := rdb.Del(ctx, key).Result()
	if err != nil {
		http.Error(w, "Data processing transaction drops exceptions error", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Resource item map deletion coordinates space target absent", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ==========================================
// UTILITY HELPERS
// ==========================================
func seedAdminAccount() {
	adminKey := "user:admin_master"
	exists, _ := rdb.Exists(ctx, adminKey).Result()
	if exists == 0 {
		pass, _ := bcrypt.GenerateFromPassword([]byte("SecureAdminPass123!"), bcrypt.DefaultCost)
		masterAdmin := User{
			Username: "admin_master",
			Password: string(pass),
			Email:    "admin@system.internal",
			Role:     RoleAdmin,
		}
		data, _ := json.Marshal(masterAdmin)
		rdb.Set(ctx, adminKey, data, 0)
		slog.Info("Pre-seeded active structural administrative supervisor record profile credentials spaces successfully into instance cache pool matrices loops")
	}
}
