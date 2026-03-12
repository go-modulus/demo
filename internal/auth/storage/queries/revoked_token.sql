-- name: CreateRevokedToken :one
INSERT INTO auth.revoked_token (token_jti, expired)
VALUES ($1, $2)
    ON CONFLICT (token_jti) DO UPDATE
                               SET
                                   expired = EXCLUDED.expired
                               RETURNING *;

-- name: DeleteRevokedToken :exec
DELETE
FROM auth.revoked_token
WHERE token_jti = $1;

-- name: DeleteExpiredRevokedTokens :exec
DELETE
FROM auth.revoked_token
WHERE NOW() > expired;


-- name: GetRevokedTokenByJti :one
SELECT *
FROM auth.revoked_token
WHERE token_jti = $1;