-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    type,
    data
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserNotifications :many
SELECT * FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUnreadNotificationsCount :one
SELECT COUNT(*) FROM notifications
WHERE user_id = $1 AND read = FALSE;

-- name: MarkNotificationAsRead :exec
UPDATE notifications
SET read = TRUE
WHERE id = $1 AND user_id = $2;

-- name: MarkAllNotificationsAsRead :exec
UPDATE notifications
SET read = TRUE
WHERE user_id = $1 AND read = FALSE;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = $1 AND user_id = $2;