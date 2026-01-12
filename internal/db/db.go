// Package db provides the database pool and transactional store helpers.
package db

import (
	"context"
	"edugov-back-v2/internal/db/sqlc"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a configured pgx connection pool.
func NewPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnIdleTime = 5 * time.Minute

	return pgxpool.NewWithConfig(ctx, cfg)
}

// Store wraps sqlc queries and the underlying pool for transactions.
type Store struct {
	*sqlc.Queries
	pool *pgxpool.Pool
}

// NewStore builds a Store from the pgx pool.
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Queries: sqlc.New(pool),
		pool:    pool,
	}
}

// ExecTx executes a function within database transaction
// Ensures fn are fully commited or rolled back together
// Uses pgxpool for transaction methods
func (s *Store) ExecTx(
	ctx context.Context,
	fn func(q *sqlc.Queries) error,
) (err error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	committed := false
	defer func() {
		// handle panic: rollback then rethrow
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
		// rollback if not committed
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	q := s.Queries.WithTx(tx)

	if err = fn(q); err != nil {
		return fmt.Errorf("tx fn: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	committed = true
	return nil
}
