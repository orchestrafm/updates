ALTER TABLE `updates` ADD `issuer` INT(8) UNSIGNED NOT NULL,
`url_hash` VARCHAR(64) NOT NULL,
`sig_url` VARCHAR(512) NOT NULL,
`sig_url_hash` VARCHAR(64) NOT NULL;
