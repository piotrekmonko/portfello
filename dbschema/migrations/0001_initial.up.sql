create table users
(
    id           varchar(512)            not null
        constraint users_pk
            primary key,
    display_name text                    not null,
    is_active    bool      default false not null,
    created_at   timestamp default now() not null
);

