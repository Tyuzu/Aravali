package main

import (
	"context"
	"fmt"
	"log"
	"nae/handlers"
	"nae/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// -----------------------------------------------------------------------------
// Endpoint Handlers
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// Main Entrypoint
// -----------------------------------------------------------------------------

func main() {
	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)

	// Protected Routes (Guarded by authMiddleware)
	mux.HandleFunc("/feed", middleware.AuthMiddleware(handlers.FeedHandler))
	mux.HandleFunc("/profile", middleware.AuthMiddleware(handlers.ProfileHandler))
	mux.HandleFunc("/analytics", middleware.AuthMiddleware(handlers.AnalyticsHandler))
	mux.HandleFunc("/settings", middleware.AuthMiddleware(handlers.SettingsHandler))

	port := ":8080"

	// Production-ready HTTP Server Configuration
	srv := &http.Server{
		Addr:         port,
		Handler:      middleware.CorsMiddleware(mux),
		ReadTimeout:  5 * time.Second,  // Prevents slow request attacks
		WriteTimeout: 10 * time.Second, // Prevents hanging responses
		IdleTimeout:  120 * time.Second,
	}

	// Server startup logic
	go func() {
		fmt.Println("==================================================")
		fmt.Printf("🚀 High-Performance Go API running on port %s\n", port)
		fmt.Println("==================================================")
		fmt.Printf("➜ Public Login:     http://localhost%s/login\n", port)
		fmt.Printf("➜ Public Register:  http://localhost%s/register\n", port)
		fmt.Printf("➜ View 1 Feed:      http://localhost%s/feed (Protected)\n", port)
		fmt.Printf("➜ View 2 Profile:   http://localhost%s/profile (Protected)\n", port)
		fmt.Printf("➜ View 3 Analytics: http://localhost%s/analytics (Protected)\n", port)
		fmt.Printf("➜ View 4 Settings:  http://localhost%s/settings (Protected)\n", port)
		fmt.Println("==================================================")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Listen error: %v", err)
		}
	}()

	// Graceful Shutdown Listening
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down API server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Forced to Shutdown: %v", err)
	}

	log.Println("Server gracefully stopped.")
}
