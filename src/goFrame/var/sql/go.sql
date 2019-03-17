/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 50715
Source Host           : 127.0.0.1:3306
Source Database       : go

Target Server Type    : MYSQL
Target Server Version : 50715
File Encoding         : 65001

Date: 2019-03-17 13:06:12
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for go_action
-- ----------------------------
DROP TABLE IF EXISTS `go_action`;
CREATE TABLE `go_action` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `controller` varchar(100) NOT NULL,
  `action` varchar(100) NOT NULL,
  `icon` varchar(32) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `logic_delete` tinyint(1) unsigned DEFAULT '0' COMMENT '逻辑删除',
  `creat_time` int(11) unsigned DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of go_action
-- ----------------------------

-- ----------------------------
-- Table structure for go_config
-- ----------------------------
DROP TABLE IF EXISTS `go_config`;
CREATE TABLE `go_config` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `name` varchar(30) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '配置名称',
  `desc` varchar(50) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '配置说明',
  `settings` text CHARACTER SET utf8 NOT NULL COMMENT '配置值',
  `type` tinyint(4) DEFAULT NULL COMMENT '配置类型',
  `status` tinyint(4) DEFAULT '1',
  `sort` smallint(3) unsigned DEFAULT '0' COMMENT '排序',
  `group_type` tinyint(4) DEFAULT '1' COMMENT '配置分组',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) unsigned DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`) USING BTREE,
  KEY `group` (`group_type`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=79 DEFAULT CHARSET=utf8mb4 COMMENT='网站系统配置表';

-- ----------------------------
-- Records of go_config
-- ----------------------------
INSERT INTO `go_config` VALUES ('78', 'SITE_NAME', '站点名称', '落雨成花', '0', '1', '0', '0', '1545551309', '1545551309');

-- ----------------------------
-- Table structure for go_email
-- ----------------------------
DROP TABLE IF EXISTS `go_email`;
CREATE TABLE `go_email` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of go_email
-- ----------------------------

-- ----------------------------
-- Table structure for go_log
-- ----------------------------
DROP TABLE IF EXISTS `go_log`;
CREATE TABLE `go_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) unsigned NOT NULL,
  `action` varchar(50) DEFAULT NULL COMMENT '方法',
  `ip` varchar(50) DEFAULT NULL,
  `params` text,
  `error` text,
  `create_time` int(11) unsigned DEFAULT NULL,
  `type` int(11) unsigned DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='日志表';

-- ----------------------------
-- Records of go_log
-- ----------------------------

-- ----------------------------
-- Table structure for go_menu
-- ----------------------------
DROP TABLE IF EXISTS `go_menu`;
CREATE TABLE `go_menu` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '上一级菜单id',
  `title_cn` varchar(50) NOT NULL COMMENT '菜单名称',
  `title_en` varchar(50) DEFAULT NULL,
  `icon` varchar(32) NOT NULL DEFAULT '',
  `url_for` varchar(50) DEFAULT NULL COMMENT '菜单url链接',
  `type` int(11) unsigned DEFAULT '1',
  `logic_delete` tinyint(1) unsigned DEFAULT '0' COMMENT '逻辑删除',
  `update_time` int(11) unsigned DEFAULT NULL COMMENT '更新时间',
  `create_time` int(11) unsigned DEFAULT NULL COMMENT '创建时间',
  `sort` int(11) unsigned DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8 COMMENT='后台菜单表';

