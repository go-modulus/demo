-- migrate:up
CREATE schema IF NOT EXISTS "blog";

CREATE TYPE blog."post_status" AS ENUM ('draft', 'published');

CREATE TABLE IF NOT EXISTS blog."post" (
  "id" uuid PRIMARY KEY,
  "title" TEXT NOT NULL,
  "body" TEXT NOT NULL,
  "author_id" uuid NOT NULL,
  "slug" TEXT NOT NULL,
  "status" blog.post_status NOT NULL DEFAULT 'draft',
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "published_at" TIMESTAMP with time zone DEFAULT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT NOW()
);

-- migrate:down
DROP TABLE blog."post";
DROP SCHEMA "blog";
