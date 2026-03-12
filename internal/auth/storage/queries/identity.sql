-- name: CreateIdentity :one
insert into "auth"."identity"
    (id, user_id, identity, type)
values (@id::uuid, @user_id::uuid, trim(lower(@identity::text)), @type)
RETURNING *;

-- name: VerifyIdentity :one
UPDATE "auth"."identity"
SET status = 'verified'
WHERE id = @id::uuid
RETURNING *;

-- name: CreatePassword :one
insert into "auth"."password"
    (id, user_id, password_hash)
values (@id::uuid, @user_id::uuid, @password_hash::text)
RETURNING *;

-- name: DeleteIdentity :exec
delete from "auth"."identity"
where id = @id::uuid;

-- name: DeletePassword :exec
delete from "auth"."password"
where id = @id::uuid;

-- name: ChangePasswordStatus :exec
update "auth"."password"
set status = @status
where id = @id::uuid;

-- name: DeleteUserIdentities :exec
DELETE FROM "auth"."identity"
WHERE user_id = @user_id::uuid;

-- name: DeleteUserPasswords :exec
DELETE FROM "auth"."password"
WHERE user_id = @user_id::uuid;

-- name: SelectIdentity :one
select *
from "auth"."identity"
where identity = trim(lower(@identity::text));

-- name: SelectUserPasswords :many
select *
from "auth"."password"
where user_id = @user_id::uuid
order by created_at desc
limit 2;