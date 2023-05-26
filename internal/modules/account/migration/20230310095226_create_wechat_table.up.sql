CREATE TABLE wechat (
    id BIGINT(20) UNSIGNED PRIMARY KEY,
    account_id BIGINT(20) UNSIGNED NOT NULL COMMENT '账户id',
    open_id VARCHAR(255) NOT NULL COMMENT 'openid',
    type VARCHAR(255) NOT NULL COMMENT '微信类型',
    union_id VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'unionid 备用',
    nickname VARCHAR(255) NOT NULL DEFAULT '' COMMENT '昵称',
    avatar_url VARCHAR(255) DEFAULT '' COMMENT '头像地址',
    mobile VARCHAR(20) DEFAULT '' COMMENT '手机号',
    created_at BIGINT(20) NOT NULL DEFAULT 0 COMMENT '创建时间',
    updated_at BIGINT(20) NOT NULL DEFAULT 0 COMMENT '更新时间'
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT = '微信登录信息表' COLLATE = utf8mb4_unicode_ci;