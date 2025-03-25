-- name: CreateNewUser :one
WITH new_user AS (
  INSERT INTO users (email, password)
  VALUES ($1, $2)
  RETURNING id, email, created_at
)
INSERT INTO profiles (user_id, first_name, last_name, bio, gender, age, image_url, interests)
VALUES (
  (SELECT id FROM new_user),
  $3, $4, $5, $6, $7, $8, $9
)
RETURNING 
  (SELECT email FROM new_user) AS email,
  (SELECT created_at FROM new_user) AS user_created_at,
  profiles.*;
