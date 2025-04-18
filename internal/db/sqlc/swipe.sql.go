// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: swipe.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkMutualSwipe = `-- name: CheckMutualSwipe :one
SELECT EXISTS (
    SELECT 1 FROM swipes 
    WHERE 
        swiper_id = $1 AND  -- The other user
        swipee_id = $2      -- Current user
) AS is_mutual
`

type CheckMutualSwipeParams struct {
	SwiperID pgtype.UUID `json:"swiper_id"`
	SwipeeID pgtype.UUID `json:"swipee_id"`
}

func (q *Queries) CheckMutualSwipe(ctx context.Context, arg CheckMutualSwipeParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkMutualSwipe, arg.SwiperID, arg.SwipeeID)
	var is_mutual bool
	err := row.Scan(&is_mutual)
	return is_mutual, err
}

const createMatch = `-- name: CreateMatch :one
WITH new_match AS (
  INSERT INTO matches (user1_id, user2_id)
  VALUES (
    LEAST($1::uuid, $2::uuid),
    GREATEST($1::uuid, $2::uuid)
  )
  ON CONFLICT (user1_id, user2_id) DO NOTHING
  RETURNING id, user1_id, user2_id, matched_at
)
INSERT INTO chats (user1_id, user2_id)
SELECT user1_id, user2_id FROM new_match
ON CONFLICT (user1_id, user2_id) DO NOTHING
RETURNING (SELECT id FROM new_match)
`

type CreateMatchParams struct {
	Column1 pgtype.UUID `json:"column_1"`
	Column2 pgtype.UUID `json:"column_2"`
}

func (q *Queries) CreateMatch(ctx context.Context, arg CreateMatchParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, createMatch, arg.Column1, arg.Column2)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const newSwipe = `-- name: NewSwipe :exec
INSERT INTO swipes (swiper_id, swipee_id) VALUES(
    $1, $2 -- Swiper ID and Swipee ID
) ON CONFLICT (swiper_id, swipee_id) DO NOTHING
`

type NewSwipeParams struct {
	SwiperID pgtype.UUID `json:"swiper_id"`
	SwipeeID pgtype.UUID `json:"swipee_id"`
}

// Handle swipes
func (q *Queries) NewSwipe(ctx context.Context, arg NewSwipeParams) error {
	_, err := q.db.Exec(ctx, newSwipe, arg.SwiperID, arg.SwipeeID)
	return err
}
