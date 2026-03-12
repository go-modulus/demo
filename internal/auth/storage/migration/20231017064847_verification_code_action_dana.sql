-- migrate:up
ALTER TYPE auth.verification_action ADD VALUE IF NOT EXISTS 'confirm_dana_user';

-- migrate:down

