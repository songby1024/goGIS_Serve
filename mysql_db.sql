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

CREATE TABLE messages (
                            `id`          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    -- 主键，自增
                            `created_at`  DATETIME DEFAULT CURRENT_TIMESTAMP,            -- 创建时间
                            `updated_at`  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
                            `deleted_at`  DATETIME,                                      -- 删除时间(用于软删除)
                            `geo_id`      BIGINT NOT NULL,                               -- 围栏标识
                            `client_id`   BIGINT NOT NULL,                               -- 信息接收者标识
                            `address_name` VARCHAR(255),                                 -- 地址名称
                            `alert_time`  VARCHAR(255),                                  -- 警告时间
                            `alert_class` VARCHAR(255),                                  -- 警告类别
                            `alert_dic`   VARCHAR(255),                                  -- 警告字典
                            `min_distance` DOUBLE,                                       -- 最小距离
                            `state`       INT,                                           -- 状态(0为已处理，1为未处理)
                            `point_lat`   DOUBLE,                                        -- 点的纬度
                            `point_lng`   DOUBLE,                                        -- 点的经度
                            INDEX `idx_deleted_at` (`deleted_at`)                        -- 为deleted_at列添加索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;