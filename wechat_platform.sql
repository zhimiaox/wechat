/*
 Navicat Premium Data Transfer

 Source Server         : mysql.zhimiao.org
 Source Server Type    : MySQL
 Source Server Version : 80013
 Source Host           : mysql.zhimiao.org:3306
 Source Schema         : wechat_platform

 Target Server Type    : MySQL
 Target Server Version : 80013
 File Encoding         : 65001

 Date: 09/05/2020 14:57:57
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for miniprogram
-- ----------------------------
DROP TABLE IF EXISTS `miniprogram`;
CREATE TABLE `miniprogram`  (
  `app_id` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '小程序appid',
  `platform_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '开放平台ID',
  `mch_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '支付商户号id',
  `original_id` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '原始ID',
  `refresh_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '接口调用凭据刷新令牌',
  `secret` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '小程序secret',
  `ext_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '小程序扩展配置',
  `state` tinyint(3) NOT NULL DEFAULT 1 COMMENT '-1-授权失效 1授权成功，2审核中，3审核通过，4审核失败，5已发布 6已撤审',
  `version` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '当前版本',
  `now_template_id` int(10) UNSIGNED NULL DEFAULT NULL COMMENT '当前模板ID',
  `template_listen` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '模板监听开发小程序(appid)',
  `audit_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '审核编号',
  `auto_audit` tinyint(2) NOT NULL DEFAULT -1 COMMENT '自动提审(升级) -1否 1是',
  `auto_release` tinyint(2) NOT NULL DEFAULT -1 COMMENT '自动发布-1否 1是',
  `create_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0),
  PRIMARY KEY (`app_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '微信小程序授权表列表' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for miniprogram_audit
-- ----------------------------
DROP TABLE IF EXISTS `miniprogram_audit`;
CREATE TABLE `miniprogram_audit`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `app_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '小程序appId',
  `original_id` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '小程序原始id',
  `audit_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '审核编号',
  `state` tinyint(3) NOT NULL DEFAULT 1 COMMENT '审核状态，-1-撤销审核 1为审核中，2为审核成功，3为审核失败',
  `reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '当status=1，审核被拒绝时，返回的拒绝原因',
  `screen_shot` varchar(3000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '附件素材',
  `template_id` int(11) NOT NULL COMMENT '最新提交审核或者发布的模板id',
  `template_app_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '模板开发小程序ID',
  `template_app_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '开发小程序名',
  `template_app_developer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '开发者',
  `template_desc` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '模板描述',
  `template_version` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '模板版本号',
  `create_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0),
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '微信小程序提交审核的小程序' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for pay
-- ----------------------------
DROP TABLE IF EXISTS `pay`;
CREATE TABLE `pay`  (
  `mch_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '支付商户号id',
  `token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '支付密钥',
  `cert` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '支付证书',
  `pay_notify_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '支付回调',
  `pay_refund_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '退款回调',
  `create_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0),
  PRIMARY KEY (`mch_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '微信支付' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for pay_order
-- ----------------------------
DROP TABLE IF EXISTS `pay_order`;
CREATE TABLE `pay_order`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '订单号，主键',
  `uid` int(11) UNSIGNED NULL DEFAULT NULL COMMENT '购买的用户',
  `pay_type` tinyint(3) UNSIGNED NULL DEFAULT NULL COMMENT '支付方式 2-支付宝 3-微信 4-现金收银',
  `buy_type` tinyint(3) UNSIGNED NULL DEFAULT 1 COMMENT '购买商品的类型：1-购买商品 2-充值 3-发票 4-会员 5-分销升级',
  `status` tinyint(3) NULL DEFAULT 0 COMMENT '状态 0-待支付 1-成功',
  `amount` int(11) UNSIGNED NULL DEFAULT NULL COMMENT '商品金额，单位分',
  `buy_goods_key` int(11) UNSIGNED NULL DEFAULT 0 COMMENT '子订单号',
  `extra` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '附加字段，备用',
  `transaction_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '支付平台交易号',
  `pay_app_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '商户账号',
  `create_time` int(11) NULL DEFAULT 0 COMMENT '创建订单的时间',
  `pay_succ_time` int(11) NULL DEFAULT 0 COMMENT '支付成功的时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `buy_goods_key`(`buy_goods_key`, `pay_succ_time`, `buy_type`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '支付订单' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for pay_refund
-- ----------------------------
DROP TABLE IF EXISTS `pay_refund`;
CREATE TABLE `pay_refund`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `uid` int(11) NULL DEFAULT 0 COMMENT '用户id',
  `refund_status` tinyint(1) NULL DEFAULT 0 COMMENT '退款状态 -1-失败 0-退款中 1-成功',
  `refund_money` int(11) NULL DEFAULT 0 COMMENT '当前退款单退款金额',
  `refund_msg` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '退款备注',
  `refund_total` int(11) NULL DEFAULT 0 COMMENT '已退款金额(不包含此笔退款)',
  `pay_money` int(11) NULL DEFAULT 0 COMMENT '此笔交易金额',
  `pay_type` tinyint(1) NULL DEFAULT 0 COMMENT '支付类型 2-支付宝  3-微信',
  `pay_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '支付单号',
  `pay_app_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '应用id',
  `transaction_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '第三方流水号',
  `refund_from` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '退款来源',
  `result_log` varchar(3000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '返回结果日志',
  `create_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0),
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '第三方退款' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for platform
-- ----------------------------
DROP TABLE IF EXISTS `platform`;
CREATE TABLE `platform`  (
  `platform_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '平台 appid',
  `platform_secret` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '平台 appsecret',
  `platform_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '平台 token',
  `platform_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '平台 消息解密key',
  `server_domain` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '服务器域名',
  `biz_domain` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '业务域名',
  `auth_redirect_url` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '用户授权成功回跳地址',
  `create_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0),
  PRIMARY KEY (`platform_id`) USING BTREE,
  UNIQUE INDEX `IX_APPID`(`platform_id`) USING BTREE COMMENT '平台唯一索引'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '平台注册信息主表' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
