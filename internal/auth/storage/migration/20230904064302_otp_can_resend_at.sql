-- migrate:up

ALTER TABLE auth.one_time_password ADD COLUMN can_resend_at timestamptz NULL;
UPDATE auth.one_time_password SET can_resend_at = created_at + interval '1 minutes';
ALTER TABLE auth.one_time_password ALTER COLUMN can_resend_at SET NOT NULL;

ALTER TABLE auth.verification_code ADD COLUMN can_resend_at timestamptz NULL;
UPDATE auth.verification_code SET can_resend_at = created_at + interval '1 minutes';
ALTER TABLE auth.verification_code ALTER COLUMN can_resend_at SET NOT NULL;

-- migrate:down

ALTER TABLE auth.one_time_password DROP COLUMN can_resend_at;

ALTER TABLE auth.verification_code DROP COLUMN can_resend_at;