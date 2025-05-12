-- name: CreateChat :one
INSERT INTO chats (user1_id, user2_id)
VALUES (
  LEAST($1::uuid, $2::uuid),
  GREATEST($1::uuid, $2::uuid)
)
ON CONFLICT (user1_id, user2_id) DO NOTHING
RETURNING *;

-- name: GetChatByUsers :one
SELECT * FROM chats 
WHERE 
  user1_id = LEAST($1::uuid, $2::uuid) AND
  user2_id = GREATEST($1::uuid, $2::uuid);

-- name: GetUserChats :many
SELECT * FROM chats
WHERE user1_id = $1 OR user2_id = $1
ORDER BY created_at DESC;

-- name: CreateMessage :one
INSERT INTO messages (chat_id, sender_id, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetChatMessages :many
SELECT * FROM messages
WHERE chat_id = $1
ORDER BY created_at ASC;