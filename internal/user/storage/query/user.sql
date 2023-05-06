set search_path="user";
-- name: CreateUser :one
insert into "user"."user" (id, name, email, registered_at, settings, contacts)
values ($1, $2, $3, now(), null, null) RETURNING *;

-- name: DeleteUser :exec
delete from "user"."user" where id = @id::uuid;

-- name: DeleteUserByEmail :exec
delete from "user"."user" where email = @email::text;




-- name: GetUser :one
select * from "user"."user" where id = @id::uuid LIMIT 1;

-- name: GetNewerUsers :many
select * from "user"."user" order by "user".registered_at DESC LIMIT @count;

-- name: GetUsersFirstPage :many
select * from "user"."user" order by "user".registered_at DESC, "user".id DESC LIMIT @count;

-- name: GetUsersAfterCursor :many
select * from "user"."user"
where
        registered_at < @registered_at OR
    (registered_at = @registered_at AND id < @id::uuid)
order by "user".registered_at DESC, "user".id DESC LIMIT @count;

-- name: GetUsersByIds :many
select * from "user"."user" WHERE id = ANY (@ids::uuid[]);

