-- migrate:up
ALTER TABLE auth.one_time_password ADD COLUMN email varchar(255) not null default '';

ALTER TABLE auth.one_time_password ALTER COLUMN email drop default;

-- migrate:down
ALTER TABLE auth.one_time_password DROP COLUMN email;