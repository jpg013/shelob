use shelob;

CREATE TABLE `proxy` (
  `ip_address` VARCHAR(255) NOT NULL,
	`port` MEDIUMINT NOT NULL,
  `protocol` VARCHAR(255),
  `location` VARCHAR(255),
  `created_at` DATETIME, 
  `is_active` BOOLEAN DEFAULT TRUE,
  
  PRIMARY KEY (`ip_address`, `port`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;