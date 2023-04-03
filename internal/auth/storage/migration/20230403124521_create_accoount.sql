-- migrate:up
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.local_account
(
    user_id                    uuid primary key,
    email                 text             DEFAULT NULL
        CONSTRAINT local_account_email_uniq UNIQUE,
    nickname                 text             DEFAULT NULL
        CONSTRAINT local_account_nickname_uniq UNIQUE,
    phone                 text             DEFAULT NULL
        CONSTRAINT local_account_phone_uniq UNIQUE,
    "password"            text             DEFAULT NULL,
    created_at         timestamptz      DEFAULT NULL
);

CREATE TABLE auth.account
(
    id                    uuid PRIMARY KEY,
    email                 text             DEFAULT NULL
        CONSTRAINT account_email_uniq UNIQUE,
    "password"            text             DEFAULT NULL,
    confirm_selector      text             DEFAULT NULL
        CONSTRAINT account_confirm_selector_uniq UNIQUE,
    confirm_verifier      text             DEFAULT NULL,
    confirmed             boolean NOT NULL DEFAULT false,
    attempt_count         int     NOT NULL DEFAULT 0,
    last_attempt_at       timestamptz      DEFAULT NULL,
    locked_at             timestamptz      DEFAULT NULL,
    recover_selector      text             DEFAULT NULL
        CONSTRAINT account_recover_selector_uniq UNIQUE,
    recover_verifier      text             DEFAULT NULL,
    recover_token_expiry  timestamptz      DEFAULT NULL,
    oauth2_uid            text             DEFAULT NULL,
    oauth2_provider       text             DEFAULT NULL,
    oauth2_access_token   text             DEFAULT NULL,
    oauth2_refresh_token  text             DEFAULT NULL,
    oauth2_expiry         timestamptz      DEFAULT NULL,
    totp_secret_key       text             DEFAULT NULL,
    sms_phone_number      text             DEFAULT NULL,
    sms_seed_phone_number text             DEFAULT NULL,
    recovery_codes        text             DEFAULT NULL,
    CONSTRAINT account_oauth_uniq UNIQUE (oauth2_uid, oauth2_provider)
);

CREATE TABLE auth.remember_token
(
    id         uuid PRIMARY KEY,
    account_id uuid NOT NULL,
    token      text NOT NULL
);

CREATE INDEX remember_token_account_id ON auth.remember_token(account_id);

-- migrate:down

DROP TABLE auth.remember_token;
DROP TABLE auth.account;
DROP TABLE auth.local_account;
DROP SCHEMA auth;