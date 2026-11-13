package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                   string
	MongoURI              string
	MongoDB               string
	RedisAddr             string
	RedisPassword         string
	RedisDB               int
	DatabaseURL           string
	PostgresUser          string
	PostgresPassword      string
	PostgresHost          string
	PostgresPort          string
	PostgresDB            string
	STRIPE_WEBHOOK_SECRET string
	HTTPPort              string
	AllowedOrigins        []string
	RTMPIngestURL         string
	TURNServers           string
	CDNBaseURL            string
	// TLSCertPath           string
	// TLSKeyPath            string
	AllowCredentials bool
	// TerminateTLSAtLB      bool // Added to support external TLS termination
}

// InitConfig loads environment variables and performs basic validation.
// Returns an error if required configuration is missing or inconsistent.
func InitConfig() (*Config, error) {
	_ = godotenv.Load()

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		if env == "production" {
			return nil, fmt.Errorf("MONGO_URI must be set in production")
		}
		mongoURI = "mongodb://localhost:27017"
	}

	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		mongoDB = "eventdb"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0
	if redisAddr == "" {
		if rawURL := os.Getenv("REDIS_URL"); rawURL != "" {
			parsedAddr, parsedPassword, parsedDB, err := parseRedisURL(rawURL)
			if err != nil {
				return nil, fmt.Errorf("invalid REDIS_URL: %w", err)
			}
			redisAddr = parsedAddr
			if redisPassword == "" {
				redisPassword = parsedPassword
			}
			redisDB = parsedDB
		}
	}
	if redisAddr == "" {
		if env == "production" {
			return nil, fmt.Errorf("REDIS_URL or REDIS_ADDR must be set in production")
		}
		redisAddr = "localhost:6379"
	}
	if v := os.Getenv("REDIS_DB"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			redisDB = parsed
		}
	}

	allowedOrigins := parseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"))

	allowCreds := true
	if v := os.Getenv("ALLOW_CREDENTIALS"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			allowCreds = b
		}
	}

	// // Parse TLS termination setting
	// terminateTLSAtLB := false
	// if v := os.Getenv("TERMINATE_TLS_AT_LB"); v != "" {
	// 	if b, err := strconv.ParseBool(v); err == nil {
	// 		terminateTLSAtLB = b
	// 	}
	// }

	// // Load Postgres discrete parameters with defaults
	// pgUser := envOrDefault("POSTGRES_USER", "apeman")
	// pgPass := envOrDefault("POSTGRES_PASSWORD", "ningning")
	// pgHost := envOrDefault("POSTGRES_HOST", "localhost")
	// pgPort := envOrDefault("POSTGRES_PORT", "5432")
	// pgDB := envOrDefault("POSTGRES_DB", "eventdb")

	// // Fallback to explicit connection URLs if available
	// dbURL := os.Getenv("POSTGRES_URL")
	// if dbURL == "" {
	// 	dbURL = os.Getenv("DATABASE_URL")
	// }

	// // Construct connection string if no full URL is provided explicitly
	// if dbURL == "" {
	// 	dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPass, pgHost, pgPort, pgDB)
	// }

	cfg := &Config{
		Env:           env,
		MongoURI:      mongoURI,
		MongoDB:       mongoDB,
		RedisAddr:     redisAddr,
		RedisPassword: redisPassword,
		RedisDB:       redisDB,
		// DatabaseURL:           dbURL,
		// PostgresUser:          pgUser,
		// PostgresPassword:      pgPass,
		// PostgresHost:          pgHost,
		// PostgresPort:          pgPort,
		// PostgresDB:            pgDB,
		HTTPPort:              port,
		AllowedOrigins:        allowedOrigins,
		STRIPE_WEBHOOK_SECRET: os.Getenv("STRIPE_WEBHOOK_SECRET"),
		RTMPIngestURL:         os.Getenv("RTMP_INGEST_URL"),
		TURNServers:           os.Getenv("TURN_SERVERS"),
		CDNBaseURL:            os.Getenv("CDN_BASE_URL"),
		// TLSCertPath:           os.Getenv("TLS_CERT_PATH"),
		// TLSKeyPath:            os.Getenv("TLS_KEY_PATH"),
		AllowCredentials: allowCreds,
		// TerminateTLSAtLB:      terminateTLSAtLB,
	}

	// Basic validation
	if len(cfg.AllowedOrigins) == 0 {
		return nil, errors.New("ALLOWED_ORIGINS must be set to a comma-separated list of allowed origins")
	}

	if cfg.Env == "production" {
		// // ensure essential services are configured in production
		// if cfg.DatabaseURL == "" {
		// 	return nil, fmt.Errorf("POSTGRES_URL or discrete POSTGRES_* environment variables must be set in production")
		// }
		if os.Getenv("REDIS_URL") == "" && os.Getenv("REDIS_ADDR") == "" {
			return nil, fmt.Errorf("REDIS_URL or REDIS_ADDR must be set in production")
		}

		// // Only enforce internal TLS if we aren't terminating it upstream at a load balancer
		// if !cfg.TerminateTLSAtLB {
		// 	if cfg.TLSCertPath == "" || cfg.TLSKeyPath == "" {
		// 		return nil, fmt.Errorf("TLS_CERT_PATH and TLS_KEY_PATH must be set in production or set TERMINATE_TLS_AT_LB=true")
		// 	}
		// }

		// Do not allow wildcard origins with credentials in production
		if cfg.AllowCredentials {
			for _, o := range cfg.AllowedOrigins {
				if strings.TrimSpace(o) == "*" {
					return nil, fmt.Errorf("ALLOW_CREDENTIALS=true is incompatible with ALLOWED_ORIGINS='*' in production")
				}
			}
		}
	}

	return cfg, nil
}

func parseAllowedOrigins(env string) []string {
	if env == "" {
		return []string{"http://localhost:5173"}
	}
	parts := strings.Split(env, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func parseRedisURL(rawURL string) (string, string, int, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", "", 0, nil
	}
	if !strings.Contains(rawURL, "://") {
		return rawURL, "", 0, nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", 0, err
	}
	if u.Host == "" {
		return "", "", 0, fmt.Errorf("missing host in REDIS_URL")
	}

	password := ""
	if u.User != nil {
		if p, ok := u.User.Password(); ok {
			password = p
		}
	}

	db := 0
	if u.Path != "" && u.Path != "/" {
		pathDB := strings.TrimPrefix(u.Path, "/")
		if pathDB != "" {
			if parsedDB, err := strconv.Atoi(pathDB); err == nil {
				db = parsedDB
			}
		}
	}

	return u.Host, password, db, nil
}

// func envOrDefault(key string, fallback string) string {
// 	if v := os.Getenv(key); v != "" {
// 		return v
// 	}
// 	return fallback
// }
