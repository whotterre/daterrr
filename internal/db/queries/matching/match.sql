-- name: GetUserMatches :many
SELECT 
    match_data.match_id,
    match_data.match_date,
    match_data.matched_user_id,
    u.email,
    p.first_name,
    p.last_name,
    p.image_url,
    p.bio,
    c.id AS chat_id,
    COALESCE(
        (
            SELECT MAX(created_at)
            FROM messages
            WHERE chat_id = c.id
        ),
        match_data.match_date
    ) AS last_interaction
FROM (
    SELECT
        m.id AS match_id,
        m.matched_at AS match_date,
        CASE
            WHEN m.user1_id = $1 THEN m.user2_id
            ELSE m.user1_id
        END AS matched_user_id,
        m.user1_id,
        m.user2_id
    FROM matches m
    WHERE m.user1_id = $1 OR m.user2_id = $1
) AS match_data
JOIN users u ON u.id = match_data.matched_user_id
JOIN profiles p ON p.user_id = u.id
LEFT JOIN chats c ON (
    (c.user1_id = match_data.user1_id AND c.user2_id = match_data.user2_id) OR
    (c.user1_id = match_data.user2_id AND c.user2_id = match_data.user1_id)
)
ORDER BY last_interaction DESC;
