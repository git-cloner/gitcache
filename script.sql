CREATE TABLE `gitcache_repos` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `path` varchar(255) NOT NULL DEFAULT '',
  `ctime` timestamp NULL DEFAULT NULL,
  `utime` timestamp NULL DEFAULT NULL,
  `hitcount` int(11) DEFAULT '0',
  `tcount` int(11) NOT NULL DEFAULT '10',
  `starcount` int(11) DEFAULT '0',
  `language` varchar(20) DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `owner` bigint(20) DEFAULT NULL,
  `description` varchar(1000) DEFAULT NULL,
  `last_recommendtime` bigint(20) DEFAULT '0',
  `size` int(11) DEFAULT '0',
  PRIMARY KEY (`Id`),
  UNIQUE KEY `idx_gitcache_repos` (`path`),
  KEY `idx_gitcache_repos1` (`name`),
  KEY `idx_gitcache_repos_l` (`last_recommendtime`)
) ENGINE=InnoDB AUTO_INCREMENT=47295 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `gitcache_stats` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `stime` timestamp NULL DEFAULT NULL,
  `cachehit` decimal(10,0) DEFAULT '0',
  `redirect` decimal(10,0) DEFAULT '0',
  `visit` decimal(10,0) DEFAULT '0',
  `vipvisit` decimal(10,0) DEFAULT '0',
  `search` decimal(10,0) DEFAULT '0',
  `imagetest` decimal(10,0) DEFAULT '0',
  `githubapp` decimal(10,0) DEFAULT '0',
  `githubdesktop` decimal(10,0) DEFAULT '0',
  `githubcli` decimal(10,0) DEFAULT '0',
  `gitexe` decimal(10,0) DEFAULT '0',
  PRIMARY KEY (`Id`),
  UNIQUE KEY `inx_stats_time` (`stime`)
) ENGINE=InnoDB AUTO_INCREMENT=276 DEFAULT CHARSET=utf8mb4;

