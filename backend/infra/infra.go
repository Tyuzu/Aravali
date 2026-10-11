// File: infra/infra.go
package infra

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"scav/config"
	"scav/infra/cache"
	"scav/infra/logger"
	"scav/infra/mq"
	"scav/infra/sqldb"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Deps struct {
	SQLDB  sqldb.Database
	Cache  cache.Cache
	MQ     mq.MQ
	Config config.Config
	Log    *zap.Logger

	// Underlying raw client for graceful shutdown (Redis MQ relies on it).
	RedisClient *redis.Client
}

/* -------------------- Constructor -------------------- */

func New(cfg *config.Config) (*Deps, error) {
	if cfg == nil {
		return nil, errors.New("nil config provided")
	}

	d := &Deps{
		Config: *cfg,
	}

	// Initialize logging early so components can use it.
	if err := logger.Init(); err != nil {
		return nil, fmt.Errorf("logger init: %w", err)
	}
	d.Log = logger.L

	// Helper to clean up partially initialized resources on error.
	cleanup := func() {
		_ = d.Close(context.Background())
	}

	/* -------- Redis -------- */

	redisAddr := env("REDIS_ADDR", "localhost:6379")
	redisPassword := env("REDIS_PASSWORD", "")
	redisDB := 0

	redisClient, err := NewRedis(redisAddr, redisPassword, redisDB)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("redis setup: %w", err)
	}

	d.RedisClient = redisClient
	d.Cache = cache.NewCache(redisClient)

	d.MQ = mq.NewRedisMQ(redisClient)
	if d.MQ == nil {
		cleanup()
		return nil, errors.New("redis MQ initialization returned nil")
	}

	/* -------- Postgres & Migrations -------- */

	pgURI := env("DATABASE_URL", "")
	if strings.TrimSpace(pgURI) != "" {
		pool, err := NewPostgres(pgURI)
		if err != nil {
			cleanup()
			return nil, fmt.Errorf("postgres setup: %w", err)
		}

		sqlDB, err := sqldb.NewDatabase(pool)
		if err != nil {
			pool.Close()
			cleanup()
			return nil, fmt.Errorf("sqldb init: %w", err)
		}
		d.SQLDB = sqlDB

		// Run DB migrations if present.
		migrationsPath := env("MIGRATIONS_PATH", "./migrations")
		absMigrations, err := filepath.Abs(migrationsPath)
		if err == nil {
			fixed := filepath.ToSlash(absMigrations)
			src := "file://" + fixed
			m, err := migrate.New(src, pgURI)
			if err != nil {
				d.Log.Sugar().Warnw("migration init failed", "error", err)
			} else {
				if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
					d.Log.Sugar().Warnw("migrations up failed", "error", err)
				} else {
					d.Log.Sugar().Infow("migrations applied", "path", absMigrations)
				}
				// Ensure migration instance is closed to avoid leaks
				_, _ = m.Close()
			}
		}
	}

	return d, nil
}

/* -------------------- Graceful Shutdown -------------------- */

func (d *Deps) Close(ctx context.Context) error {
	if d == nil {
		return nil
	}

	if ctx == nil {
		ctx = context.Background()
	}

	var errs []string

	/*
		1. Close PostgreSQL via the Database interface.
	*/
	if d.SQLDB != nil {
		d.SQLDB.Close()
		d.SQLDB = nil
	}

	/*
		2. Close Redis.
	*/
	if d.RedisClient != nil {
		if err := d.RedisClient.Close(); err != nil && !errors.Is(err, redis.ErrClosed) {
			errs = append(errs, fmt.Sprintf("redis close: %v", err))
		}
		d.RedisClient = nil
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %s", strings.Join(errs, "; "))
	}

	return nil
}

/* -------------------- Redis -------------------- */

func NewRedis(addr string, password string, dbIndex int) (*redis.Client, error) {
	if strings.TrimSpace(addr) == "" {
		return nil, errors.New("redis address is empty")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbIndex,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return client, nil
}

/* -------------------- Postgres -------------------- */

func NewPostgres(uri string) (*pgxpool.Pool, error) {
	if strings.TrimSpace(uri) == "" {
		return nil, errors.New("postgres URL is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

/* -------------------- Helpers -------------------- */

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
