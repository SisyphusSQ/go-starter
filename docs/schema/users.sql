CREATE TABLE `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名称',
    `email` varchar(128) NOT NULL DEFAULT '' COMMENT '用户邮箱',
    `password` varchar(128) NOT NULL DEFAULT '' COMMENT '密码(bcrypt哈希)',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unq_email` (`email`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户信息表';