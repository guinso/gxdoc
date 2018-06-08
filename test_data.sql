-- Adminer 4.3.1 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `doc_schema`;
CREATE TABLE `doc_schema` (
  `id` int(11) NOT NULL,
  `name` char(100) NOT NULL,
  `description` text NOT NULL,
  `is_active` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `doc_schema` (`id`, `name`, `description`, `is_active`) VALUES
(1,	'invoice',	'invoice....',	1),
(2,	'pr',	'purchase requisite',	1);

DROP TABLE IF EXISTS `doc_schema_revision`;
CREATE TABLE `doc_schema_revision` (
  `schema_id` int(11) NOT NULL,
  `revision` int(11) NOT NULL,
  `xml_definition` text NOT NULL,
  `remark` text NOT NULL,
  PRIMARY KEY (`schema_id`,`revision`),
  CONSTRAINT `doc_schema_revision_ibfk_1` FOREIGN KEY (`schema_id`) REFERENCES `doc_schema` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `doc_schema_revision` (`schema_id`, `revision`, `xml_definition`, `remark`) VALUES
(1,	1,	'<dxdoc name=\"invoice\" revision=\"1\" id=\"1\"><dxstr name=\"invNo\"></dxstr><dxint name=\"totalQty\" isOptional=\"true\"></dxint><dxdecimal name=\"price\" precision=\"2\"></dxdecimal></dxdoc>',	''),
(1,	2,	'<dxdoc name=\"invoice\" revision=\"2\" id=\"1\"><dxstr name=\"invNo\"></dxstr><dxint name=\"totalQty\" isOptional=\"true\"></dxint><dxdecimal name=\"price\" precision=\"2\"></dxdecimal></dxdoc>',	''),
(2,	-1,	'<dxdoc name=\"pr\" revision=\"1\" id=\"2\">\r\n<dxint name=\"qty\"></dxint>\r\n<dxstr name=\"pr number\" lenLimit=\"6\"></dxstr>\r\n</dxdoc>',	''),
(2,	1,	'<dxdoc name=\"pr\" revision=\"1\" id=\"2\">\r\n<dxint name=\"qty\"></dxint>\r\n<dxstr name=\"pr number\" lenLimit=\"6\"></dxstr>\r\n</dxdoc>',	'');

-- 2018-06-08 10:17:35