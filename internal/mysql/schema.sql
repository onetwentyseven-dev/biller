CREATE TABLE `providers` (
	`id` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_general_ci',
	`name` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_general_ci',
	`ts_created` DATETIME(2) NOT NULL,
	`ts_updated` DATETIME(2) NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `name` (`name`) USING BTREE
)
COMMENT='Providers is any entity that provides a service for which a bill is recieved'
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `bills` (
	`id` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_general_ci',
	`name` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_general_ci',
	`provider_id` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_general_ci',
	`ts_created` DATETIME(2) NOT NULL,
	`ts_updated` DATETIME(2) NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `name_provider_id` (`name`, `provider_id`) USING BTREE,
	INDEX `FK_bills_providers` (`provider_id`) USING BTREE,
	CONSTRAINT `FK_bills_providers` FOREIGN KEY (`provider_id`) REFERENCES `providers` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
