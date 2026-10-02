package infra

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

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
