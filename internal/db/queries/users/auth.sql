-----------------------------------------
-- 1. LOGIN FUNCTIONALITY
-----------------------------------------

-- For login page: Get user by email + password hash
-- name: GetUserForLogin :one
SELECT id, email, password FROM users 
WHERE email = $1 
LIMIT 1;

-- After successful login: Create session
-- name: CreateSession :one
INSERT INTO user_sessions (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-----------------------------------------
-- 2. SESSION MANAGEMENT
-----------------------------------------

-- Middleware: Check if token is valid
-- name: GetSessionByToken :one
SELECT 
  s.*,
  u.email  -- Include user info
FROM user_sessions s
JOIN users u ON s.user_id = u.id
WHERE s.token = $1 
  AND s.expires_at > now()
  AND s.is_revoked = false;

-- Logout: Delete session
-- name: DeleteSession :exec
DELETE FROM user_sessions 
WHERE token = $1;

-----------------------------------------
-- 3. PASSWORD RESET
-----------------------------------------

-- Generate password reset token
-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- Validate reset token (24-hour expiry)
-- name: GetValidPasswordResetToken :one
SELECT * FROM password_reset_tokens
WHERE token = $1 
  AND expires_at > now()
  AND used = false
LIMIT 1;

-- After password update: Mark token as used
-- name: MarkResetTokenUsed :exec
UPDATE password_reset_tokens 
SET used = true 
WHERE token = $1;