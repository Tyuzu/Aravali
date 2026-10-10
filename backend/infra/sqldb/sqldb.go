package sqldb

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database is the application-level database interface.
type Database interface {
	// Close closes the underlying database connection pool.
	Close()
	Ping(ctx context.Context) error
}

// PostgresDB is a thin wrapper around a pgxpool.Pool.
type PostgresDB struct {
	Pool *pgxpool.Pool
}

func (p *PostgresDB) Close() {
	if p == nil || p.Pool == nil {
		return
	}
	p.Pool.Close()
}

// Ping verifies connectivity to the database.
func (p *PostgresDB) Ping(ctx context.Context) error {
	if p == nil || p.Pool == nil {
		return nil
	}
	return p.Pool.Ping(ctx)
}

// NewDatabase creates a Database adapter around a pgx pool.
func NewDatabase(pool *pgxpool.Pool) Database {
	if pool == nil {
		return nil
	}
	return &PostgresDB{Pool: pool}
}
