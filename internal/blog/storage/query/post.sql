-- name: CreatePost :one
insert into blog."post" (id, title, body, author_id, slug)
values ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeletePost :exec
delete
from blog."post"
where id = @id::uuid;

-- name: GetPost :one
select *
from blog."post"
where id = @id::uuid
LIMIT 1;

-- name: ListPosts :many
select p.*
from blog."post" as p
order by p.published_at desc
limit @count offset @after;

-- name: CountPosts :one
select count(*) as count
from blog."post" as p;

-- name: PublishPost :one
update blog."post"
set published_at = now(),
status = 'published'
where id = @id::uuid
RETURNING *;