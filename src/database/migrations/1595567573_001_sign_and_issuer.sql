ALTER TABLE `updates` ADD `issuer` INT(8) UNSIGNED NOT NULL,
	ADD `url_hash` VARCHAR(64) NOT NULL,
	ADD `sig_url` VARCHAR(512) NOT NULL,
	ADD `sig_url_hash` VARCHAR(64) NOT NULL;
