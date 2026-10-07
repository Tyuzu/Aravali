package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"scav/infra/sqldb/schema"
)

func EnsureUserColumns(ctx context.Context, pool *pgxpool.Pool) error {
	columns := []schema.ColumnDefinition{
		{Name: "userid", Definition: "TEXT"},
		{Name: "username", Definition: "TEXT"},
		{Name: "email", Definition: "TEXT"},
		{Name: "phone", Definition: "TEXT"},
		{Name: "phone_number", Definition: "TEXT"},
		{Name: "name", Definition: "TEXT"},
		{Name: "avatar", Definition: "TEXT"},
		{Name: "banner", Definition: "TEXT"},
		{Name: "bio", Definition: "TEXT"},
		{Name: "status", Definition: "TEXT"},
		{Name: "address", Definition: "TEXT"},
		{Name: "password", Definition: "TEXT"},
		{Name: "password_hash", Definition: "TEXT"},
		{Name: "role", Definition: "TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[]"},
		{Name: "email_verified", Definition: "BOOLEAN NOT NULL DEFAULT FALSE"},
		{Name: "is_verified", Definition: "BOOLEAN NOT NULL DEFAULT FALSE"},
		{Name: "online", Definition: "BOOLEAN NOT NULL DEFAULT FALSE"},
		{Name: "last_login", Definition: "TIMESTAMPTZ"},
		{Name: "refresh_token", Definition: "TEXT"},
		{Name: "refresh_prev", Definition: "TEXT"},
		{Name: "refresh_expiry", Definition: "TIMESTAMPTZ"},
		{Name: "refresh_ua", Definition: "TEXT"},
		{Name: "refresh_ip", Definition: "TEXT"},
		{Name: "profile_views", Definition: "INTEGER NOT NULL DEFAULT 0"},
		{Name: "followerscount", Definition: "INTEGER NOT NULL DEFAULT 0"},
		{Name: "followscount", Definition: "INTEGER NOT NULL DEFAULT 0"},
		{Name: "social_links", Definition: "JSONB DEFAULT '{}'::JSONB"},
		{Name: "wallet_balance", Definition: "NUMERIC(18,2) NOT NULL DEFAULT 0"},
		{Name: "created_at", Definition: "TIMESTAMPTZ NOT NULL DEFAULT NOW()"},
		{Name: "updated_at", Definition: "TIMESTAMPTZ NOT NULL DEFAULT NOW()"},
		{Name: "metadata", Definition: "JSONB NOT NULL DEFAULT '{}'::JSONB"},
		{Name: "followings", Definition: "TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[]"},
	}

	return schema.EnsureColumns(ctx, pool, "users", columns)
}
