-- Create "users" table
CREATE TABLE "users" ("id" character varying NOT NULL, "password" character varying NOT NULL, "buy_count" bigint NOT NULL DEFAULT 0, "weekly_purchase" boolean NOT NULL DEFAULT false, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
