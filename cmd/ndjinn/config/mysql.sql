-- Adminer 4.2.4 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

CREATE DATABASE `ndjinn` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
USE `ndjinn`;

CREATE TABLE `listing` (
  `id` tinyint(1) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `deleted` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `user_id` tinyint(1) unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `NickName` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  `MembershipLevel` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `email` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `password` char(60) COLLATE utf8_unicode_ci NOT NULL,
  `status_id` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `deleted` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `fuserstatus` (`status_id`),
  CONSTRAINT `fuserstatus` FOREIGN KEY (`status_id`) REFERENCES `userstatus` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO `user` (`id`, `NickName`, `MembershipLevel`, `email`, `password`, `status_id`, `created_at`, `updated_at`, `deleted`) VALUES
(1,	'admin',	0,	'test@example.com',	'$2a$10$XOfFtc9HPL3ZEi3r93J/H.RL1QscsEAYkzY3zPmIklJtUYzFQO/8a',	1,	'2016-04-01 12:55:58',	'0000-00-00 00:00:00',	0);

CREATE TABLE `userstatus` (
  `id` tinyint(1) unsigned NOT NULL AUTO_INCREMENT,
  `status` varchar(25) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `deleted` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO `userstatus` (`id`, `status`, `created_at`, `updated_at`, `deleted`) VALUES
(1,	'active',	'2016-04-01 11:08:52',	'2016-04-01 11:08:52',	0),
(2,	'inactive',	'2016-04-01 11:08:52',	'2016-04-01 11:08:52',	0);

-- 2016-04-01 12:58:53
