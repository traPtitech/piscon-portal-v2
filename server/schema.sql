CREATE TABLE IF NOT EXISTS `teams` (
    `id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `team_id` VARCHAR(36),
    `is_admin` TINYINT(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`),
    FOREIGN KEY (`team_id`) REFERENCES `teams` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `sessions` (
    `id` VARCHAR(50) NOT NULL,
    `user_id` VARCHAR(36) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `expired_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `instances` (
    `id` VARCHAR(36) NOT NULL,
    `team_id` VARCHAR(36) NOT NULL,
    `instance_number` INT NOT NULL,
    `status` ENUM('running', 'building', 'starting', 'stopping', 'stopped', 'deleting', 'deleted') NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `public_ip` VARCHAR(15),
    `private_ip` VARCHAR(15),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `benchmarks` (
    `id` VARCHAR(36) NOT NULL,
    `instance_id` VARCHAR(36) NOT NULL,
    `team_id` VARCHAR (36) NOT NULL,
    `user_id` VARCHAR(36) NOT NULL,
    `status` ENUM('waiting', 'running', 'finished') NOT NULL DEFAULT 'waiting',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `started_at` TIMESTAMP,
    `finished_at` TIMESTAMP,
    `score` BIGINT NOT NULL DEFAULT 0,
    `result` ENUM('passed', 'failed', 'error'),
    `error_message` TEXT,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`instance_id`) REFERENCES `instances` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `benchmark_logs` (
    `benchmark_id` VARCHAR(36) NOT NULL,
    `user_log` TEXT NOT NULL,
    `admin_log` TEXT NOT NULL,
    PRIMARY KEY (`benchmark_id`),
    FOREIGN KEY (`benchmark_id`) REFERENCES `benchmarks` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `documents` (
    `id` VARCHAR(36) NOT NULL,
    `body` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;