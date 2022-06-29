set search_path="user";

create table "user"."user"
(
    id            uuid                     not null
        constraint user_pk
            primary key,
    name          varchar(50)              not null,
    email         varchar(127)             not null,
    registered_at timestamp with time zone not null,
    settings      jsonb,
    contacts      text[]
);

