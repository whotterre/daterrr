-- 1. Users table
CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" text UNIQUE NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "last_active" timestamp DEFAULT (now())
);

-- 2. Profiles table
CREATE TABLE "profiles" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid UNIQUE NOT NULL,
  "first_name" text NOT NULL,
  "last_name" text NOT NULL,
  "bio" text,
  "gender" text NOT NULL,
  "age" integer NOT NULL,
  "image_url" text,
  "location" point,
  "interests" text[]
);

-- 3. Swipes table (only stores right swipes)
CREATE TABLE "swipes" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "swiper_id" uuid NOT NULL,
  "swipee_id" uuid NOT NULL,
  "swiped_at" timestamp DEFAULT (now()),
  UNIQUE ("swiper_id", "swipee_id")
);

-- 4. Matches table with consistent ordering
CREATE TABLE "matches" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user1_id" uuid NOT NULL,  
  "user2_id" uuid NOT NULL,
  "matched_at" timestamp DEFAULT (now()),
  UNIQUE ("user1_id", "user2_id")
);

-- 5. Chats table with same ordering as matches
CREATE TABLE "chats" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user1_id" uuid NOT NULL,
  "user2_id" uuid NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  UNIQUE ("user1_id", "user2_id")
);

-- 6. Messages table
CREATE TABLE "messages" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "chat_id" uuid NOT NULL,
  "sender_id" uuid NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "read_at" timestamp
);

-- 7. Notifications table
CREATE TABLE notifications (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type text NOT NULL,
  data JSONB NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  read BOOLEAN DEFAULT FALSE
);

-- 8. Auth tables
CREATE TABLE user_sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token text NOT NULL UNIQUE,
  created_at timestamp DEFAULT now(),
  expires_at timestamp NOT NULL,
  ip_address inet,
  user_agent text,
  is_revoked boolean DEFAULT false
);

CREATE TABLE password_reset_tokens (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash text NOT NULL UNIQUE,
  created_at timestamp DEFAULT now(),
  expires_at timestamp NOT NULL,
  used boolean DEFAULT false
);

-- ===== INDEXES =====
CREATE UNIQUE INDEX ON "users" ("email");
CREATE UNIQUE INDEX ON "profiles" ("user_id");
CREATE INDEX ON "profiles" USING GIST ("location");
CREATE INDEX ON "profiles" USING GIN ("interests");
CREATE INDEX ON "swipes" ("swiper_id");
CREATE INDEX ON "swipes" ("swipee_id");
CREATE INDEX ON "matches" ("user1_id");
CREATE INDEX ON "matches" ("user2_id");
CREATE INDEX ON "chats" ("user1_id");
CREATE INDEX ON "chats" ("user2_id");
CREATE INDEX ON "messages" ("chat_id", "created_at");
CREATE INDEX idx_session_token ON user_sessions(token);
CREATE INDEX idx_session_user ON user_sessions(user_id) WHERE NOT is_revoked;
CREATE INDEX idx_reset_tokens ON password_reset_tokens(token_hash) WHERE NOT used;

-- ===== CONSTRAINTS =====
ALTER TABLE "profiles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "swipes" ADD FOREIGN KEY ("swiper_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "swipes" ADD FOREIGN KEY ("swipee_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "matches" ADD FOREIGN KEY ("user1_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "matches" ADD FOREIGN KEY ("user2_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "chats" ADD FOREIGN KEY ("user1_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "chats" ADD FOREIGN KEY ("user2_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "messages" ADD FOREIGN KEY ("chat_id") REFERENCES "chats" ("id") ON DELETE CASCADE;
ALTER TABLE "messages" ADD FOREIGN KEY ("sender_id") REFERENCES "users" ("id") ON DELETE CASCADE;