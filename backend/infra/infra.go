package infra

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"scav/config"
	"scav/infra/cache"
	"scav/infra/db"
	"scav/infra/mq"
	"scav/infra/sqldb"
	"scav/utils/logger"
)

type Deps struct {
	SQLDB    sqldb.Database
	DB       db.Database
	Cache    cache.Cache
	MQ       mq.MQ
	NatsConn *nats.Conn
	Config   config.Config

	// Underlying raw clients for graceful shutdown
	PGPool      *pgxpool.Pool
	MongoClient *mongo.Client
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

	// Helper to cleanup partially initialized resources on error
	cleanup := func() {
		_ = d.Close(context.Background())
	}

	/* -------- Mongo -------- */
	mongoURI := env("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := env("MONGO_DB", "eventdb")

	client, database, err := NewMongo(mongoURI, mongoDB)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("mongo setup: %w", err)
	}
	d.MongoClient = client
	d.DB = db.NewMongoDatabase(database, client, 100)

	/* -------- Redis -------- */
	redisAddr := env("REDIS_ADDR", "localhost:6379")
	redisPassword := env("REDIS_PASSWORD", "")
	redisDB := 0

	rclient, err := NewRedis(redisAddr, redisPassword, redisDB)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("redis setup: %w", err)
	}
	d.RedisClient = rclient
	d.Cache = cache.NewRedisCache(rclient)

	/* -------- NATS JetStream (optional) -------- */
	natsURL := env("NATS_URL", "")
	if natsURL != "" {
		conn, js, err := NewJetStream(natsURL)
		if err != nil {
			cleanup()
			return nil, fmt.Errorf("nats setup: %w", err)
		}
		d.NatsConn = conn
		d.MQ = mq.NewJetStreamMQ(js)
	}

	/* -------- Postgres -------- */
	postgresURL := cfg.DatabaseURL
	if postgresURL == "" {
		postgresURL = env("POSTGRES_URL", env("DATABASE_URL", ""))
	}

	if postgresURL == "" {
		user := env("POSTGRES_USER", "apeman")
		pass := env("POSTGRES_PASSWORD", "ningning")
		host := env("POSTGRES_HOST", "localhost")
		port := env("POSTGRES_PORT", "5432")
		dbname := env("POSTGRES_DB", "eventdb")

		postgresURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, dbname)
	}

	pool, err := NewPostgres(postgresURL)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("postgres setup: %w", err)
	}
	d.PGPool = pool
	d.SQLDB = sqldb.NewPostgresDatabase(pool, 100)

	logger.L.Sugar().Infow("infra initialized", "nats_enabled", natsURL != "")

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

	// 1. Drain NATS first so subscriptions stop receiving new work
	if d.NatsConn != nil {
		if err := d.NatsConn.Drain(); err != nil {
			d.NatsConn.Close()
			errs = append(errs, fmt.Sprintf("nats drain: %v", err))
		}
	}

	// 2. Close Database & Cache Connections
	if d.PGPool != nil {
		d.PGPool.Close()
	}

	if d.RedisClient != nil {
		if err := d.RedisClient.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("redis close: %v", err))
		}
	}

	if d.MongoClient != nil {
		if err := d.MongoClient.Disconnect(ctx); err != nil {
			errs = append(errs, fmt.Sprintf("mongo disconnect: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

/* -------------------- Helpers -------------------- */

func env(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func NewMongo(uri string, dbName string) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(uri).
			SetMaxPoolSize(100).
			SetMinPoolSize(10).
			SetRetryWrites(true),
	)
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, nil, err
	}

	return client, client.Database(dbName), nil
}

func NewRedis(addr string, password string, dbIndex int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbIndex,
	})

	ctx, cancel := cancelTimeout(5 * time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return client, nil
}

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

func NewPostgres(uri string) (*pgxpool.Pool, error) {
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

func cancelTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}
