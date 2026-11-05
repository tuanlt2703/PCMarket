package infrastructure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func Connect(ctx context.Context, url string) (*DB, error) {
	pgPool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("Failed to create pool: %v", err)
	}

	if err := pgPool.Ping(ctx); err != nil {
		pgPool.Close()
		return nil, fmt.Errorf("Failed to connect database: %v", err)
	}
	return &DB{pool: pgPool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DB) Conn() *pgxpool.Pool {
	return db.pool
}
