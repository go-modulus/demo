-- migrate:up
CREATE SCHEMA IF NOT EXISTS "user";

CREATE TABLE "user"."user"
(
    id            uuid                     NOT NULL
        CONSTRAINT user_pk
            primary key,
    name          varchar(50)              NOT NULL,
    email         varchar(127)             NOT NULL,
    registered_at timestamp with time zone NOT NULL,
    settings      jsonb,
    contacts      text[]
);



-- migrate:down
DROP TABLE "user"."user";
DROP SCHEMA "user";
