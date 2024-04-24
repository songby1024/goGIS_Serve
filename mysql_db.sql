-- auto-generated definition
create table user
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    deleted_at datetime(3) null,
    username   longtext    null,
    password   longtext    null,
    email      longtext    null,
    ruler      int         null
)
    charset = utf8mb3
    row_format = COMPACT;

create index idx_user_deleted_at
    on user (deleted_at);