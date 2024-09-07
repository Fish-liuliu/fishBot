/*
 Navicat Premium Data Transfer

 Source Server         : wsl
 Source Server Type    : MySQL
 Source Server Version : 80039
 Source Host           : localhost:3306
 Source Schema         : fishbot

 Target Server Type    : MySQL
 Target Server Version : 80039
 File Encoding         : 65001

 Date: 07/09/2024 15:20:40
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for attendance
-- ----------------------------
DROP TABLE IF EXISTS `attendance`;
CREATE TABLE `attendance`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(255) NOT NULL COMMENT '成员id',
  `name` varchar(255)  NOT NULL COMMENT '成员名称',
  `attendance_date` date DEFAULT NULL COMMENT '打卡日期',
  `type` tinyint DEFAULT 0 COMMENT '0-未打卡，1-已打卡',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for random_numbers
-- ----------------------------
DROP TABLE IF EXISTS `random_numbers`;
CREATE TABLE `random_numbers`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `num` int DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
