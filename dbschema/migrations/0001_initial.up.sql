create table local_user
(
    email               varchar(256)            not null
        primary key,
    display_name        text                    not null,
    roles               text                    not null,
    pwdhash             text                    not null,
    created_at          timestamp default now() not null
);

comment on table local_user is 'Holds user data for auth.LocalProvider. Development use only.';


create table history
(
    id           varchar(512)            not null
            primary key,
    namespace    varchar(64)             not null,
    reference    varchar(22)             not null,
    event        text                    not null,
    email        text                    not null,
    created_at   timestamp default now() not null
);

comment on column history.namespace is 'Holds the resource name related resource, such as table name or auth provider name.';

comment on column history.reference is 'Holds the resource ID of the related resource, such as table PK or auth provider id.';

comment on column history.event is 'Describes the event.';

comment on column history.email is 'Identifies the user who triggered the event.';


create table wallet
(
    id      varchar(22)  not null
        constraint wallet_pk
            primary key,
    user_id      varchar(256)            not null,
    balance      double precision        not null,
    currency     varchar(8)              not null,
    created_at   timestamp default now() not null
);

comment on table wallet is 'Holds together a user and different expenses.';

comment on column wallet.id is 'A base57-encoded uuid.';

comment on column wallet.user_id is 'User ID reference to auth provider. This is this wallet Owner.';


create table expense
(
    id         varchar(22)             not null
        constraint expense_pk
            primary key,
    wallet_id  varchar(22)             not null
        constraint expense_wallet_id_fk
            references wallet,
    amount     double precision        not null,
    description text,
    created_at timestamp default now() not null
);

comment on table expense is 'Tracks expenses.';

comment on column expense.id is 'A base57 encoded uuid.';

comment on column expense.wallet_id is 'Reference to the wallet.';