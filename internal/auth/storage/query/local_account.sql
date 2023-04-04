set search_path="auth";
-- name: CreateLocalAccount :one
INSERT INTO auth.local_account (user_id, email, nickname, phone, password, created_at)
VALUES (@user_id::uuid, @email, @nickname,
        @phone,
        @password_hash::text, NOW())
 RETURNING *;

-- name: DeleteLocalAccount :execrows
delete from auth.local_account where user_id = @user_id::uuid;
