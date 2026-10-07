package schema

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ColumnDefinition describes a single database column to ensure exists.
type ColumnDefinition struct {
	Name       string
	Definition string
}

// TableDefinition describes a table and the columns expected to exist.
type TableDefinition struct {
	Name    string
	Columns []ColumnDefinition
}

// EnsureColumns adds missing columns to an existing table without failing if the table is absent.
func EnsureColumns(ctx context.Context, pool *pgxpool.Pool, tableName string, columns []ColumnDefinition) error {
	if pool == nil {
		return fmt.Errorf("nil postgres pool")
	}
	if tableName == "" {
		return fmt.Errorf("table name is empty")
	}

	exists, err := TableExists(ctx, pool, tableName)
	if err != nil {
		return fmt.Errorf("check %s table: %w", tableName, err)
	}
	if !exists {
		return nil
	}

	for _, column := range columns {
		present, err := ColumnExists(ctx, pool, tableName, column.Name)
		if err != nil {
			return fmt.Errorf("check %s.%s: %w", tableName, column.Name, err)
		}
		if present {
			continue
		}
		if _, err := pool.Exec(ctx, fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s %s", QuoteIdent(tableName), QuoteIdent(column.Name), column.Definition)); err != nil {
			return fmt.Errorf("add %s.%s: %w", tableName, column.Name, err)
		}
	}

	return nil
}

// EnsureTable is a convenience wrapper for a complete table definition.
func EnsureTable(ctx context.Context, pool *pgxpool.Pool, table TableDefinition) error {
	return EnsureColumns(ctx, pool, table.Name, table.Columns)
}

func TableExists(ctx context.Context, pool *pgxpool.Pool, tableName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1
		FROM information_schema.tables
		WHERE table_schema = current_schema()
		  AND table_name = $1
	)`
	err := pool.QueryRow(ctx, query, NormalizeIdentifier(tableName)).Scan(&exists)
	return exists, err
}

func ColumnExists(ctx context.Context, pool *pgxpool.Pool, tableName, columnName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE table_schema = current_schema()
		  AND table_name = $1
		  AND column_name = $2
	)`
	err := pool.QueryRow(ctx, query, NormalizeIdentifier(tableName), NormalizeIdentifier(columnName)).Scan(&exists)
	return exists, err
}

func NormalizeIdentifier(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, `"`)
	return name
}

func QuoteIdent(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	normalized := re.ReplaceAllString(name, "_")
	return `"` + strings.ToLower(normalized) + `"`
}
