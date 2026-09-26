package sqldb

import (
	"context"
	"database/sql"
)

// OrderBy defines sorting criteria for queries.
type OrderBy struct {
	Column     string
	Descending bool
}

// FindManyOptions provides pagination, sorting, and field selection.
type FindManyOptions struct {
	Limit   int
	Offset  int
	OrderBy string    // Updated to support string formats like "created_at ASC"
	Sort    []OrderBy // Kept for structured multi-column sorting if needed
	Columns []string
}

// Database defines a standard PostgreSQL database abstraction layer.
type Database interface {
	/* Lifecycle & Transactions */
	Ping(ctx context.Context) error
	WithDB(ctx context.Context, op func(ctx context.Context) error) error
	RunTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error

	/* Create */
	InsertOne(ctx context.Context, table string, record any) error
	InsertMany(ctx context.Context, table string, records []any) error
	BulkWrite(ctx context.Context, table string, operations []any) error

	/* Read */
	FindOne(ctx context.Context, table string, query string, args []any, result any) error
	FindOneWithProjection(ctx context.Context, table string, columns []string, query string, args []any, result any) error

	FindMany(ctx context.Context, table string, query string, args []any, result any) error
	FindManyWithOptions(ctx context.Context, table string, query string, args []any, opts FindManyOptions, result any) error
	FindManyWithProjection(
		ctx context.Context,
		table string,
		query string,
		args []any,
		columns []string,
		opts FindManyOptions,
		result any,
	) error

	Distinct(ctx context.Context, table string, column string, query string, args []any, result any) error

	/* Update */
	UpdateOne(ctx context.Context, table string, query string, args []any, updateValues map[string]any) (int64, error)
	UpdateMany(ctx context.Context, table string, query string, args []any, updateValues map[string]any) (int64, error)
	Upsert(ctx context.Context, table string, conflictColumn string, record any) error
	Inc(ctx context.Context, table string, query string, args []any, column string, value int64) error

	/* Delete */
	DeleteOne(ctx context.Context, table string, query string, args []any) (int64, error)
	DeleteMany(ctx context.Context, table string, query string, args []any) (int64, error)

	/* Atomic Operations */
	FindOneAndUpdate(ctx context.Context, table string, query string, args []any, updateValues map[string]any, result any) error

	/* Raw Execution / Count / Aggregate */
	QueryRaw(ctx context.Context, sqlQuery string, args []any, result any) error
	Count(ctx context.Context, table string, query string, args []any) (int64, error)
	CountDocuments(ctx context.Context, table string, query string, args []any) (int64, error)
	Aggregate(ctx context.Context, table string, query string, args []any) error
}
