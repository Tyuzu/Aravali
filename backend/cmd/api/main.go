// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"scav/config"
	"scav/infra"
	mq "scav/infra/mq"
	"scav/infra/workers"
	"scav/internal/mechat"
	"scav/internal/newchat"
	"scav/middleware"
	"scav/routes"
	"scav/subscribers"
	"scav/utils/logger"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	// =====================
	// Logger
	// =====================
	if err := logger.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		_ = logger.Sync()
	}()

	// =====================
	// Configuration
	// =====================
	cfg, err := config.InitConfig()
	if err != nil {
		logger.L.Sugar().Fatalw(
			"config validation failed",
			"error", err,
		)
	}

	// =====================
	// Infrastructure
	// =====================
	app, err := infra.New(cfg)
	if err != nil {
		logger.L.Sugar().Fatalw(
			"failed to initialize infrastructure",
			"error", err,
		)
	}

	// Ensure infrastructure is eventually closed even if main
	// returns unexpectedly after initialization.
	defer func() {
		closeCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := app.Close(closeCtx); err != nil {
			logger.L.Sugar().Errorw(
				"failed to close infrastructure",
				"error", err,
			)
		}
	}()

	// =====================
	// Application Lifecycle
	// =====================
	//
	// This context is shared by background workers and MQ subscribers.
	// When appCancel() is called during shutdown, subscribers and workers stop.
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	// Media worker subscription.
	var mediaSub mq.Subscription

	// =====================
	// MQ Subscribers
	// =====================
	if app.MQ != nil {
		if err := subscribers.RegisterAll(appCtx, app); err != nil {
			logger.L.Sugar().Fatalw(
				"failed to register MQ subscribers",
				"error", err,
			)
		}

		logger.L.Sugar().Infow(
			"MQ subscribers registered",
			"mq", "redis_pubsub",
		)

		// Start media worker consumer so background workers process media.jobs.
		if sub, err := workers.StartMediaWorker(
			appCtx,
			app.MQ,
		); err != nil {
			logger.L.Sugar().Errorw(
				"failed to start media worker",
				"error", err,
			)
		} else {
			mediaSub = sub

			logger.L.Sugar().Infow(
				"media worker started",
			)
		}
	} else {
		logger.L.Sugar().Warnw(
			"MQ is not configured; skipping MQ subscribers",
		)
	}

	// =====================
	// Rate Limiter
	// =====================
	rateLimiter := middleware.NewRateLimiter(
		1,
		12,
		10*time.Minute,
		10000,
	)

	// =====================
	// Chat Hubs
	// =====================
	hub := newchat.NewHub()
	go hub.Run()

	mehub := mechat.NewHub()
	go mehub.Run()

	// =====================
	// Router Setup
	// =====================
	router := routes.SetupRouter(
		app,
		rateLimiter,
	)

	newchat.AddNewChatRoutes(
		router,
		hub,
		app,
		rateLimiter,
	)

	mechat.AddMeChatRoutes(
		router,
		mehub,
		app,
		rateLimiter,
	)

	routes.AddStaticRoutes(router)

	// =====================
	// Readiness Probe
	// =====================
	router.GET(
		"/ready",
		func(
			w http.ResponseWriter,
			r *http.Request,
			_ httprouter.Params,
		) {
			ctx, cancel := context.WithTimeout(
				r.Context(),
				2*time.Second,
			)
			defer cancel()

			// Relational / SQL DB Check
			if app.SQLDB != nil {
				if err := app.SQLDB.Ping(ctx); err != nil {
					http.Error(
						w,
						"sqldb_unavailable",
						http.StatusServiceUnavailable,
					)
					return
				}
			} else if app.DB != nil {
				if err := app.DB.Ping(ctx); err != nil {
					http.Error(
						w,
						"db_unavailable",
						http.StatusServiceUnavailable,
					)
					return
				}
			}

			// Cache Check
			if app.Cache != nil {
				if _, err := app.Cache.Ping(ctx); err != nil {
					http.Error(
						w,
						"cache_unavailable",
						http.StatusServiceUnavailable,
					)
					return
				}
			}

			// Message Queue Check
			if app.MQ != nil {
				if err := app.MQ.Ping(ctx); err != nil {
					http.Error(
						w,
						"mq_unavailable",
						http.StatusServiceUnavailable,
					)
					return
				}
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		},
	)

	// =====================
	// Middleware & CORS
	// =====================
	handler := middleware.LoggingMiddleware(
		middleware.SecurityHeaders(router),
	)

	corsOpts := cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: []string{
			"HEAD",
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"Idempotency-Key",
			"X-Requested-With",
			"Accept",
			"Origin",
		},
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           300,
	}

	corsHandler := cors.New(corsOpts).Handler(handler)

	// =====================
	// HTTP Server
	// =====================
	server := &http.Server{
		Addr:              cfg.HTTPPort,
		Handler:           corsHandler,
		ReadTimeout:       7 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// =====================
	// Start HTTP Server
	// =====================
	go func() {
		logger.L.Sugar().Infow(
			"API server listening",
			"addr", cfg.HTTPPort,
			"protocol", "http",
		)

		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			logger.L.Sugar().Fatalw(
				"HTTP server error",
				"error", err,
			)
		}
	}()

	// =====================
	// Wait for Shutdown Signal
	// =====================
	sigCh := make(chan os.Signal, 1)

	signal.Notify(
		sigCh,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-sigCh

	logger.L.Sugar().Infow(
		"Shutting down server...",
	)

	// =====================
	// Graceful Shutdown
	// =====================
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer shutdownCancel()

	// Stop accepting new HTTP requests first.
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.L.Sugar().Errorw(
			"HTTP server shutdown error",
			"error", err,
		)
	}

	// =====================
	// Stop MQ Subscribers & Background Workers
	// =====================
	logger.L.Sugar().Infow(
		"Stopping MQ subscribers and background workers...",
	)

	appCancel()

	// Unsubscribe media worker explicitly if active.
	if mediaSub != nil {
		if err := mediaSub.Unsubscribe(); err != nil {
			logger.L.Sugar().Errorw(
				"media worker unsubscribe failed",
				"error", err,
			)
		}
	}

	// Stop application-level background components.
	rateLimiter.Stop()
	hub.Stop()
	mehub.Stop()

	// =====================
	// Close Infrastructure
	// =====================
	//
	// This closes:
	// - PostgreSQL
	// - Redis (including Redis Pub/Sub resources)
	// - MongoDB
	//
	// There is no NATS connection anymore.
	logger.L.Sugar().Infow(
		"Closing infrastructure...",
	)

	closeCtx, closeCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer closeCancel()

	if err := app.Close(closeCtx); err != nil {
		logger.L.Sugar().Errorw(
			"infrastructure shutdown error",
			"error", err,
		)
	}

	logger.L.Sugar().Infow(
		"Server stopped successfully",
	)
}
