CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" integer NOT NULL,
  "email" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_bloked"  boolean NOT NULL DEFAULT false,
  "expired_at" timestamptz NOT NULL,
  "created_at" timestamptz DEFAULT 'now'
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
-- ALTER TABLE "users" ADD FOREIGN KEY ("user_type_id") REFERENCES "user_types" ("id");