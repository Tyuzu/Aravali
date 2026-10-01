// File: infra/infra.go
package infra

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"scav/config"
	"scav/infra/cache"
	"scav/infra/mq"
	"scav/infra/sqldb"
	"scav/utils/logger"
)

type Deps struct {
	SQLDB  sqldb.Database
	Cache  cache.Cache
	MQ     mq.MQ
	Config config.Config

	// Underlying raw clients for graceful shutdown.
	PGPool      *pgxpool.Pool
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

	// Helper to clean up partially initialized resources on error.
	cleanup := func() {
		_ = d.Close(context.Background())
	}

	/* -------- Redis -------- */

	redisAddr := env(
		"REDIS_ADDR",
		"localhost:6379",
	)

	redisPassword := env(
		"REDIS_PASSWORD",
		"",
	)

	redisDB := 0

	redisClient, err := NewRedis(
		redisAddr,
		redisPassword,
		redisDB,
	)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("redis setup: %w", err)
	}

	d.RedisClient = redisClient

	// Redis cache uses the same go-redis client.
	d.Cache = cache.NewRedisCache(redisClient)

	/*
		Redis Pub/Sub is now the application's MQ.

		The same Redis server can safely be used for:
		- Cache
		- Pub/Sub

		go-redis handles the required Pub/Sub connections
		internally.
	*/
	d.MQ = mq.NewRedisMQ(redisClient)

	if d.MQ == nil {
		cleanup()
		return nil, errors.New("redis MQ initialization returned nil")
	}

	/* -------- Postgres -------- */

	postgresURL := cfg.DatabaseURL

	if postgresURL == "" {
		postgresURL = env(
			"POSTGRES_URL",
			env("DATABASE_URL", ""),
		)
	}

	if postgresURL == "" {
		user := env(
			"POSTGRES_USER",
			"apeman",
		)

		pass := env(
			"POSTGRES_PASSWORD",
			"ningning",
		)

		host := env(
			"POSTGRES_HOST",
			"localhost",
		)

		port := env(
			"POSTGRES_PORT",
			"5432",
		)

		dbname := env(
			"POSTGRES_DB",
			"eventdb",
		)

		postgresURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			user,
			pass,
			host,
			port,
			dbname,
		)
	}

	pool, err := NewPostgres(postgresURL)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("postgres setup: %w", err)
	}

	d.PGPool = pool
	d.SQLDB = sqldb.NewPostgresDatabase(
		pool,
		100,
	)

	logger.L.Sugar().Infow(
		"infra initialized",
		"redis_enabled", true,
		"redis_addr", redisAddr,
		"mq", "redis_pubsub",
	)

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
		1. Close PostgreSQL.
	*/
	if d.PGPool != nil {
		d.PGPool.Close()
	}

	/*
		2. Close Redis.

		The Redis MQ implementation creates/owns its Pub/Sub
		connections through the underlying Redis client, so closing
		the client shuts down the Redis resources as well.
	*/
	if d.RedisClient != nil {
		if err := d.RedisClient.Close(); err != nil {
			errs = append(
				errs,
				fmt.Sprintf("redis close: %v", err),
			)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf(
			"close errors: %s",
			strings.Join(errs, "; "),
		)
	}

	return nil
}

/* -------------------- Redis -------------------- */

func NewRedis(
	addr string,
	password string,
	dbIndex int,
) (*redis.Client, error) {
	if strings.TrimSpace(addr) == "" {
		return nil, errors.New("redis address is empty")
	}

	client := redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       dbIndex,
		},
	)

	ctx, cancel := cancelTimeout(5 * time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()

		return nil, err
	}

	return client, nil
}

/* -------------------- Postgres -------------------- */

func NewPostgres(
	uri string,
) (*pgxpool.Pool, error) {
	if strings.TrimSpace(uri) == "" {
		return nil, errors.New(
			"postgres URL is empty",
		)
	}

	ctx, cancel := cancelTimeout(10 * time.Second)
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

func env(
	key string,
	fallback string,
) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func cancelTimeout(
	duration time.Duration,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		context.Background(),
		duration,
	)
}
