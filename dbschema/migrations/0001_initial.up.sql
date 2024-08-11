-- Holds user data for auth.LocalProvider. Development use only.
create table local_user
(
    id                  varchar(512)            not null,
    email               varchar(256)            not null
        primary key,
    display_name        text                    not null,
    roles               text                    not null,
    pwdhash             text                    not null,
    created_at          timestamp default CURRENT_TIMESTAMP not null
);

-- Records events.
create table history
(
    id           varchar(512)            not null
        primary key,
    namespace    varchar(64)             not null, /* Holds the resource name related resource, such as table name or auth provider name. */
    reference    varchar(22)             not null, /* Holds the resource ID of the related resource, such as table PK or auth provider id. */
    event        text                    not null, /* Describes the event. */
    email        text                    not null, /* Identifies the user who triggered the event. */
    created_at   timestamp default CURRENT_TIMESTAMP not null
);

-- Holds together a user and different expenses.
create table wallet
(
    id      varchar(22)  not null
        constraint wallet_pk
            primary key, /* A base57-encoded uuid. */
    user_id      varchar(256)            not null, /* User ID reference to auth provider. This is this wallet Owner. */
    balance      double precision        not null,
    currency     varchar(8)              not null,
    created_at   timestamp default CURRENT_TIMESTAMP not null
);

-- Tracks expenses.
create table expense
(
    id         varchar(22)             not null
        constraint expense_pk
            primary key, /* A base57 encoded uuid. */
    wallet_id  varchar(22)             not null
        constraint expense_wallet_id_fk
            references wallet,
    amount     double precision        not null,
    description text,
    created_at timestamp default CURRENT_TIMESTAMP not null
);
