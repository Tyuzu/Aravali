package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// Product represents our data model
type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Global context and Redis client
var (
	ctx = context.Background()
	rdb *redis.Client
)

func main() {
	// 1. Initialize Redis Connection
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Change if using Docker/Cloud Redis
		Password: "",               // No password by default
		DB:       0,                // Default database
	})

	// Test the Redis connection on startup
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully.")

	// 2. Initialize the Modern ServeMux
	mux := http.NewServeMux()

	// 3. Register Routes using net/http Method & Wildcard matching
	mux.HandleFunc("GET /products{$}", listProducts)
	mux.HandleFunc("POST /products{$}", createProduct)
	mux.HandleFunc("GET /products/{id}", getProduct)
	mux.HandleFunc("PUT /products/{id}", updateProduct)
	mux.HandleFunc("DELETE /products/{id}", deleteProduct)

	// 4. Wrap the router in custom logging middleware
	loggedMux := loggerMiddleware(mux)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// ==========================================
// MIDDLEWARE
// ==========================================

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Serve the request
		next.ServeHTTP(w, r)

		// Log the request properties
		log.Printf("%s %s took %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// ==========================================
// CRUD HANDLERS (Redis Driven)
// ==========================================

// GET /products
func listProducts(w http.ResponseWriter, r *http.Request) {
	// Find all keys matching product:*
	keys, err := rdb.Keys(ctx, "product:*").Result()
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}

	list := []Product{}
	if len(keys) > 0 {
		// MGet fetches multiple string keys in a single roundtrip
		vals, err := rdb.MGet(ctx, keys...).Result()
		if err != nil {
			http.Error(w, "Failed to fetch product data", http.StatusInternalServerError)
			return
		}

		for _, val := range vals {
			if valStr, ok := val.(string); ok {
				var p Product
				_ = json.Unmarshal([]byte(valStr), &p)
				list = append(list, p)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// POST /products
func createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if p.ID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	// Serialize struct to JSON string for Redis
	data, _ := json.Marshal(p)
	redisKey := fmt.Sprintf("product:%s", p.ID)

	// Save to Redis (no expiration)
	err := rdb.Set(ctx, redisKey, data, 0).Err()
	if err != nil {
		http.Error(w, "Failed to save product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// GET /products/{id}
func getProduct(w http.ResponseWriter, r *http.Request) {
	// Pull {id} parameter using the native r.PathValue()
	id := r.PathValue("id")
	redisKey := fmt.Sprintf("product:%s", id)

	val, err := rdb.Get(ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(val))
}

// PUT /products/{id}
func updateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	redisKey := fmt.Sprintf("product:%s", id)

	// Verify it exists first
	_, err := rdb.Get(ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	p.ID = id

	data, _ := json.Marshal(p)
	if err := rdb.Set(ctx, redisKey, data, 0).Err(); err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// DELETE /products/{id}
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	redisKey := fmt.Sprintf("product:%s", id)

	// Del returns number of deleted keys. If 0, key didn't exist.
	deleted, err := rdb.Del(ctx, redisKey).Result()
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}
	if deleted == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
