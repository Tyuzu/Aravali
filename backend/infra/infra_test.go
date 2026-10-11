package infra

import (
	"context"
	"testing"

	"scav/config"

	"github.com/redis/go-redis/v9"
)

func TestResolveDatabaseURI_PrefersConfigValue(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://wrong:pass@localhost:5432/wrongdb")
	cfg := &config.Config{DatabaseURL: "postgres://apeman:ningning@localhost:5432/eventdb"}

	if got := resolveDatabaseURI(cfg); got != cfg.DatabaseURL {
		t.Fatalf("resolveDatabaseURI() = %q, want %q", got, cfg.DatabaseURL)
	}
}

func TestResolveDatabaseURI_FallsBackToEnvironment(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://env:pass@localhost:5432/envdb")
	cfg := &config.Config{}

	if got := resolveDatabaseURI(cfg); got != "postgres://env:pass@localhost:5432/envdb" {
		t.Fatalf("resolveDatabaseURI() = %q, want env DATABASE_URL", got)
	}
}

func TestDepsClose_IsIdempotent(t *testing.T) {
	d := &Deps{
		RedisClient: redis.NewClient(&redis.Options{Addr: "localhost:6379"}),
	}

	if err := d.Close(context.Background()); err != nil {
		t.Fatalf("first close should succeed: %v", err)
	}

	if err := d.Close(context.Background()); err != nil {
		t.Fatalf("second close should be a no-op: %v", err)
	}
}
