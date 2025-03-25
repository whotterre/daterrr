-- Drop foreign keys first to avoid dependency issues
ALTER TABLE "comments" DROP CONSTRAINT "comments_user_id_fkey";
ALTER TABLE "comments" DROP CONSTRAINT "comments_post_id_fkey";
ALTER TABLE "posts" DROP CONSTRAINT "posts_user_id_fkey";
ALTER TABLE "messages" DROP CONSTRAINT "messages_sender_id_fkey";
ALTER TABLE "messages" DROP CONSTRAINT "messages_chat_id_fkey";
ALTER TABLE "chats" DROP CONSTRAINT "chats_user2_id_fkey";
ALTER TABLE "chats" DROP CONSTRAINT "chats_user1_id_fkey";
ALTER TABLE "matches" DROP CONSTRAINT "matches_user2_id_fkey";
ALTER TABLE "matches" DROP CONSTRAINT "matches_user1_id_fkey";
ALTER TABLE "swipes" DROP CONSTRAINT "swipes_swipee_id_fkey";
ALTER TABLE "swipes" DROP CONSTRAINT "swipes_swiper_id_fkey";
ALTER TABLE "profiles" DROP CONSTRAINT "profiles_user_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "swipes_swiper_id_swipee_id_idx";
DROP INDEX IF EXISTS "matches_user1_id_user2_id_idx";
DROP INDEX IF EXISTS "chats_user1_id_user2_id_idx";

-- Drop tables
DROP TABLE IF EXISTS "comments";
DROP TABLE IF EXISTS "posts";
DROP TABLE IF EXISTS "messages";
DROP TABLE IF EXISTS "chats";
DROP TABLE IF EXISTS "matches";
DROP TABLE IF EXISTS "swipes";
DROP TABLE IF EXISTS "profiles";
DROP TABLE IF EXISTS "users";
