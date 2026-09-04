package infra

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"

	"scav/config"
	"scav/infra/cache"
	"scav/infra/db"
	"scav/infra/mq"
	"scav/utils/logger"
)

type Deps struct {
	DB       db.Database
	Cache    cache.Cache
	MQ       mq.MQ
	NatsConn *nats.Conn
	Config   config.Config
}

/* -------------------- Constructor -------------------- */

func New(cfg *config.Config) (*Deps, error) {
	/* -------- Postgres -------- */

	postgresURL := cfg.DatabaseURL
	if postgresURL == "" {
		postgresURL = env("POSTGRES_URL", env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/scav"))
	}

	pool, err := NewPostgres(postgresURL)
	if err != nil {
		return nil, err
	}

	dbLayer := db.NewPostgresDatabase(pool, 100)

	/* -------- Redis -------- */

	redisAddr := env("REDIS_ADDR", "localhost:6379")
	redisPassword := env("REDIS_PASSWORD", "")
	redisDB := 0

	rclient, err := NewRedis(redisAddr, redisPassword, redisDB)
	if err != nil {
		return nil, err
	}
	cacheLayer := cache.NewRedisCache(rclient)

	/* -------- NATS JetStream (optional) -------- */

	var mqLayer mq.MQ
	var nc *nats.Conn

	natsURL := env("NATS_URL", "")
	if natsURL != "" {
		conn, js, err := NewJetStream(natsURL)
		if err != nil {
			return nil, err
		}

		mqLayer = mq.NewJetStreamMQ(js)
		nc = conn
	}

	logger.L.Sugar().Infow("infra initialized", "nats_enabled", natsURL != "")

	return &Deps{
		DB:       dbLayer,
		Cache:    cacheLayer,
		MQ:       mqLayer,
		NatsConn: nc,
		Config:   *cfg,
	}, nil
}

/* -------------------- Helpers -------------------- */

func env(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

/* -------------------- Postgres -------------------- */

func NewPostgres(uri string) (*pgxpool.Pool, error) {
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

/* -------------------- Redis -------------------- */

func NewRedis(addr string, password string, dbIndex int) (*redis.Client, error) {
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

/* -------------------- NATS -------------------- */

func NewJetStream(url string) (*nats.Conn, nats.JetStreamContext, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		_ = nc.Drain()
		return nil, nil, err
	}

	return nc, js, nil
}
