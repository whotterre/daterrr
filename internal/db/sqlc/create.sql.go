// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: create.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createNewUser = `-- name: CreateNewUser :one
WITH new_user AS (
  INSERT INTO users (email, password)
  VALUES ($1, $2)
  RETURNING id, email, created_at
)
INSERT INTO profiles (user_id, first_name, last_name, bio, gender, age, image_url, location, interests)
VALUES (
  (SELECT id FROM new_user),
  $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING 
  (SELECT email FROM new_user) AS email,
  (SELECT created_at FROM new_user) AS user_created_at,
  profiles.id, profiles.user_id, profiles.first_name, profiles.last_name, profiles.bio, profiles.gender, profiles.age, profiles.image_url, profiles.location, profiles.interests
`

type CreateNewUserParams struct {
	Email     string       `json:"email"`
	Password  string       `json:"password"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Bio       pgtype.Text  `json:"bio"`
	Gender    string       `json:"gender"`
	Age       int32        `json:"age"`
	ImageUrl  pgtype.Text  `json:"image_url"`
	Location  pgtype.Point `json:"location"`
	Interests []string     `json:"interests"`
}

type CreateNewUserRow struct {
	Email         string           `json:"email"`
	UserCreatedAt pgtype.Timestamp `json:"user_created_at"`
	ID            pgtype.UUID      `json:"id"`
	UserID        pgtype.UUID      `json:"user_id"`
	FirstName     string           `json:"first_name"`
	LastName      string           `json:"last_name"`
	Bio           pgtype.Text      `json:"bio"`
	Gender        string           `json:"gender"`
	Age           int32            `json:"age"`
	ImageUrl      pgtype.Text      `json:"image_url"`
	Location      pgtype.Point     `json:"location"`
	Interests     []string         `json:"interests"`
}

func (q *Queries) CreateNewUser(ctx context.Context, arg CreateNewUserParams) (CreateNewUserRow, error) {
	row := q.db.QueryRow(ctx, createNewUser,
		arg.Email,
		arg.Password,
		arg.FirstName,
		arg.LastName,
		arg.Bio,
		arg.Gender,
		arg.Age,
		arg.ImageUrl,
		arg.Location,
		arg.Interests,
	)
	var i CreateNewUserRow
	err := row.Scan(
		&i.Email,
		&i.UserCreatedAt,
		&i.ID,
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.Bio,
		&i.Gender,
		&i.Age,
		&i.ImageUrl,
		&i.Location,
		&i.Interests,
	)
	return i, err
}
