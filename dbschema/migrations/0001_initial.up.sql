create table history
(
    id           varchar(512)            not null
            primary key,
    event        text                    not null,
    email        text                    not null,
    created_at   timestamp default now() not null
);
