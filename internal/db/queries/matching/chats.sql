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

-- name: GetUserChatsWithChatID :many
SELECT id AS chat_id, user1_id, user2_id, created_at
FROM chats
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

-- name: GetConversationsForUser :many
SELECT 
    c.id AS chat_id,
    c.created_at,
    u1.id AS user1_id,
    p1.first_name AS user1_first_name,
    p1.last_name AS user1_last_name,
    p1.image_url AS user1_profile_picture,
    u2.id AS user2_id,
    p2.first_name AS user2_first_name,
    p2.last_name AS user2_last_name,
    p2.image_url AS user2_profile_picture,
    m.id AS last_message_id,
    m.chat_id AS last_message_chat_id,
    m.sender_id AS last_message_sender_id,
    m.content AS last_message_content,
    m.created_at AS last_message_created_at,
    COUNT(m.id) FILTER (WHERE m.read_at IS NULL AND m.sender_id != $1) AS unread_count
FROM chats c
JOIN users u1 ON c.user1_id = u1.id
JOIN users u2 ON c.user2_id = u2.id
JOIN profiles p1 ON u1.id = p1.user_id
JOIN profiles p2 ON u2.id = p2.user_id
LEFT JOIN messages m ON c.id = m.chat_id
WHERE c.user1_id = $1 OR c.user2_id = $1
GROUP BY 
    c.id, 
    c.created_at, 
    u1.id, 
    p1.first_name, 
    p1.last_name, 
    p1.image_url, 
    u2.id, 
    p2.first_name, 
    p2.last_name, 
    p2.image_url, 
    m.id, 
    m.chat_id, 
    m.sender_id, 
    m.content, 
    m.created_at
ORDER BY c.created_at DESC;
