-- migrate:up
CREATE TABLE auth.one_time_password (
    token varchar(255) PRIMARY KEY,
    user_id uuid NULL,
    used_at timestamptz NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL
);

-- migrate:down
DROP TABLE auth.one_time_password
