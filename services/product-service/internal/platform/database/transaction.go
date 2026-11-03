package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// transactionKey is the context key that stores the active pgx.Tx. Using a
// private type guarantees that no other package can overwrite the value.
type transactionKey struct{}

func transactionFrom(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(transactionKey{}).(pgx.Tx)
	if !ok {
		return nil
	}
	return tx
}

// RunInTransaction runs fn inside a transaction.
//
// The transaction is passed to fn through its ctx argument, so every
// repository call made with that ctx (via DB.Conn) joins the transaction.
// It commits when fn returns nil and rolls back when fn returns an error or
// panics. A nested call joins the outer transaction instead of starting a new one.
func (db *DB) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if transactionFrom(ctx) != nil {
		return fn(ctx)
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	committed := false
	defer func() {
		// Runs on every exit path, including a panic inside fn.
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	if err := fn(context.WithValue(ctx, transactionKey{}, tx)); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	committed = true
	return nil
}
