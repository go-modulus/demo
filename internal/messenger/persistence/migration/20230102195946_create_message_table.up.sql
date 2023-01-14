create table "messenger"."message"
(
    id              uuid      not null
        constraint message_pk primary key,
    conversation_id uuid      not null,
    sender_id       uuid      not null,
    text            text,
    text_parts      jsonb,
    status          text      not null,
    type            text      not null,
    updated_at      timestamp not null,
    created_at      timestamp not null
);