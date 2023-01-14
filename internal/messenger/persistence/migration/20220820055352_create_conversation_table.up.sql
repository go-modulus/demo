create schema "messenger";

create table "messenger"."conversation"
(
    id          uuid      not null
        constraint conversation_pk primary key,
    sender_id   uuid      not null,
    receiver_id uuid      not null,
    updated_at  timestamp not null,
    created_at  timestamp not null
);

create unique index "conversation_sender_and_receiver_ids" on "messenger"."conversation" (sender_id, receiver_id);

