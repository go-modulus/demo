-- migrate:up
CREATE INDEX auth_verification_code_user_idx
    ON auth.verification_code (user_id);

-- migrate:down
DROP INDEX auth.auth_verification_code_user_idx;