-- ----------------------------
-- Records of go_menu
-- ----------------------------
INSERT INTO `go_menu` VALUES ('1', '0', '系统管理', null, '', '', '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('2', '1', '菜单管理', null, '', 'MenuController.Index', '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('3', '1', '权限管理', null, '', 'AuthController.Index', '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('4', '1', '角色管理', null, '', 'RoleController.Index', '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('5', '1', '用户管理', null, '', 'UserController.Index', '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('6', '0', '日志管理', null, '', 'LogController.Index', '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('7', '6', '系统日志', null, '', 'LogController.System', '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('8', '6', '操作日志', null, '', 'LogController.Opera', '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('9', '0', '资源管理', null, '', '', '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('10', '9', '文件管理', null, '', null, '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('11', '9', '图片管理', null, '', 'PictureController.Index', '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('12', '9', '视频管理', null, '', null, '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('13', '0', '邮件管理', null, '', null, '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('14', '13', '邮件模板', null, '', null, '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('15', '13', '邮件发送记录', null, '', null, '1', '0', null, null, '0');
INSERT INTO `go_menu` VALUES ('16', '13', '收信箱', null, '', null, '1', '0', null, null, '0');

-- ----------------------------
-- Table structure for go_menu_action_rel
-- ----------------------------
DROP TABLE IF EXISTS `go_menu_action_rel`;
CREATE TABLE `go_menu_action_rel` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `menu_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '方法所属菜单id',
  `action_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '方法id',
  `logic_delete` tinyint(1) unsigned DEFAULT '0' COMMENT '逻辑删除',
  `creat_time` int(11) unsigned DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单和菜单页操作方法的关联表';

-- ----------------------------
-- Records of go_menu_action_rel
-- ----------------------------

-- ----------------------------
-- Table structure for go_picture
-- ----------------------------
DROP TABLE IF EXISTS `go_picture`;
CREATE TABLE `go_picture` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id自增',
  `title` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `path` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '路径',
  `thumb` varchar(255) DEFAULT NULL COMMENT '缩略图',
  `hash` char(64) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '文件md5',
  `size` int(10) unsigned DEFAULT '0',
  `status` int(2) unsigned DEFAULT '0',
  `type` varchar(6) DEFAULT NULL COMMENT '图片类型',
  `create_time` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of go_picture
-- ----------------------------

-- ----------------------------
-- Table structure for go_role
-- ----------------------------
DROP TABLE IF EXISTS `go_role`;
CREATE TABLE `go_role` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name_cn` varchar(50) NOT NULL,
  `name_en` varchar(50) DEFAULT NULL,
  `status` tinyint(4) unsigned DEFAULT '0' COMMENT '状态: 1禁用， 0正常',
  `create_time` int(11) unsigned DEFAULT NULL COMMENT '创建时间',
  `update_time` int(11) unsigned DEFAULT NULL COMMENT '更新时间',
  `logic_delete` tinyint(1) unsigned DEFAULT '0' COMMENT '逻辑删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='用户角色表';

-- ----------------------------
-- Records of go_role
-- ----------------------------
INSERT INTO `go_role` VALUES ('1', '超级管理员', 'superadmin', '0', null, null, '0');
INSERT INTO `go_role` VALUES ('2', '管理员', 'admin', '0', null, null, '0');
INSERT INTO `go_role` VALUES ('3', 'vip会员', 'vip', '0', '1551592291', '0', '1');
INSERT INTO `go_role` VALUES ('4', '普通用户', 'normal user', '0', '1551592684', '0', '0');

-- ----------------------------
-- Table structure for go_role_actione_rel
-- ----------------------------
DROP TABLE IF EXISTS `go_role_actione_rel`;
CREATE TABLE `go_role_actione_rel` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` int(11) unsigned NOT NULL,
  `action_id` varchar(100) NOT NULL,
  `logic_delete` tinyint(1) unsigned DEFAULT '0' COMMENT '逻辑删除',
  `creat_time` int(11) unsigned DEFAULT NULL COMMENT '控制器-方法名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色和角色可执行的方法的关联表';

-- ----------------------------
-- Records of go_role_actione_rel
-- ----------------------------

-- ----------------------------
-- Table structure for go_role_menu_rel
-- ----------------------------
DROP TABLE IF EXISTS `go_role_menu_rel`;
CREATE TABLE `go_role_menu_rel` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` int(11) unsigned NOT NULL,
  `menu_id` int(11) unsigned NOT NULL,
  `logic_delete` tinyint(1) unsigned DEFAULT '0' COMMENT '逻辑删除',
  `create_time` int(11) unsigned DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of go_role_menu_rel
-- ----------------------------
INSERT INTO `go_role_menu_rel` VALUES ('1', '1', '0', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('2', '1', '1', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('3', '1', '2', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('4', '1', '3', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('5', '1', '4', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('6', '1', '5', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('7', '1', '6', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('8', '1', '7', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('9', '1', '8', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('10', '1', '9', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('11', '1', '10', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('12', '1', '11', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('13', '1', '12', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('14', '1', '13', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('15', '1', '14', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('16', '1', '15', '0', null);
INSERT INTO `go_role_menu_rel` VALUES ('17', '1', '16', '0', null);

-- ----------------------------
-- Table structure for go_role_user_rel
-- ----------------------------
DROP TABLE IF EXISTS `go_role_user_rel`;
CREATE TABLE `go_role_user_rel` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` int(11) unsigned NOT NULL COMMENT '角色id',
  `user_id` int(11) unsigned NOT NULL COMMENT '用户id',
  `logic_delete` tinyint(1) unsigned DEFAULT '0' COMMENT '逻辑删除',
  `create_time` int(11) unsigned DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='角色和角色可访问菜单的关联表';

-- ----------------------------
-- Records of go_role_user_rel
-- ----------------------------
INSERT INTO `go_role_user_rel` VALUES ('1', '1', '1', '0', null);
INSERT INTO `go_role_user_rel` VALUES ('2', '1', '13', '0', null);
INSERT INTO `go_role_user_rel` VALUES ('3', '3', '13', '0', null);
INSERT INTO `go_role_user_rel` VALUES ('4', '2', '1', '0', null);

-- ----------------------------
-- Table structure for go_task
-- ----------------------------
DROP TABLE IF EXISTS `go_task`;
CREATE TABLE `go_task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `group_id` int(11) NOT NULL DEFAULT '0' COMMENT '分组ID',
  `task_name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名称',
  `task_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '任务类型',
  `description` varchar(200) NOT NULL DEFAULT '' COMMENT '任务描述',
  `cron_spec` varchar(100) NOT NULL DEFAULT '' COMMENT '时间表达式',
  `concurrent` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否只允许一个实例',
  `command` text NOT NULL COMMENT '命令详情',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0停用 1启用',
  `notify` tinyint(4) NOT NULL DEFAULT '0' COMMENT '通知设置',
  `notify_email` text NOT NULL COMMENT '通知人列表',
  `timeout` smallint(6) NOT NULL DEFAULT '0' COMMENT '超时设置',
  `execute_times` int(11) NOT NULL DEFAULT '0' COMMENT '累计执行次数',
  `prev_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '上次执行时间',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of go_task
-- ----------------------------

-- ----------------------------
-- Table structure for go_task_group
-- ----------------------------
DROP TABLE IF EXISTS `go_task_group`;
CREATE TABLE `go_task_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `group_name` varchar(50) NOT NULL DEFAULT '' COMMENT '组名',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of go_task_group
-- ----------------------------

-- ----------------------------
-- Table structure for go_task_log
-- ----------------------------
DROP TABLE IF EXISTS `go_task_log`;
CREATE TABLE `go_task_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `output` mediumtext NOT NULL COMMENT '任务输出',
  `error` text NOT NULL COMMENT '错误信息',
  `status` tinyint(4) NOT NULL COMMENT '状态',
  `process_time` int(11) NOT NULL DEFAULT '0' COMMENT '消耗时间/毫秒',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of go_task_log
-- ----------------------------

-- ----------------------------
-- Table structure for go_user
-- ----------------------------
DROP TABLE IF EXISTS `go_user`;
CREATE TABLE `go_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(20) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户名',
  `password` char(32) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '密码',
  `salt` char(10) CHARACTER SET utf8 DEFAULT '' COMMENT '密码盐',
  `email` varchar(50) CHARACTER SET utf8 DEFAULT '' COMMENT '邮箱',
  `mobile` varchar(16) DEFAULT NULL COMMENT '联系电话',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像图片路径',
  `last_login` int(11) DEFAULT '0' COMMENT '最后登录时间',
  `last_ip` varchar(15) CHARACTER SET utf8 DEFAULT '' COMMENT '最后登录IP',
  `role_id` int(11) unsigned DEFAULT '1' COMMENT '关联用户等级id',
  `logic_delete` tinyint(1) unsigned DEFAULT '0' COMMENT '逻辑删除',
  `status` tinyint(4) unsigned DEFAULT '0' COMMENT '状态: 1禁用， 0正常',
  `is_super` tinyint(1) unsigned DEFAULT '0',
  `create_time` int(11) unsigned DEFAULT NULL,
  `update_time` int(10) unsigned DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_name` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of go_user
-- ----------------------------
INSERT INTO `go_user` VALUES ('1', 'admin', '7fef6171469e80d32c0559f88b377245', '', 'admin@example.com', '123456', '', '1552744805', '[', '0', '0', '1', '0', '0', '0');
INSERT INTO `go_user` VALUES ('2', '张三', '7fef6171469e80d32c0559f88b377245', '', 'admin@example.com', '123456', null, '1547963364', '', '1', '0', '1', null, null, '0');
INSERT INTO `go_user` VALUES ('4', '李四', 'd41d8cd98f00b204e9800998ecf8427e', '', 'admin@example.com', '123456', null, '1547963364', '', '1', '0', '0', null, null, '0');
INSERT INTO `go_user` VALUES ('5', '王五', 'd41d8cd98f00b204e9800998ecf8427e', '', 'admin@example.com', '123456', null, '1547963364', '', '1', '0', '0', null, null, '0');
INSERT INTO `go_user` VALUES ('6', '赵六', 'd41d8cd98f00b204e9800998ecf8427e', '', 'admin@example.com', '123456', null, '1547963364', '', '1', '0', '0', null, null, '0');
INSERT INTO `go_user` VALUES ('7', 'first', 'd41d8cd98f00b204e9800998ecf8427e', '', 'admin@example.com', '123456', null, '1547963364', '', '1', '0', '1', null, null, '0');
INSERT INTO `go_user` VALUES ('8', 'two', 'd41d8cd98f00b204e9800998ecf8427e', '', 'admin@example.com', '123456', null, '1547963364', '', '1', '0', '0', null, null, '0');
INSERT INTO `go_user` VALUES ('9', 'third', 'd41d8cd98f00b204e9800998ecf8427e', '', 'admin@example.com', '123456', null, '1547963364', '', '1', '0', '0', null, null, '0');
INSERT INTO `go_user` VALUES ('10', 'four', 'd41d8cd98f00b204e9800998ecf8427e', '', 'admin@example.com', '13696888888', '', '0', '', '1', '0', '1', '1', '0', '0');
INSERT INTO `go_user` VALUES ('11', 'five', 'd41d8cd98f00b204e9800998ecf8427e', '', 'admin@example.com', '123456', null, '1547963364', '', '1', '0', '0', null, null, '0');
INSERT INTO `go_user` VALUES ('12', 'six', 'd41d8cd98f00b204e9800998ecf8427e', '', 'admin@example.com', '123456', null, '1547963364', '', '1', '0', '0', null, null, '0');
INSERT INTO `go_user` VALUES ('13', 'seven', 'd41d8cd98f00b204e9800998ecf8427e', '', 'admin@example.com', '123456', '', '1547963364', '127.0.0.1', '1', '0', '0', '1', '1547963364', '1547963364');
INSERT INTO `go_user` VALUES ('14', 'sss', 'd41d8cd98f00b204e9800998ecf8427e', '', '16523@qq.com', '13222245444', '', '0', '', '1', '0', '0', '1', null, '0');
INSERT INTO `go_user` VALUES ('15', 'www', 'd41d8cd98f00b204e9800998ecf8427e', '', '222222@qq.com', '13654121222', '', '0', '', '1', '1', '0', '0', null, '0');
INSERT INTO `go_user` VALUES ('16', 'qqq', 'd41d8cd98f00b204e9800998ecf8427e', '', '1@qq.com', '13659999999', '', '0', '', '1', '0', '0', '1', null, '0');
INSERT INTO `go_user` VALUES ('17', '1111', 'd41d8cd98f00b204e9800998ecf8427e', '', '123@qq.com', '13222245444', '', '0', '', '1', '1', '0', '1', null, '0');
INSERT INTO `go_user` VALUES ('18', '11', 'd41d8cd98f00b204e9800998ecf8427e', '', '2@qq.com', '13623232333', '', '0', '', '1', '1', '1', '0', '1550149216', '0');
INSERT INTO `go_user` VALUES ('20', 'ddddd', 'd41d8cd98f00b204e9800998ecf8427e', '', 'a@qq.com', '13655555555', '', '0', '', '1', '0', '0', '0', '0', '1551593942');
