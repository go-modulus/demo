-- name: GetConversation :one
select *
from "messenger"."conversation"
where id = $1
limit 1;

-- name: CreateOrGetConversation :one
with new_conversation as (
    insert into "messenger"."conversation" (id, sender_id, receiver_id, updated_at, created_at)
        select @id, @sender_id, @receiver_id, @updated_at, @created_at
        where not exists(
                select 1
                from "messenger"."conversation"
                where (sender_id = @sender_id and receiver_id = @receiver_id)
                   or (receiver_id = @sender_id and sender_id = @receiver_id)
            )
        on conflict do nothing
        returning *) (select *
                      from "messenger"."conversation"
                      where (sender_id = @sender_id and receiver_id = @receiver_id)
                         or (receiver_id = @sender_id and sender_id = @receiver_id)

                      union all

                      select *
                      from new_conversation) limit 1;

-- name: PaginateMyConversations :many
select *
from "messenger"."conversation"
where (sender_id = @viewer_id or receiver_id = @viewer_id)
  and (
        sqlc.narg(after_updated_at)::timestamp is null
        or updated_at < @after_updated_at
        or (sqlc.narg(after_id)::uuid is null or (updated_at = @after_updated_at and id > @after_id))
    )
order by updated_at desc, id
limit @first;