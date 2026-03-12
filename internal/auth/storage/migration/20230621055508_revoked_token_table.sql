-- migrate:up
CREATE TABLE auth.revoked_token (
    token_jti text PRIMARY KEY,
    expired timestamptz NOT NULL
);


-- migrate:down
DROP TABLE auth.revoked_token;
