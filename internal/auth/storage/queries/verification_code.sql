-- name: CreateVerificationCode :one
INSERT INTO auth.verification_code (code, action, email, user_id, used_at, payload, expires_at, can_resend_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    ON CONFLICT (code) DO UPDATE
                               SET
                              action = EXCLUDED.action,
                              email = EXCLUDED.email,
                              user_id = EXCLUDED.user_id,
                              used_at = EXCLUDED.used_at,
                              payload = EXCLUDED.payload,
                              expires_at = EXCLUDED.expires_at,
                              can_resend_at = EXCLUDED.can_resend_at
                               RETURNING *;

-- name: GetVerificationCodeByCodeAndAction :one
SELECT *
FROM auth.verification_code
WHERE
    code = $1
AND
    action = $2;

-- name: DeleteVerificationCode :exec
DELETE
FROM
    auth.verification_code
WHERE
    code = $1;

-- name: GetLastVerificationCodeByAction :one
SELECT *
FROM auth.verification_code
WHERE action = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: GetLastActiveVerificationCodeByUserId :one
SELECT *
FROM auth.verification_code
WHERE
    user_id = $1
AND
    used_at is null
AND
    expires_at > NOW()
AND
    can_resend_at > NOW()
ORDER BY created_at desc
LIMIT 1;

-- name: GetLastActiveVerificationCodeForTransferNft :one
SELECT *
FROM auth.verification_code
WHERE
    user_id = $1
  AND
    used_at is null
  AND
    expires_at > NOW()
  AND
    can_resend_at > NOW()
  AND
    payload ->> 'nft_id' = @nft_id::text
  AND
    payload ->> 'wallet_id' = @wallet_id::text
LIMIT 1;

-- name: GetLastActiveVerificationCodeForTransferNftByAddress :one
SELECT *
FROM auth.verification_code
WHERE
        user_id = $1
  AND
    used_at is null
  AND
        expires_at > NOW()
  AND
        can_resend_at > NOW()
  AND
        payload ->> 'nft_id' = @nft_id::text
  AND
    payload ->> 'to_address' = @to_address::text
    LIMIT 1;

-- name: DeleteVerificationCodeForTransferNft :exec
DELETE
FROM auth.verification_code
WHERE
        user_id = $1
  AND
        payload ->> 'nft_id' = @nft_id::text
  AND
        payload ->> 'wallet_id' = @wallet_id::text;

-- name: DeleteVerificationCodeForTransferNftByAddress :exec
DELETE
FROM auth.verification_code
WHERE
        user_id = $1
  AND
        payload ->> 'nft_id' = @nft_id::text
  AND
    payload ->> 'to_address' = @to_address::text;

-- name: DeleteVerificationCodeForConfirmDanaUser :exec
DELETE
FROM auth.verification_code
WHERE
    payload ->> 'dana_user_id' = @dana_user_id::text;

-- name: GetVerificationCodeByDanaUserId :one
SELECT *
FROM auth.verification_code
WHERE
    payload ->> 'dana_user_id' = @dana_user_id::text
ORDER BY created_at desc
LIMIT 1;

-- name: UpdateVerificationCodeExpires :one
UPDATE "auth"."verification_code"
SET expires_at = @expires_at::timestamptz
WHERE code = @code
    RETURNING *;

-- name: SetUsedVerificationCode :one
UPDATE "auth"."verification_code"
SET user_id = @user_id::uuid,
used_at = NOW()
WHERE code = @code
    RETURNING *;