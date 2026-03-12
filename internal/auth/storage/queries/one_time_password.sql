-- name: CreateOneTimePassword :one
INSERT INTO auth.one_time_password (token, email, user_id, used_at, expires_at, can_resend_at)
VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (token) DO UPDATE
                               SET
                               user_id = EXCLUDED.user_id,
                               email = EXCLUDED.email,
                               used_at = EXCLUDED.used_at,
                               expires_at = EXCLUDED.expires_at,
                               can_resend_at = EXCLUDED.can_resend_at
                               RETURNING *;

-- name: GetOneTimePasswordByToken :one
SELECT *
FROM auth.one_time_password
WHERE token = $1;

-- name: GetLastOneTimePassword :one
SELECT *
FROM auth.one_time_password
ORDER BY created_at desc
LIMIT 1;

-- name: GetLastActiveOneTimePasswordByEmail :one
SELECT *
FROM auth.one_time_password
WHERE
    email = $1
AND
    used_at is null
AND
    expires_at > NOW()
AND
    can_resend_at > NOW()
ORDER BY created_at desc
    LIMIT 1;

-- name: SetUsedOneTimePassword :one
UPDATE "auth"."one_time_password"
SET user_id = @user_id::uuid,
used_at = NOW()
WHERE token = @token
    RETURNING *;

-- name: DeleteOneTimePasswordByEmail :exec
DELETE FROM "auth"."one_time_password"
WHERE email = $1;

-- name: UpdateOneTimePasswordExpiresAt :one
UPDATE "auth"."one_time_password"
SET expires_at = @expires_at::timestamptz
WHERE token = @token
    RETURNING *;

-- name: DeleteOneTimePasswordByToken :exec
DELETE
FROM auth.one_time_password
WHERE token = $1;
