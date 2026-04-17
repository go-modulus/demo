-- name: SaveUserInfo :one
INSERT INTO "auth"."user_info" (id, name)
VALUES (@id::uuid, @name::text)
RETURNING *;

-- name: DeleteUserInfo :exec
DELETE FROM "auth".user_info WHERE id = @id::uuid;
