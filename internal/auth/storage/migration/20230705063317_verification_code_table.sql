-- migrate:up
CREATE TYPE "auth".verification_action
    AS ENUM (
    'transfer_nft'
    );

CREATE TABLE auth.verification_code (
    code text PRIMARY KEY,
    action auth.verification_action NOT NULL,
    email text NOT NULL,
    user_id uuid NULL,
    used_at timestamptz NULL,
    payload jsonb NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL
);

-- migrate:down

DROP TABLE auth.verification_code;
DROP TYPE auth.verification_action;