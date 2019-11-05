use shelob;

CREATE TABLE `unknown_proxy_source` (
  `base64_source` TEXT NOT NULL,
	`ocr_text` VARCHAR(255) DEFAULT ''
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
