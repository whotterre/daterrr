-- name: NewSwipe :exec
-- Handle swipes 
INSERT INTO swipes (swiper_id, swipee_id) VALUES(
    $1, $2 -- Swiper ID and Swipee ID
) ON CONFLICT (swiper_id, swipee_id) DO NOTHING;

-- name: CheckMutualSwipe :one
SELECT EXISTS (
    SELECT 1 FROM swipes 
    WHERE 
        swiper_id = $1 AND  -- The other user
        swipee_id = $2      -- Current user
) AS is_mutual;

-- name: CreateMatch :one
WITH new_match AS (
  INSERT INTO matches (user1_id, user2_id)
  VALUES (
    LEAST($1::uuid, $2::uuid),
    GREATEST($1::uuid, $2::uuid)
  )
  ON CONFLICT (user1_id, user2_id) DO NOTHING
  RETURNING *
)
INSERT INTO chats (user1_id, user2_id)
SELECT user1_id, user2_id FROM new_match
ON CONFLICT (user1_id, user2_id) DO NOTHING
RETURNING (SELECT id FROM new_match);