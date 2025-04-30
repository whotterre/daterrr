package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxQuerier represents a transaction with all query methods.
type TxQuerier interface {
	Querier
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Store interface embeds the generated sqlc Querier and adds transaction methods
type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(Querier) error) error
	BeginTx(ctx context.Context) (TxQuerier, error)
}

// SQLStore holds a connection pool and implements Store
type SQLStore struct {
	pool    *pgxpool.Pool
	*Queries
}

// txQuerier implements TxQuerier
type txQuerier struct {
	*Queries
	tx pgx.Tx
}

func (t *txQuerier) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *txQuerier) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

// NewStore creates a new store using a pgxpool.Pool
func NewStore(pool *pgxpool.Pool) Store {
	var dbi DBTX = pool
	return &SQLStore{
		pool:    pool,
		Queries: New(dbi),
	}
}

// BeginTx starts a new transaction and returns a TxQuerier
func (s *SQLStore) BeginTx(ctx context.Context) (TxQuerier, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	
	queries := New(tx)
	return &txQuerier{
		Queries: queries,
		tx:      tx,
	}, nil
}

// ExecTx executes a function within a database transaction
func (s *SQLStore) ExecTx(ctx context.Context, fn func(Querier) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Rollback is safe even if Commit succeeds

	q := New(tx) // Wrap tx with sqlc Queries

	err = fn(q)
	if err != nil {
		return err // Rollback will trigger due to defer
	}

	return tx.Commit(ctx)
}