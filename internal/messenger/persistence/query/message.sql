-- name: GetMessage :one
select *
from "messenger"."message"
where id = $1
limit 1;

-- name: FindLastMessages :many
select m.*
from "messenger"."conversation" c
         join lateral (
    select mj.*
    from "messenger"."message" mj
    where mj.conversation_id = c.id
      and mj.type = 'text'
    order by mj.created_at desc
    limit 1
    ) m on true
where c.id = ANY (@conversation_ids::uuid[])
order by c.id;

-- name: PaginateMessages :many
select *
from "messenger"."message"
where conversation_id = $1
  and (
        sqlc.narg(after_created_at)::timestamp is null
        or created_at < @after_created_at
        or (sqlc.narg(after_id)::uuid is null or (created_at = @after_created_at and id > @after_id))
    )
order by created_at desc, id
limit @first;

-- name: CreateMessage :exec
insert into "messenger"."message" (id, conversation_id, sender_id, text, text_parts, status, type, updated_at,
                                   created_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: FindMessageForUpdate :one
select *
from "messenger"."message"
where id = $1 for update;

-- name: UpdateMessage :exec
update "messenger"."message"
set text       = @text,
    text_parts = @text_parts,
    updated_at = @updated_at
where id = @id;