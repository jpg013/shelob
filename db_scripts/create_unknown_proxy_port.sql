use shelob;

CREATE TABLE `unknown_proxy_port` (
  `file_path` VARCHAR(255) NOT NULL,
	`ocr_text` VARCHAR(255) DEFAULT '',
  `port` MEDIUMINT,
	PRIMARY KEY (`file_path`),
	INDEX `file_path_idx` (`file_path`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
