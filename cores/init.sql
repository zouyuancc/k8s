-- 创建Kubernetes数据库
alter user 'root'@'localhost' identified with mysql_native_password by '123456';
CREATE DATABASE Kubernetes;
USE Kubernetes;

-- 创建task表
create table `task` (
    `task_id` char(36) NOT NULL,
    `kind` varchar(256) NOT NULL,
    `user` varchar(256) NOT NULL,
    `submit_time` datetime NOT NULL, 
    PRIMARY KEY (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 创建deployment表
create table `deployment` (
    `task_id` char(36) NOT NULL,
    `name` varchar(256) NOT NULL,
    `metadata` json,
    `replicas` int NOT NULL,
    `selector` json,
    `template_metadata` json,
    FOREIGN KEY (`task_id`) REFERENCES task(`task_id`),
    PRIMARY KEY (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 创建container表
create table `container` (
    `container_id` char(36) NOT NULL,
    `task_id` char(36) NOT NULL,
    `image` varchar(256) NOT NULL,
    `name` varchar(256) NOT NULL,
    `command` json,
    `args` json,
    `resources` json,
    FOREIGN KEY (`task_id`) REFERENCES task(`task_id`),
    PRIMARY KEY (`container_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 创建service表
create table `service` (
    `task_id` char(36) NOT NULL,
    `name` varchar(256) NOT NULL,
    -- `apiVersion` varchar(16),
    `metadata` json,
    `type` varchar(256) NOT NULL,
    `selector` json,
    FOREIGN KEY (`task_id`) REFERENCES task(`task_id`),
    PRIMARY KEY(`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 创建service_port表
create table `service_port` (
    `p_id` char(36) NOT NULL,
    `task_id` char(36) NOT NULL,
    `port` int,
    `nodePort` int,
    `targetPort` int,
    FOREIGN KEY (`task_id`) REFERENCES task(`task_id`),
    PRIMARY KEY(`p_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;