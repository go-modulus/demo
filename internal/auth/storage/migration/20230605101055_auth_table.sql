-- migrate:up

CREATE SCHEMA IF NOT EXISTS "auth";

CREATE TYPE "auth".password_status
    AS ENUM (
    'active', -- The password is active and can be used for authentication.
    'old' -- The password has been changed before and is not actual for today.
    );

CREATE TYPE "auth".password_reset_status
    AS ENUM (
    'active', -- The password reset token is active and can be used.
    'used' -- The password reset token is used.
    );
CREATE TYPE "auth".identity_type
    AS ENUM (
    'email',
    'phone',
    'username'
    );

CREATE TYPE "auth".identity_status
    AS ENUM (
    'not_verified',
    'verified',
    'blocked'
    );

create table auth.identity
(
    id         uuid                 NOT NULL
        primary key,
    user_id    uuid                 NOT NULL,
    identity   text                 NOT NULL
        unique check (identity = trim(lower(identity))),
    "type"     auth.identity_type   NOT NULL,
    status     auth.identity_status NOT NULL
        default 'not_verified'::auth.identity_status,
    created_at timestamptz          NOT NULL
        DEFAULT NOW(),
    updated_at timestamptz          NOT NULL
        DEFAULT NOW()
);
create index identity_user_id_idx
    on auth.identity (user_id);

create table auth.password
(
    id            uuid         NOT NULL
        primary key,
    user_id       uuid         NOT NULL,
    password_hash varchar(255) NOT NULL,
    status        auth.password_status
        default 'active'::auth.password_status NOT NULL,
    created_at    timestamptz  NOT NULL
        DEFAULT NOW(),
    updated_at    timestamptz  NOT NULL
        DEFAULT NOW()
);

create index password_user_id_idx
    on auth.password (user_id);

create table auth.password_reset
(
    id         uuid        NOT NULL
        primary key,
    user_id    uuid        NOT NULL,
    token      text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    status     auth.password_reset_status
                                    default 'active'::auth.password_reset_status NOT NULL
);

create index password_reset_user_id_idx
    on auth.password_reset (user_id);

create table auth.social_auth
(
    id             uuid         NOT NULL
        primary key,
    user_id        uuid,
    external_id    varchar(255) NOT NULL unique,
    external_type  varchar(25)  NOT NULL,
    email          varchar(255) NOT NULL,
    email_verified boolean
        default false NOT NULL,
    first_name     varchar(255) NOT NULL,
    last_name      varchar(255) NOT NULL,
    picture        varchar(512) NOT NULL,
    created_at     timestamptz  NOT NULL
        DEFAULT NOW(),
    updated_at     timestamptz  NOT NULL
        DEFAULT NOW()
);

create index social_auth_user_id_idx
    on auth.social_auth (user_id);

-- migrate:down

DROP table auth.identity;
DROP table auth.password;
DROP TABLE auth.password_reset;
DROP table auth.social_auth;

drop type auth.password_status;
drop type auth.password_reset_status;
drop type auth.identity_type;
drop type auth.identity_status;


DROP SCHEMA auth CASCADE;