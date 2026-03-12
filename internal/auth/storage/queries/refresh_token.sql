-- name: CreateRefreshToken :one
INSERT INTO auth.refresh_token (hash, user_id, session_id, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteRefreshToken :exec
DELETE
FROM auth.refresh_token
WHERE hash = $1;

-- name: DeleteRefreshTokensBySessionId :exec
DELETE
FROM auth.refresh_token
WHERE session_id = $1;


-- name: DeleteRefreshTokensByUserId :exec
DELETE
FROM auth.refresh_token
WHERE user_id = $1;


-- name: GetRefreshTokenByHash :one
SELECT *
FROM auth.refresh_token
WHERE hash = $1
and status = 'active';

-- name: RevokeRefreshTokenBySessionId :one
UPDATE "auth"."refresh_token"
SET status = 'revoked'
WHERE session_id = $1
    RETURNING *;

-- name: UpdateRefreshTokenExpiresAt :one
UPDATE "auth"."refresh_token"
SET expires_at = $1
WHERE hash = $2
    RETURNING *;

-- name: UpdateRefreshTokenUserId :one
UPDATE "auth"."refresh_token"
SET user_id = $1
WHERE hash = $2
    RETURNING *;
