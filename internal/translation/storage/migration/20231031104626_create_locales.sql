-- migrate:up
CREATE SCHEMA "translation";
CREATE TYPE "translation"."locale" AS ENUM ('en', 'id');
CREATE TYPE "translation"."path" AS ENUM (
    'blog.post.title',
    'blog.post.body'
    );

-- migrate:down
DROP TYPE "translation"."locale";
DROP TYPE "translation"."path";
DROP SCHEMA "translation" CASCADE;

