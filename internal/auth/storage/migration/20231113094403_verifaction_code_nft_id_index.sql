-- migrate:up
CREATE INDEX payload_nft_id_idx ON auth.verification_code((payload->>'nft_id'));


-- migrate:down
DROP INDEX auth.payload_nft_id_idx;

