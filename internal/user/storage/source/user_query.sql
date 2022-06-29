-- name: GetUser :one
select * from "user"."user" where id = $1 LIMIT 1;

-- name: GetNewerUsers :many
select * from "user"."user" order by "user".registered_at DESC LIMIT $1;

-- name: GetUsersFirstPage :many
select * from "user"."user" order by "user".registered_at DESC, "user".id DESC LIMIT $1;

-- name: GetUsersAfterCursor :many
select * from "user"."user"
where
    registered_at < $1 OR
    (registered_at = $2 AND id < $3)
order by "user".registered_at DESC, "user".id DESC LIMIT $4;

-- name: GetUsersByIds :many
select * from "user"."user" WHERE id = ANY ($1::uuid[]);
