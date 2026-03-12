-- migrate:up

CREATE TYPE "auth".rt_status
    AS ENUM (
    'active', -- The token can be used to refresh an access token.
    'revoked' -- The token is marked as revoked. It cannot be used to refresh an access token.
    );

CREATE TABLE auth.refresh_token (
                                    hash text PRIMARY KEY,
                                    session_id uuid NOT NULL,
                                    user_id uuid NOT NULL,
                                    status auth.rt_status NOT NULL DEFAULT 'active',
                                    created_at timestamptz NOT NULL DEFAULT now(),
                                    expires_at timestamptz NOT NULL
);

create index refresh_token_user_id_idx
    on auth.refresh_token (user_id);

create index refresh_token_session_id_idx
    on auth.refresh_token (session_id);

-- migrate:down

DROP TABLE auth.refresh_token;
DROP TYPE auth.rt_status;