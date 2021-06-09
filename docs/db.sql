-- ----------------------------
-- Table structure for yapi_test_api
-- ----------------------------
DROP TABLE IF EXISTS `yapi_test_api`;
CREATE TABLE `yapi_test_api`
(
    `id`      int(10) unsigned NOT NULL AUTO_INCREMENT,
    `project` varchar(40)  NOT NULL,
    `sid`     int(11) NOT NULL,
    `yid`     int(11) NOT NULL,
    `method`  varchar(40)  NOT NULL,
    `title`   varchar(255) NOT NULL,
    `path`    varchar(255) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=491 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for yapi_report_api
-- ----------------------------
DROP TABLE IF EXISTS `yapi_report_api`;
CREATE TABLE `yapi_report_api`
(
    `id`          int(11) unsigned NOT NULL AUTO_INCREMENT,
    `project`     varchar(40)  NOT NULL,
    `task_id`     varchar(128),
    `sid`         int(11) NOT NULL,
    `yid`         int(11) NOT NULL,
    `method`      varchar(40)  NOT NULL,
    `title`       varchar(128) NOT NULL,
    `path`        varchar(255) NOT NULL,
    `full_path`   varchar(255) NOT NULL,
    `is_coverage` TINYINT(1) DEFAULT 0,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for yapi_test_member
-- ----------------------------
DROP TABLE IF EXISTS `yapi_test_member`;
CREATE TABLE `yapi_test_member`
(
    `id`     int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`   varchar(50) NOT NULL,
    `alias`  varchar(50),
    `group`  varchar(20) NOT NULL,
    `status` TINYINT(1) DEFAULT 0,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;