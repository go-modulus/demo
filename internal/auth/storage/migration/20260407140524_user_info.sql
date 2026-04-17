-- migrate:up

CREATE TABLE "auth"."user_info"
(
    -- the same as auth.account.id
    id         uuid PRIMARY KEY,
    name       text      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- migrate:down
DROP TABLE "auth"."user_info";
