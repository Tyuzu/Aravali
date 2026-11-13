package config

import "testing"

func TestInitConfigLoadsDatabaseSettings(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("MONGO_URI", "mongodb://mongo:27017")
	t.Setenv("MONGO_DB", "scav")
	t.Setenv("REDIS_URL", "redis://:secret@redis:6379/2")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_DB", "2")

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("InitConfig returned error: %v", err)
	}

	if cfg.MongoURI != "mongodb://mongo:27017" {
		t.Fatalf("MongoURI = %q, want %q", cfg.MongoURI, "mongodb://mongo:27017")
	}
	if cfg.MongoDB != "scav" {
		t.Fatalf("MongoDB = %q, want %q", cfg.MongoDB, "scav")
	}
	if cfg.RedisAddr != "redis:6379" {
		t.Fatalf("RedisAddr = %q, want %q", cfg.RedisAddr, "redis:6379")
	}
	if cfg.RedisPassword != "secret" {
		t.Fatalf("RedisPassword = %q, want %q", cfg.RedisPassword, "secret")
	}
	if cfg.RedisDB != 2 {
		t.Fatalf("RedisDB = %d, want %d", cfg.RedisDB, 2)
	}
}

func TestInitConfigRequiresProductionDatabaseSettings(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv("ALLOWED_ORIGINS", "https://example.com")
	t.Setenv("REDIS_URL", "")
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("MONGO_URI", "")

	if _, err := InitConfig(); err == nil {
		t.Fatal("InitConfig() = nil error, want production validation error")
	}
}
