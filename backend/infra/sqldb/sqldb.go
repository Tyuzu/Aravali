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
	Offset  int       // Replaces MongoDB's Skip
	Sort    []OrderBy // Replaces ordered MongoDB sorting
	Columns []string  // Replaces MongoDB Projection
}

// Database defines a standard PostgreSQL database abstraction layer.
type Database interface {
	/* Lifecycle */
	Ping(ctx context.Context) error
	WithDB(ctx context.Context, op func(ctx context.Context) error) error
	RunTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error

	/* Create */
	Insert(ctx context.Context, table string, record any) error
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
	Update(ctx context.Context, table string, query string, args []any, updateValues map[string]any) (int64, error)
	UpdateOne(ctx context.Context, table string, query string, args []any, updateValues map[string]any) (int64, error)
	UpdateMany(ctx context.Context, table string, query string, args []any, updateValues map[string]any) (int64, error)
	Upsert(ctx context.Context, table string, conflictColumn string, record any) error
	Inc(ctx context.Context, table string, query string, args []any, column string, value int64) error
	AddToSet(ctx context.Context, table string, query string, args []any, arrayColumn string, value any) error

	/* Delete */
	Delete(ctx context.Context, table string, query string, args []any) (int64, error)
	DeleteOne(ctx context.Context, table string, query string, args []any) (int64, error)
	DeleteMany(ctx context.Context, table string, query string, args []any) (int64, error)

	/* Atomic */
	FindOneAndUpdate(ctx context.Context, table string, query string, args []any, updateValues map[string]any, result any) error

	/* Aggregate / Count */
	QueryRaw(ctx context.Context, sqlQuery string, args []any, result any) error
	Count(ctx context.Context, table string, query string, args []any) (int64, error)
	CountDocuments(ctx context.Context, table string, query string, args []any) (int64, error)
	EstimatedDocumentCount(ctx context.Context, table string) (int64, error)
}
