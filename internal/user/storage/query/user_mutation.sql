set search_path="user";
-- name: CreateUser :one
insert into "user"."user" (id, name, email, registered_at, settings, contacts)
values ($1, $2, $3, now(), null, null) RETURNING *;

-- name: DeleteUser :exec
delete from "user"."user" where id = @id::uuid;
