

-- Get complete user profile by ID
-- name: GetUserProfile :one
SELECT 
  u.id,
  u.email,
  u.created_at,
  u.last_active,
  p.first_name,
  p.last_name,
  p.bio,
  p.gender,
  p.age,
  p.image_url,
  p.interests,
  p.location
FROM users u
JOIN profiles p ON u.id = p.user_id
WHERE u.id = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT 
  id, 
  email, 
  created_at,
  last_active
FROM users 
WHERE id = $1
LIMIT 1;

-- name: UpdateProfile :one
UPDATE profiles
SET 
  first_name = COALESCE($2, first_name),
  last_name = COALESCE($3, last_name),
  bio = COALESCE($4, bio),
  image_url = COALESCE($5, image_url),
  interests = COALESCE($6, interests),
  location = COALESCE(ST_SetSRID(ST_MakePoint($7, $8), 4326), location)
WHERE user_id = $1
RETURNING *;

-- name: UpdateLastActive :exec
UPDATE users
SET last_active = now()
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password FROM users
WHERE email = $1
LIMIT 1;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- Delete user and profile with confirmation
-- name: DeleteUserWithProfile :one
WITH deleted_user AS (
  DELETE FROM users WHERE users.id = $1
  RETURNING id, email, created_at
),
deleted_profile AS (
  DELETE FROM profiles WHERE user_id = $1
  RETURNING *
)
SELECT 
  du.id AS user_id,
  du.email,
  du.created_at AS user_created_at,
  dp.first_name,
  dp.last_name,
  dp.gender,
  dp.age,
  dp.image_url
FROM deleted_user du
LEFT JOIN deleted_profile dp ON true;