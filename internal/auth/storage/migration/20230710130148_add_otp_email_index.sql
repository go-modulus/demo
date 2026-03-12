-- migrate:up
CREATE INDEX auth_one_time_password_email_idx
    ON auth.one_time_password (email);

-- migrate:down
DROP INDEX auth.auth_one_time_password_email_idx;

