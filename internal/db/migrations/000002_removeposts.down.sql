-- 1. Remove automatic match trigger first
DROP TRIGGER IF EXISTS trg_create_match ON swipes;
DROP FUNCTION IF EXISTS create_match_on_mutual_swipe();

-- 2. Drop helper views
DROP VIEW IF EXISTS user_matches;
DROP VIEW IF EXISTS potential_matches;

-- 3. Recreate original swipes table with 'liked' column
CREATE TABLE swipes_original (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  swiper_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  swipee_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  swiped_at timestamp DEFAULT now(),
  UNIQUE (swiper_id, swipee_id)
);

-- 4. Migrate data (only right swipes exist in new system)
INSERT INTO swipes_original (id, swiper_id, swipee_id, liked, swiped_at)
SELECT id, swiper_id, swipee_id, true, swiped_at FROM swipes;

-- 5. Replace the table
DROP TABLE swipes;
ALTER TABLE swipes_original RENAME TO swipes;

-- 6. Recreate indexes
CREATE INDEX ON swipes (swiper_id);
CREATE INDEX ON swipes (swipee_id);