// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: auth.sql

package db

import (
	"context"
	"net/netip"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPasswordResetToken = `-- name: CreatePasswordResetToken :one

INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, token_hash, created_at, expires_at, used
`

type CreatePasswordResetTokenParams struct {
	UserID    pgtype.UUID      `json:"user_id"`
	TokenHash string           `json:"token_hash"`
	ExpiresAt pgtype.Timestamp `json:"expires_at"`
}

// ---------------------------------------
// 3. PASSWORD RESET
// ---------------------------------------
// Generate password reset token
func (q *Queries) CreatePasswordResetToken(ctx context.Context, arg CreatePasswordResetTokenParams) (PasswordResetToken, error) {
	row := q.db.QueryRow(ctx, createPasswordResetToken, arg.UserID, arg.TokenHash, arg.ExpiresAt)
	var i PasswordResetToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TokenHash,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.Used,
	)
	return i, err
}

const createSession = `-- name: CreateSession :one
INSERT INTO user_sessions (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, token, created_at, expires_at, ip_address, user_agent, is_revoked
`

type CreateSessionParams struct {
	UserID    pgtype.UUID      `json:"user_id"`
	Token     string           `json:"token"`
	ExpiresAt pgtype.Timestamp `json:"expires_at"`
}

// After successful login: Create session
func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (UserSession, error) {
	row := q.db.QueryRow(ctx, createSession, arg.UserID, arg.Token, arg.ExpiresAt)
	var i UserSession
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IpAddress,
		&i.UserAgent,
		&i.IsRevoked,
	)
	return i, err
}

const deleteSession = `-- name: DeleteSession :exec
DELETE FROM user_sessions 
WHERE token = $1
`

// Logout: Delete session
func (q *Queries) DeleteSession(ctx context.Context, token string) error {
	_, err := q.db.Exec(ctx, deleteSession, token)
	return err
}

const getSessionByToken = `-- name: GetSessionByToken :one

SELECT 
  s.id, s.user_id, s.token, s.created_at, s.expires_at, s.ip_address, s.user_agent, s.is_revoked,
  u.email  -- Include user info
FROM user_sessions s
JOIN users u ON s.user_id = u.id
WHERE s.token = $1 
  AND s.expires_at > now()
  AND s.is_revoked = false
`

type GetSessionByTokenRow struct {
	ID        pgtype.UUID      `json:"id"`
	UserID    pgtype.UUID      `json:"user_id"`
	Token     string           `json:"token"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	ExpiresAt pgtype.Timestamp `json:"expires_at"`
	IpAddress *netip.Addr      `json:"ip_address"`
	UserAgent pgtype.Text      `json:"user_agent"`
	IsRevoked pgtype.Bool      `json:"is_revoked"`
	Email     string           `json:"email"`
}

// ---------------------------------------
// 2. SESSION MANAGEMENT
// ---------------------------------------
// Middleware: Check if token is valid
func (q *Queries) GetSessionByToken(ctx context.Context, token string) (GetSessionByTokenRow, error) {
	row := q.db.QueryRow(ctx, getSessionByToken, token)
	var i GetSessionByTokenRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IpAddress,
		&i.UserAgent,
		&i.IsRevoked,
		&i.Email,
	)
	return i, err
}

const getUserForLogin = `-- name: GetUserForLogin :one

SELECT id, email, password FROM users 
WHERE email = $1 
LIMIT 1
`

type GetUserForLoginRow struct {
	ID       pgtype.UUID `json:"id"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}

// ---------------------------------------
// 1. LOGIN FUNCTIONALITY
// ---------------------------------------
// For login page: Get user by email + password hash
func (q *Queries) GetUserForLogin(ctx context.Context, email string) (GetUserForLoginRow, error) {
	row := q.db.QueryRow(ctx, getUserForLogin, email)
	var i GetUserForLoginRow
	err := row.Scan(&i.ID, &i.Email, &i.Password)
	return i, err
}

const getValidPasswordResetToken = `-- name: GetValidPasswordResetToken :one

SELECT id, user_id, token_hash, created_at, expires_at, used FROM password_reset_tokens
WHERE token_hash = $1 
  AND expires_at > now()
  AND used = false
LIMIT 1
`

// djd
// Validate reset token (24-hour expiry)
func (q *Queries) GetValidPasswordResetToken(ctx context.Context, tokenHash string) (PasswordResetToken, error) {
	row := q.db.QueryRow(ctx, getValidPasswordResetToken, tokenHash)
	var i PasswordResetToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TokenHash,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.Used,
	)
	return i, err
}

const markResetTokenUsed = `-- name: MarkResetTokenUsed :exec
UPDATE password_reset_tokens 
SET used = true 
WHERE token_hash = $1
`

// After password update: Mark token as used
func (q *Queries) MarkResetTokenUsed(ctx context.Context, tokenHash string) error {
	_, err := q.db.Exec(ctx, markResetTokenUsed, tokenHash)
	return err
}
