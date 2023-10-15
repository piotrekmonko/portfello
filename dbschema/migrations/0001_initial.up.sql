create table users
(
    id           varchar(512)            not null
        constraint users_pk
            primary key,
    display_name text                    not null,
    email        text                    not null,
    created_at   timestamp default now() not null
);
