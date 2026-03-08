CREATE TABLE IF NOT EXISTS `email_send_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '模板 ID',
    `trigger_code` VARCHAR(100) NOT NULL COMMENT '触发点代码',
    `language` VARCHAR(10) NOT NULL COMMENT '语言',
    `recipient` VARCHAR(500) NOT NULL COMMENT '收件人',
    `subject` VARCHAR(500) NOT NULL COMMENT '邮件主题（渲染后）',
    `params` JSON DEFAULT NULL COMMENT '发送参数',
    `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/sent/failed',
    `error_message` TEXT COMMENT '错误信息',
    `sent_at` DATETIME DEFAULT NULL COMMENT '发送时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_trigger` (`trigger_code`),
    KEY `idx_status` (`status`),
    KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='邮件发送日志';
