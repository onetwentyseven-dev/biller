/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE TABLE IF NOT EXISTS `providers` (
  `id` varchar(255) NOT NULL,
  `user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `name` varchar(255) NOT NULL,
  `web_address` varchar(255) DEFAULT NULL,
  `ts_created` datetime(2) NOT NULL,
  `ts_updated` datetime(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_user_id_unique` (`name`,`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `receipts` (
  `id` varchar(255) NOT NULL,
  `user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `provider_id` varchar(255) DEFAULT NULL,
  `label` varchar(255) NOT NULL,
  `date_paid` datetime NOT NULL,
  `amount_paid` double(10,2) NOT NULL,
  `ts_created` datetime NOT NULL,
  `ts_updated` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `provider_id_foreign_key` (`provider_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `provider_id_foreign_key` FOREIGN KEY (`provider_id`) REFERENCES `providers` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `bills` (
  `id` varchar(255) NOT NULL,
  `provider_id` varchar(255) NOT NULL,
  `user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `name` varchar(255) NOT NULL,
  `ts_created` datetime(2) NOT NULL,
  `ts_updated` datetime(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_user_id_provider_id` (`name`,`provider_id`,`user_id`) USING BTREE,
  KEY `FK_bills_providers` (`provider_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `FK_bills_providers` FOREIGN KEY (`provider_id`) REFERENCES `providers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `bill_sheets` (
  `id` varchar(255) NOT NULL,
  `user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `name` varchar(255) NOT NULL,
  `ts_created` datetime(2) NOT NULL,
  `ts_updated` datetime(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_name_unique` (`user_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `bill_sheet_entries` (
  `entry_id` varchar(255) NOT NULL,
  `sheet_id` varchar(255) NOT NULL,
  `bill_id` varchar(255) NOT NULL,
  `date_due` datetime NOT NULL,
  `amount_due` double(10,2) NOT NULL,
  `receipt_id` varchar(255) DEFAULT NULL,
  `date_paid` datetime DEFAULT NULL,
  `amount_paid` double(10,2) DEFAULT NULL,
  `ts_created` datetime NOT NULL,
  `ts_updated` datetime NOT NULL,
  PRIMARY KEY (`entry_id`),
  KEY `bill_id_foreign_key` (`bill_id`),
  KEY `sheet_id_foreign_key` (`sheet_id`),
  KEY `receipt_id_foreign_key` (`receipt_id`),
  CONSTRAINT `bill_id_foreign_key` FOREIGN KEY (`bill_id`) REFERENCES `bills` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `receipt_id_foreign_key` FOREIGN KEY (`receipt_id`) REFERENCES `receipts` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `sheet_id_foreign_key` FOREIGN KEY (`sheet_id`) REFERENCES `bill_sheets` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
