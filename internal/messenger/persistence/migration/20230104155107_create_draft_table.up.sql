create table "messenger"."draft"
(
    id              uuid not null
        constraint draft_pk primary key,
    conversation_id uuid not null,
    author_id       uuid not null,
    text            text,
    text_parts      jsonb
);

create unique index "draft_author_and_conversation_ids" on "messenger"."draft" (author_id, conversation_id);
