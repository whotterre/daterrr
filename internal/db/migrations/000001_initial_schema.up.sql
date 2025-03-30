CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" text UNIQUE NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "last_active" timestamp DEFAULT (now())
);

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

CREATE TABLE "swipes" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "swiper_id" uuid NOT NULL,
  "swipee_id" uuid NOT NULL,
  "liked" boolean NOT NULL,
  "swiped_at" timestamp DEFAULT (now())
);

CREATE TABLE "matches" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user1_id" uuid NOT NULL,
  "user2_id" uuid NOT NULL,
  "matched_at" timestamp DEFAULT (now())
);

CREATE TABLE "chats" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user1_id" uuid NOT NULL,
  "user2_id" uuid NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "messages" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "chat_id" uuid NOT NULL,
  "sender_id" uuid NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "read_at" timestamp
);

CREATE TABLE "posts" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid NOT NULL,
  "content" text,
  "image_url" text,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp
);

CREATE TABLE "comments" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "post_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

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


CREATE UNIQUE INDEX ON "users" ("email");

CREATE UNIQUE INDEX ON "profiles" ("user_id");

CREATE INDEX ON "profiles" USING GIST ("location");

CREATE INDEX ON "profiles" USING GIN ("interests");

CREATE UNIQUE INDEX ON "swipes" ("swiper_id", "swipee_id");

CREATE INDEX ON "messages" ("chat_id", "created_at");

CREATE INDEX ON "posts" ("user_id", "created_at");

CREATE INDEX ON "comments" ("post_id", "created_at");

-- Optimize auth queries
CREATE INDEX idx_session_token ON user_sessions(token_hash);
CREATE INDEX idx_session_user ON user_sessions(user_id) WHERE NOT is_revoked;
CREATE INDEX idx_reset_tokens ON password_reset_tokens(token_hash) WHERE NOT consumed;
CREATE INDEX idx_login_attempts ON login_attempts(user_id, ip_address);

COMMENT ON COLUMN "users"."email" IS 'Validated with regex';

COMMENT ON COLUMN "users"."password" IS 'Min length 8, hashed';

COMMENT ON COLUMN "profiles"."first_name" IS 'Min length 2';

COMMENT ON COLUMN "profiles"."last_name" IS 'Min length 2';

COMMENT ON COLUMN "profiles"."bio" IS 'Max length 500';

COMMENT ON COLUMN "profiles"."gender" IS 'male|female|non-binary|other';

COMMENT ON COLUMN "profiles"."age" IS '18-120';

COMMENT ON COLUMN "messages"."content" IS 'Max length 2000';

COMMENT ON COLUMN "posts"."content" IS 'Max length 2000';

COMMENT ON COLUMN "comments"."content" IS 'Max length 500';

ALTER TABLE "profiles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "swipes" ADD FOREIGN KEY ("swiper_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "swipes" ADD FOREIGN KEY ("swipee_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "matches" ADD FOREIGN KEY ("user1_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "matches" ADD FOREIGN KEY ("user2_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "chats" ADD FOREIGN KEY ("user1_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "chats" ADD FOREIGN KEY ("user2_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "messages" ADD FOREIGN KEY ("chat_id") REFERENCES "chats" ("id") ON DELETE CASCADE;

ALTER TABLE "messages" ADD FOREIGN KEY ("sender_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE;

ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;