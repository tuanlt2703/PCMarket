package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Querier is the set of query methods shared by *pgxpool.Pool and pgx.Tx.
// Repositories only use this interface, so the same code works both inside and
// outside a transaction.
type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// DB owns the PostgreSQL connection pool.
type DB struct {
	pool *pgxpool.Pool
}

// Connect opens the pool and checks that the database is reachable.
func Connect(ctx context.Context, url string) (*DB, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return &DB{pool: pool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Conn returns the connection to run a query on: the transaction started by
// RunInTransaction when ctx carries one, otherwise the pool.
//
// Always call it with the ctx of the current request, never a stored one,
// otherwise the query silently runs outside the transaction.
func (db *DB) Conn(ctx context.Context) Querier {
	if tx := transactionFrom(ctx); tx != nil {
		return tx
	}
	return db.pool
}
