package sqldb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNilPool         = fmt.Errorf("pool is nil")
	ErrUninitializedDB = fmt.Errorf("db not initialized")
)

// Database defines the application-level database interface,
// enabling clean separation and easy mocking in unit tests.
type Database interface {
	Close()
	Ping(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// PostgresDB wraps a pgxpool.Pool and implements the Database interface.
type PostgresDB struct {
	pool *pgxpool.Pool
}

// NewDatabase creates a new Database adapter around a pgx pool.
func NewDatabase(pool *pgxpool.Pool) (Database, error) {
	if pool == nil {
		return nil, ErrNilPool
	}
	return &PostgresDB{pool: pool}, nil
}

// Close closes the underlying database connection pool.
func (p *PostgresDB) Close() {
	if p == nil || p.pool == nil {
		return
	}
	p.pool.Close()
}

// Ping verifies connectivity to the database.
func (p *PostgresDB) Ping(ctx context.Context) error {
	if p == nil || p.pool == nil {
		return ErrUninitializedDB
	}
	return p.pool.Ping(ctx)
}

// Exec executes a query without returning any rows.
func (p *PostgresDB) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, arguments...)
}

// Query executes a query that returns rows.
func (p *PostgresDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

// QueryRow executes a query that is expected to return at most one row.
func (p *PostgresDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

// SendBatch sends a batch of SQL statements.
func (p *PostgresDB) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return p.pool.SendBatch(ctx, b)
}

// Begin starts a transaction.
func (p *PostgresDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.pool.Begin(ctx)
}

// BeginTx starts a transaction with specific options.
func (p *PostgresDB) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.pool.BeginTx(ctx, txOptions)
}
