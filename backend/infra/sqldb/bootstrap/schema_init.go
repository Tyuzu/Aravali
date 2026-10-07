package bootstrap

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"scav/infra/sqldb/schema"
)

// EnsureSchema is the module-level database bootstrap hook.
func EnsureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	columns := []schema.ColumnDefinition{
		{Name: "id", Definition: "TEXT"},
		{Name: "status", Definition: "TEXT"},
		{Name: "metadata", Definition: "JSONB NOT NULL DEFAULT '{}'::JSONB"},
		{Name: "created_at", Definition: "TIMESTAMPTZ NOT NULL DEFAULT NOW()"},
		{Name: "updated_at", Definition: "TIMESTAMPTZ NOT NULL DEFAULT NOW()"},
	}

	return schema.EnsureColumns(ctx, pool, "bootstrap_state", columns)
}
