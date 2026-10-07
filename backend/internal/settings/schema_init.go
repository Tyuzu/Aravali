package settings

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EnsureSchema is the module-level database bootstrap hook.
func EnsureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("nil postgres pool")
	}

	columns := []struct {
		name       string
		definition string
	}{
		{name: "id", definition: "TEXT"},
		{name: "status", definition: "TEXT"},
		{name: "metadata", definition: "JSONB NOT NULL DEFAULT '{}'::JSONB"},
		{name: "created_at", definition: "TIMESTAMPTZ NOT NULL DEFAULT NOW()"},
		{name: "updated_at", definition: "TIMESTAMPTZ NOT NULL DEFAULT NOW()"},
	}

	return ensureTableColumns(ctx, pool, "settings", columns)
}

func ensureTableColumns(ctx context.Context, pool *pgxpool.Pool, tableName string, columns []struct {
	name       string
	definition string
}) error {
	if pool == nil {
		return fmt.Errorf("nil postgres pool")
	}

	exists, err := tableExists(ctx, pool, tableName)
	if err != nil {
		return fmt.Errorf("check %s table: %w", tableName, err)
	}
	if !exists {
		return nil
	}

	for _, column := range columns {
		present, err := columnExists(ctx, pool, tableName, column.name)
		if err != nil {
			return fmt.Errorf("check %s.%s: %w", tableName, column.name, err)
		}
		if present {
			continue
		}
		if _, err := pool.Exec(ctx, fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s %s", quoteIdent(tableName), quoteIdent(column.name), column.definition)); err != nil {
			return fmt.Errorf("add %s.%s: %w", tableName, column.name, err)
		}
	}
	return nil
}

func tableExists(ctx context.Context, pool *pgxpool.Pool, tableName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = current_schema()
          AND table_name = $1
    )`
	err := pool.QueryRow(ctx, query, normalizeIdentifier(tableName)).Scan(&exists)
	return exists, err
}

func columnExists(ctx context.Context, pool *pgxpool.Pool, tableName, columnName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = $1
          AND column_name = $2
    )`
	err := pool.QueryRow(ctx, query, normalizeIdentifier(tableName), normalizeIdentifier(columnName)).Scan(&exists)
	return exists, err
}

func normalizeIdentifier(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, `"`)
	return name
}

func quoteIdent(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	normalized := re.ReplaceAllString(name, "_")
	return `"` + strings.ToLower(normalized) + `"`
}
