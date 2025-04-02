-- name: NewSwipe :exec
-- Handle swipes 
INSERT INTO swipes (swiper_id, swipee_id) VALUES(
    $1, $2 -- Swiper ID and Swipee ID
) ON CONFLICT (swiper_id, swipee_id) DO NOTHING;