-- migrate:up
CREATE INDEX payload_dana_user_idx ON auth.verification_code((payload->>'dana_user_id'));


-- migrate:down
DROP INDEX auth.payload_dana_user_idx;

