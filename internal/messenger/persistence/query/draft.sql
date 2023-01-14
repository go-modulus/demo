-- name: FindOrCreateDraft :one
with new_draft as (
    insert into "messenger"."draft" (id, conversation_id, author_id)
        select @id, @conversation_id, @author_id
        where not exists(
                select 1
                from "messenger"."draft"
                where author_id = @author_id
                  and conversation_id = @conversation_id
            )
        on conflict do nothing
        returning *) (select *
                      from "messenger"."draft"
                      where author_id = @author_id
                        and conversation_id = @conversation_id

                      union all

                      select *
                      from new_draft) limit 1;

-- name: FindDrafts :many
select *
from "messenger"."draft"
where author_id = @author_id
  and conversation_id = ANY (@conversation_ids::uuid[]);

-- name: FindDraftForUpdate :one
select *
from "messenger"."draft"
where id = $1 for update;

-- name: UpdateDraft :exec
update "messenger"."draft"
set text       = @text,
    text_parts = @text_parts
where id = @id;

-- name: RemoveDraft :exec
delete
from "messenger"."draft"
where author_id = $1
  and conversation_id = $2;
