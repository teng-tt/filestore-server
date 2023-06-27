create database fileserver default character set utf8;

create table `tbl_file` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT "文件hash",
    `file_name` char(256) NOT NULL DEFAULT '' COMMENT "文件名",
    `file_size` bigint(20) DEFAULT "0" COMMENT "文件大小",
    `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT "文件存储位置",
    `create_at` datetime DEFAULT NOW() COMMENT "创建日期",
    `update_at` datetime DEFAULT  NOW() on update current_timestamo() COMMENT "更新日期",
    `status` int(11) NOT NULL DEFAULT '0' COMMENT "状态（可用/禁用/已删除）",
    `ext1` int(11) DEFAULT '0' COMMENT "备用字段1",
    `ext2` text COMMENT "备用字段2",
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_file_hash` (`file_sha1`),
        KEY `idx_status` (`status`)
)   ENGINE = InnoDB DEFAULT CHARSET=utf8;