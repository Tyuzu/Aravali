package farms

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"scav/infra/sqldb/schema"
)

func EnsureFarmColumns(ctx context.Context, pool *pgxpool.Pool) error {
	columns := []schema.ColumnDefinition{
		{Name: "name", Definition: "TEXT"},
		{Name: "location", Definition: "TEXT"},
		{Name: "latitude", Definition: "DOUBLE PRECISION DEFAULT 0"},
		{Name: "longitude", Definition: "DOUBLE PRECISION DEFAULT 0"},
		{Name: "description", Definition: "TEXT"},
		{Name: "owner", Definition: "TEXT"},
		{Name: "contact", Definition: "TEXT"},
		{Name: "social", Definition: "TEXT"},
		{Name: "practice", Definition: "TEXT"},
		{Name: "availability", Definition: "JSONB DEFAULT '{}'::JSONB"},
		{Name: "tags", Definition: "TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[]"},
		{Name: "crops", Definition: "JSONB DEFAULT '[]'::JSONB"},
		{Name: "banner", Definition: "TEXT"},
		{Name: "media", Definition: "TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[]"},
		{Name: "avg_rating", Definition: "DOUBLE PRECISION NOT NULL DEFAULT 0"},
		{Name: "review_count", Definition: "INTEGER NOT NULL DEFAULT 0"},
		{Name: "favorites_count", Definition: "BIGINT NOT NULL DEFAULT 0"},
		{Name: "created_by", Definition: "TEXT"},
		{Name: "updated_by", Definition: "TEXT"},
		{Name: "metadata", Definition: "JSONB NOT NULL DEFAULT '{}'::JSONB"},
	}

	return schema.EnsureColumns(ctx, pool, "farms", columns)
}
