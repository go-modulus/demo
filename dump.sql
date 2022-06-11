create schema "user";

create table "user"."user"
(
    id    uuid
        constraint user_pk
            primary key,
    name  varchar(50)  not null,
    email varchar(127) not null,
    registered_at timestamp with time zone not null
);

create unique index user_email_uindex
    on "user"."user" (email);

INSERT INTO "user"."user" (id, name, email, registered_at) VALUES ('00000000-0000-0000-0000-000000000001', 'Test1', 'test1@test.com', '2022-01-01 10:00:00');
INSERT INTO "user"."user" (id, name, email, registered_at) VALUES ('00000000-0000-0000-0000-000000000002', 'Test2', 'test2@test.com', '2022-01-01 11:00:00');
