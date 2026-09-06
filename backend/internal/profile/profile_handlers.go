package profile

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"scav/infra/cache"
	"scav/middleware"
)

/* -------------------------------------------------------
   Helpers
------------------------------------------------------- */

// validateJWT extracts + validates JWT from header
func validateJWT(r *http.Request) (*middleware.Claims, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, errors.New("no auth header")
	}
	return middleware.ValidateJWT(token)
}

/* ------------------------------------------------------
/* -------------------------------------------------------
   User profile endpoints
------------------------------------------------------- */

// isOnline checks if a user is online via cache
func isOnline(ctx context.Context, userid string, cache cache.Cache) (bool, error) {
	return cache.Exists(ctx, "online:"+userid)
}

// CacheProfile stores the serialized profile in cache
func CacheProfile(ctx context.Context, c cache.Cache, username string, data string, ttl time.Duration) error {
	return c.Set(ctx, "profile:"+username, []byte(data), ttl)
}

// GetCachedProfile fetches the cached profile
func GetCachedProfile(ctx context.Context, c cache.Cache, username string) (string, error) {
	data, err := c.Get(ctx, "profile:"+username)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// InvalidateCachedProfile deletes the cached profile
func InvalidateCachedProfile(ctx context.Context, c cache.Cache, username string) error {
	return c.Del(ctx, "profile:"+username)
}

// UpdateCachedUsername invalidates cached data keyed by user ID
func UpdateCachedUsername(ctx context.Context, c cache.Cache, userid string) error {
	return c.Del(ctx, fmt.Sprintf("users:%s", userid))
}
