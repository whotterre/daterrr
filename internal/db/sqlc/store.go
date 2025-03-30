package db

import (
    "context"
    "github.com/jackc/pgx/v5"
)

type Store interface {
    Querier
    ExecTx(ctx context.Context, fn func(Querier) error) error
	Conn() *pgx.Conn
}

type SQLStore struct {
    conn *pgx.Conn
    *Queries
}

func NewStore(conn *pgx.Conn) Store {
    return &SQLStore{
        conn: conn,
        Queries: New(conn),
    }
}

func (s *SQLStore) Conn() *pgx.Conn {
    return s.conn
}

func (s *SQLStore) ExecTx(ctx context.Context, fn func(Querier) error) error {
    tx, err := s.conn.Begin(ctx)
    if err != nil {
        return err
    }

    q := New(tx)
    err = fn(q)
    if err != nil {
        if rbErr := tx.Rollback(ctx); rbErr != nil {
            return rbErr
        }
        return err
    }

    return tx.Commit(ctx)
}