-- MySQL dump 10.13  Distrib 5.6.14, for osx10.7 (x86_64)
--
-- Host: localhost    Database: hive
-- ------------------------------------------------------
-- Server version	5.6.14

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `scd_deploy`
--

DROP TABLE IF EXISTS `scd_deploy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_deploy` (
  `scd_addr` varchar(128) NOT NULL COMMENT '调度地址',
  `scd_node_type` bigint(20) NOT NULL COMMENT '节点类型 0.主节点 1.子节点',
  `create_user_id` varchar(30) NOT NULL COMMENT '创建人',
  `create_time` date NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`scd_addr`,`scd_node_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='调度部署信息表：\n           记录调度程序部署情况。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_deploy`
--

LOCK TABLES `scd_deploy` WRITE;
/*!40000 ALTER TABLE `scd_deploy` DISABLE KEYS */;
/*!40000 ALTER TABLE `scd_deploy` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_job`
--

DROP TABLE IF EXISTS `scd_job`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_job` (
  `job_id` bigint(20) NOT NULL COMMENT '调度id',
  `job_name` varchar(256) NOT NULL COMMENT '作业名称',
  `job_desc` varchar(500) DEFAULT NULL COMMENT '作业说明',
  `prev_job_id` bigint(20) NOT NULL COMMENT '上级作业id',
  `next_job_id` bigint(20) NOT NULL COMMENT '下级作业id',
  `create_user_id` varchar(30) DEFAULT '' COMMENT '创建人',
  `create_time` date DEFAULT NULL COMMENT '创建时间',
  `modify_user_id` varchar(30) DEFAULT NULL COMMENT '修改人',
  `modify_time` date DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`job_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='作业信息：\n           调度部分，记录调度作业信息。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_job`
--

LOCK TABLES `scd_job` WRITE;
/*!40000 ALTER TABLE `scd_job` DISABLE KEYS */;
INSERT INTO `scd_job` VALUES (1,'作业1','0',0,2,'',NULL,NULL,NULL),(2,'作业2','0',1,3,'',NULL,NULL,NULL),(3,'作业3','0',2,9,'',NULL,NULL,NULL),(4,'作业4','0',0,5,'',NULL,NULL,NULL),(5,'作业5','0',5,6,'',NULL,NULL,NULL),(6,'作业6','0',6,7,'',NULL,NULL,NULL),(7,'作业7','0',7,8,'',NULL,NULL,NULL),(8,'作业8','0',8,0,'',NULL,NULL,NULL),(9,'作业9','0',3,10,'',NULL,NULL,NULL),(10,'作业10','0',9,0,'',NULL,NULL,NULL);
/*!40000 ALTER TABLE `scd_job` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_job_log`
--

DROP TABLE IF EXISTS `scd_job_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_job_log` (
  `batch_job_id` varchar(128) NOT NULL COMMENT '作业批次id，规则 批次id+作业id',
  `batch_id` varchar(128) NOT NULL COMMENT '批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)',
  `job_id` bigint(20) NOT NULL COMMENT '作业id',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '结束时间',
  `state` varchar(1) DEFAULT NULL COMMENT '状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.意外中止',
  `result` decimal(10,2) DEFAULT NULL COMMENT '结果,作业中执行成功任务的百分比',
  `batch_type` varchar(1) NOT NULL COMMENT '执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行',
  PRIMARY KEY (`batch_job_id`,`job_id`,`start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='作业执行信息表：\n           日志部分，记录作业执行情况。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_job_log`
--

LOCK TABLES `scd_job_log` WRITE;
/*!40000 ALTER TABLE `scd_job_log` DISABLE KEYS */;
INSERT INTO `scd_job_log` VALUES ('2014-06-16 09:48:00.047067 1.1','2014-06-16 09:48:00.047067 1',1,'2014-06-16 01:48:00','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:48:00.047067 1.10','2014-06-16 09:48:00.047067 1',10,'2014-06-16 01:48:40','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:48:00.047067 1.2','2014-06-16 09:48:00.047067 1',2,'2014-06-16 01:48:00','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:48:00.047067 1.3','2014-06-16 09:48:00.047067 1',3,'2014-06-16 01:48:20','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:48:00.047067 1.9','2014-06-16 09:48:00.047067 1',9,'2014-06-16 01:48:30','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:49:00.039637 1.1','2014-06-16 09:49:00.039637 1',1,'2014-06-16 01:49:00','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:49:00.039637 1.10','2014-06-16 09:49:00.039637 1',10,'0000-00-00 00:00:00','0000-00-00 00:00:00','0',0.00,'1'),('2014-06-16 09:49:00.039637 1.2','2014-06-16 09:49:00.039637 1',2,'2014-06-16 01:49:00','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:49:00.039637 1.3','2014-06-16 09:49:00.039637 1',3,'0000-00-00 00:00:00','0000-00-00 00:00:00','0',0.00,'1'),('2014-06-16 09:49:00.039637 1.9','2014-06-16 09:49:00.039637 1',9,'0000-00-00 00:00:00','0000-00-00 00:00:00','0',0.00,'1'),('2014-06-16 09:50:00.043007 1.1','2014-06-16 09:50:00.043007 1',1,'2014-06-16 01:50:00','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:50:00.043007 1.10','2014-06-16 09:50:00.043007 1',10,'2014-06-16 01:50:40','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:50:00.043007 1.2','2014-06-16 09:50:00.043007 1',2,'2014-06-16 01:50:00','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:50:00.043007 1.3','2014-06-16 09:50:00.043007 1',3,'2014-06-16 01:50:20','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:50:00.043007 1.9','2014-06-16 09:50:00.043007 1',9,'2014-06-16 01:50:30','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:51:00.041106 1.1','2014-06-16 09:51:00.041106 1',1,'2014-06-16 01:51:00','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:51:00.041106 1.10','2014-06-16 09:51:00.041106 1',10,'2014-06-16 01:51:40','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:51:00.041106 1.2','2014-06-16 09:51:00.041106 1',2,'2014-06-16 01:51:00','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:51:00.041106 1.3','2014-06-16 09:51:00.041106 1',3,'2014-06-16 01:51:20','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:51:00.041106 1.9','2014-06-16 09:51:00.041106 1',9,'2014-06-16 01:51:30','0000-00-00 00:00:00','1',0.00,'1');
/*!40000 ALTER TABLE `scd_job_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_job_task`
--

DROP TABLE IF EXISTS `scd_job_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_job_task` (
  `job_task_id` bigint(20) NOT NULL COMMENT '自增id',
  `job_id` bigint(20) NOT NULL COMMENT '调度id',
  `task_id` bigint(20) NOT NULL COMMENT '任务id',
  `job_task_no` bigint(20) NOT NULL COMMENT '序号',
  `create_user_id` varchar(30) NOT NULL COMMENT '创建人',
  `create_time` date NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`job_task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='作业任务映射表：\n           调度部分，记录作业与任务映射信息。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_job_task`
--

LOCK TABLES `scd_job_task` WRITE;
/*!40000 ALTER TABLE `scd_job_task` DISABLE KEYS */;
INSERT INTO `scd_job_task` VALUES (1,1,1,1,'1','2014-05-28'),(2,1,2,2,'1','2014-05-28'),(3,2,3,1,'1','2014-05-28'),(4,2,4,2,'1','2014-05-28'),(5,2,5,3,'1','2014-05-28'),(6,3,6,1,'1','2014-05-28'),(7,9,7,1,'1','2014-05-28'),(8,9,8,2,'1','2014-05-28'),(9,4,9,1,'1','2014-05-28'),(10,4,10,2,'1','2014-05-28'),(11,4,11,3,'1','2014-05-28'),(12,5,12,1,'1','2014-05-28'),(13,6,13,1,'1','2014-05-28'),(14,7,14,1,'1','2014-05-28'),(15,7,15,2,'1','2014-05-28'),(16,7,16,3,'1','2014-05-28'),(17,8,17,1,'1','2014-05-28'),(18,8,18,2,'1','2014-05-28'),(19,8,19,3,'1','2014-05-28'),(20,10,20,1,'1','2014-05-28');
/*!40000 ALTER TABLE `scd_job_task` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_schedule`
--

DROP TABLE IF EXISTS `scd_schedule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_schedule` (
  `scd_id` bigint(20) NOT NULL COMMENT '调度id',
  `scd_name` varchar(256) NOT NULL COMMENT '调度名称',
  `scd_num` int(11) NOT NULL COMMENT '调度次数 0.不限次数 ',
  `scd_cyc` varchar(2) NOT NULL COMMENT '调度周期 ss 秒 mi 分钟 h 小时 d 日 m 月 w 周 q 季度 y 年',
  `scd_timeout` bigint(20) DEFAULT NULL COMMENT '最大执行时间，单位 秒',
  `scd_job_id` bigint(20) DEFAULT NULL COMMENT '作业id',
  `scd_desc` varchar(500) DEFAULT NULL COMMENT '调度说明',
  `create_user_id` varchar(30) NOT NULL COMMENT '创建人',
  `create_time` date NOT NULL COMMENT '创建时间',
  `modify_user_id` varchar(30) DEFAULT NULL COMMENT '修改人',
  `modify_time` date DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`scd_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='调度信息：\n           调度部分，记录调度信息。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_schedule`
--

LOCK TABLES `scd_schedule` WRITE;
/*!40000 ALTER TABLE `scd_schedule` DISABLE KEYS */;
INSERT INTO `scd_schedule` VALUES (1,'数据仓库调度',0,'mi',0,1,'数据仓库日常调度','1','2014-05-28','1','2014-05-28'),(2,'数据市场调度',0,'h',0,4,'数据市场日常调度','1','2014-05-28','1','2014-05-28');
/*!40000 ALTER TABLE `scd_schedule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_schedule_log`
--

DROP TABLE IF EXISTS `scd_schedule_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_schedule_log` (
  `batch_id` varchar(128) NOT NULL COMMENT '批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)',
  `scd_id` bigint(20) NOT NULL COMMENT '调度id',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL COMMENT '结束时间',
  `state` varchar(1) DEFAULT NULL COMMENT '状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.失败',
  `result` decimal(10,2) DEFAULT NULL COMMENT '结果,调度中执行成功任务的百分比',
  `batch_type` varchar(1) NOT NULL COMMENT '执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行',
  PRIMARY KEY (`batch_id`,`scd_id`,`start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户调度权限表：\n           日志部分，记录调度执行情况。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_schedule_log`
--

LOCK TABLES `scd_schedule_log` WRITE;
/*!40000 ALTER TABLE `scd_schedule_log` DISABLE KEYS */;
INSERT INTO `scd_schedule_log` VALUES ('2014-06-16 09:48:00.047067 1',1,'2014-06-16 01:48:00','2014-06-16 01:48:50','3',1.00,'1'),('2014-06-16 09:49:00.039637 1',1,'2014-06-16 01:49:00','0000-00-00 00:00:00','1',0.00,'1'),('2014-06-16 09:50:00.043007 1',1,'2014-06-16 01:50:00','2014-06-16 01:50:50','3',1.00,'1'),('2014-06-16 09:51:00.041106 1',1,'2014-06-16 01:51:00','2014-06-16 01:51:50','3',1.00,'1');
/*!40000 ALTER TABLE `scd_schedule_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_start`
--

DROP TABLE IF EXISTS `scd_start`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_start` (
  `scd_id` bigint(20) NOT NULL COMMENT '调度id',
  `scd_start` bigint(20) NOT NULL COMMENT '周期内启动时间单位秒',
  `create_user_id` bigint(20) NOT NULL COMMENT '创建人',
  `create_time` date NOT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='异常处理设置信息：\n           调度部分，记录异常的处理信息。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_start`
--

LOCK TABLES `scd_start` WRITE;
/*!40000 ALTER TABLE `scd_start` DISABLE KEYS */;
INSERT INTO `scd_start` VALUES (2,60,1,'2014-06-12'),(2,2500,1,'2014-06-12'),(2,2700,1,'2014-06-12'),(2,2400,1,'2014-06-12'),(2,600,1,'2014-06-12');
/*!40000 ALTER TABLE `scd_start` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_task`
--

DROP TABLE IF EXISTS `scd_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_task` (
  `task_id` bigint(20) NOT NULL COMMENT '任务id',
  `task_address` varchar(256) NOT NULL COMMENT '任务地址',
  `task_name` varchar(256) NOT NULL COMMENT '任务名称',
  `task_cyc` varchar(2) DEFAULT '' COMMENT '调度周期 ss 秒 mi 分钟 h 小时 d 日 m 月 w 周 q 季度 y 年',
  `task_time_out` bigint(20) DEFAULT '0' COMMENT '超时时间',
  `task_start` bigint(20) DEFAULT NULL COMMENT '周期内启动时间，格式 mm-dd hh24:mi:ss，最大单位小于调度周期',
  `task_type_id` bigint(20) DEFAULT NULL COMMENT '任务类型ID',
  `task_cmd` varchar(500) NOT NULL COMMENT '任务命令行',
  `task_desc` varchar(500) DEFAULT NULL COMMENT '任务说明',
  `create_user_id` varchar(30) DEFAULT '' COMMENT '创建人',
  `create_time` date DEFAULT NULL COMMENT '创建时间',
  `modify_user_id` varchar(30) DEFAULT NULL COMMENT '修改人',
  `modify_time` date DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务信息：\r           任务部分，任务信息记录需要执行的具体任务，以及执行方式。由用户录入。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_task`
--

LOCK TABLES `scd_task` WRITE;
/*!40000 ALTER TABLE `scd_task` DISABLE KEYS */;
INSERT INTO `scd_task` VALUES (1,'127.0.0.1','任务1','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py',NULL,'1','2014-05-28',NULL,NULL),(2,'127.0.0.1','ping2','h',60,2950,1,'ping',NULL,'1','2014-05-28',NULL,NULL),(3,'127.0.0.1','任务3','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py',NULL,'1','2014-05-28',NULL,NULL),(4,'127.0.0.1','任务4','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py',NULL,'1','2014-05-28',NULL,NULL),(5,'127.0.0.1','任务5','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py',NULL,'1','2014-05-28',NULL,NULL),(6,'127.0.0.1','任务6','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(7,'127.0.0.1','任务7','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(8,'127.0.0.1','任务8','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(9,'127.0.0.1','任务9','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(10,'127.0.0.1','任务10','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(11,'127.0.0.1','任务11','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(12,'127.0.0.1','任务12','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(13,'127.0.0.1','任务13','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(14,'127.0.0.1','任务14','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(15,'127.0.0.1','任务15','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(16,'127.0.0.1','任务16','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(17,'127.0.0.1','任务17','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(18,'127.0.0.1','任务18','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(19,'127.0.0.1','任务19','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL),(20,'127.0.0.1','任务20','',60,0,1,'/Users/rp/develop/code/py/testSchedule.py\n',NULL,'1','2014-05-28',NULL,NULL);
/*!40000 ALTER TABLE `scd_task` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_task_attr`
--

DROP TABLE IF EXISTS `scd_task_attr`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_task_attr` (
  `task_attr_id` bigint(20) NOT NULL COMMENT '自增id',
  `task_id` bigint(20) NOT NULL COMMENT '任务id',
  `task_attr_name` varchar(500) NOT NULL COMMENT '任务属性名称',
  `task_attr_value` text COMMENT '任务属性值',
  `create_time` date NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`task_attr_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务属性表：\n           任务部分，记录具体任务的属性值。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_task_attr`
--

LOCK TABLES `scd_task_attr` WRITE;
/*!40000 ALTER TABLE `scd_task_attr` DISABLE KEYS */;
INSERT INTO `scd_task_attr` VALUES (1,2,'name','abc','2014-06-06'),(2,2,'type','cc','2014-06-06'),(3,2,'time','dd','2014-06-06');
/*!40000 ALTER TABLE `scd_task_attr` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_task_exception`
--

DROP TABLE IF EXISTS `scd_task_exception`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_task_exception` (
  `task_exception_id` bigint(20) NOT NULL COMMENT '自增id',
  `task_exception_no` bigint(20) NOT NULL COMMENT '序号',
  `scd_id` bigint(20) NOT NULL COMMENT '调度id',
  `job_id` bigint(20) NOT NULL COMMENT '作业id',
  `task_id` bigint(20) NOT NULL COMMENT '任务id',
  `exception_id` bigint(20) NOT NULL COMMENT '异常处理id',
  `create_user_id` varchar(30) NOT NULL COMMENT '创建人',
  `create_time` date NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`task_exception_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务异常关系表：\n           调度部分，记录任务对应异常的处理操作。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_task_exception`
--

LOCK TABLES `scd_task_exception` WRITE;
/*!40000 ALTER TABLE `scd_task_exception` DISABLE KEYS */;
/*!40000 ALTER TABLE `scd_task_exception` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_task_log`
--

DROP TABLE IF EXISTS `scd_task_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_task_log` (
  `batch_task_id` varchar(128) NOT NULL COMMENT '任务批次id，规则作业批次id+任务id',
  `batch_job_id` varchar(128) NOT NULL COMMENT '作业批次id，规则 批次id+作业id',
  `batch_id` varchar(128) NOT NULL COMMENT '批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)',
  `task_id` bigint(20) NOT NULL COMMENT '任务id',
  `start_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '开始时间',
  `end_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '结束时间',
  `state` varchar(1) DEFAULT NULL COMMENT '状态 0.初始状态 1. 执行中 2. 暂停 3. 完成 4.忽略 5.意外中止',
  `batch_type` varchar(1) NOT NULL COMMENT '执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行',
  PRIMARY KEY (`batch_task_id`,`task_id`,`start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务执行信息表：\n           日志部分，记录任务执行情况。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_task_log`
--

LOCK TABLES `scd_task_log` WRITE;
/*!40000 ALTER TABLE `scd_task_log` DISABLE KEYS */;
INSERT INTO `scd_task_log` VALUES ('2014-06-16 09:48:00.047067 1.1.1','2014-06-16 09:48:00.047067 1.1','2014-06-16 09:48:00.047067 1',1,'2014-06-16 01:48:00','2014-06-16 01:48:10','3','1'),('2014-06-16 09:48:00.047067 1.1.2','2014-06-16 09:48:00.047067 1.1','2014-06-16 09:48:00.047067 1',2,'2014-06-16 01:48:00','2014-06-16 01:48:00','4','1'),('2014-06-16 09:48:00.047067 1.10.20','2014-06-16 09:48:00.047067 1.10','2014-06-16 09:48:00.047067 1',20,'2014-06-16 01:48:40','2014-06-16 01:48:50','3','1'),('2014-06-16 09:48:00.047067 1.2.3','2014-06-16 09:48:00.047067 1.2','2014-06-16 09:48:00.047067 1',3,'2014-06-16 01:48:10','2014-06-16 01:48:20','3','1'),('2014-06-16 09:48:00.047067 1.2.4','2014-06-16 09:48:00.047067 1.2','2014-06-16 09:48:00.047067 1',4,'2014-06-16 01:48:10','2014-06-16 01:48:20','3','1'),('2014-06-16 09:48:00.047067 1.2.5','2014-06-16 09:48:00.047067 1.2','2014-06-16 09:48:00.047067 1',5,'2014-06-16 01:48:00','2014-06-16 01:48:10','3','1'),('2014-06-16 09:48:00.047067 1.3.6','2014-06-16 09:48:00.047067 1.3','2014-06-16 09:48:00.047067 1',6,'2014-06-16 01:48:20','2014-06-16 01:48:30','3','1'),('2014-06-16 09:48:00.047067 1.9.7','2014-06-16 09:48:00.047067 1.9','2014-06-16 09:48:00.047067 1',7,'2014-06-16 01:48:30','2014-06-16 01:48:40','3','1'),('2014-06-16 09:48:00.047067 1.9.8','2014-06-16 09:48:00.047067 1.9','2014-06-16 09:48:00.047067 1',8,'2014-06-16 01:48:30','2014-06-16 01:48:40','3','1'),('2014-06-16 09:49:00.039637 1.1.1','2014-06-16 09:49:00.039637 1.1','2014-06-16 09:49:00.039637 1',1,'2014-06-16 01:49:00','2014-06-16 01:49:10','3','1'),('2014-06-16 09:49:00.039637 1.1.2','2014-06-16 09:49:00.039637 1.1','2014-06-16 09:49:00.039637 1',2,'2014-06-16 01:49:00','2014-06-16 01:49:05','3','1'),('2014-06-16 09:49:00.039637 1.10.20','2014-06-16 09:49:00.039637 1.10','2014-06-16 09:49:00.039637 1',20,'0000-00-00 00:00:00','0000-00-00 00:00:00','0','1'),('2014-06-16 09:49:00.039637 1.2.3','2014-06-16 09:49:00.039637 1.2','2014-06-16 09:49:00.039637 1',3,'2014-06-16 01:49:10','0000-00-00 00:00:00','1','1'),('2014-06-16 09:49:00.039637 1.2.4','2014-06-16 09:49:00.039637 1.2','2014-06-16 09:49:00.039637 1',4,'2014-06-16 01:49:10','0000-00-00 00:00:00','1','1'),('2014-06-16 09:49:00.039637 1.2.5','2014-06-16 09:49:00.039637 1.2','2014-06-16 09:49:00.039637 1',5,'2014-06-16 01:49:00','2014-06-16 01:49:10','3','1'),('2014-06-16 09:49:00.039637 1.3.6','2014-06-16 09:49:00.039637 1.3','2014-06-16 09:49:00.039637 1',6,'0000-00-00 00:00:00','0000-00-00 00:00:00','0','1'),('2014-06-16 09:49:00.039637 1.9.7','2014-06-16 09:49:00.039637 1.9','2014-06-16 09:49:00.039637 1',7,'0000-00-00 00:00:00','0000-00-00 00:00:00','0','1'),('2014-06-16 09:49:00.039637 1.9.8','2014-06-16 09:49:00.039637 1.9','2014-06-16 09:49:00.039637 1',8,'0000-00-00 00:00:00','0000-00-00 00:00:00','0','1'),('2014-06-16 09:50:00.043007 1.1.1','2014-06-16 09:50:00.043007 1.1','2014-06-16 09:50:00.043007 1',1,'2014-06-16 01:50:00','2014-06-16 01:50:10','3','1'),('2014-06-16 09:50:00.043007 1.1.2','2014-06-16 09:50:00.043007 1.1','2014-06-16 09:50:00.043007 1',2,'2014-06-16 01:50:00','2014-06-16 01:50:00','4','1'),('2014-06-16 09:50:00.043007 1.10.20','2014-06-16 09:50:00.043007 1.10','2014-06-16 09:50:00.043007 1',20,'2014-06-16 01:50:40','2014-06-16 01:50:50','3','1'),('2014-06-16 09:50:00.043007 1.2.3','2014-06-16 09:50:00.043007 1.2','2014-06-16 09:50:00.043007 1',3,'2014-06-16 01:50:10','2014-06-16 01:50:20','3','1'),('2014-06-16 09:50:00.043007 1.2.4','2014-06-16 09:50:00.043007 1.2','2014-06-16 09:50:00.043007 1',4,'2014-06-16 01:50:10','2014-06-16 01:50:20','3','1'),('2014-06-16 09:50:00.043007 1.2.5','2014-06-16 09:50:00.043007 1.2','2014-06-16 09:50:00.043007 1',5,'2014-06-16 01:50:00','2014-06-16 01:50:10','3','1'),('2014-06-16 09:50:00.043007 1.3.6','2014-06-16 09:50:00.043007 1.3','2014-06-16 09:50:00.043007 1',6,'2014-06-16 01:50:20','2014-06-16 01:50:30','3','1'),('2014-06-16 09:50:00.043007 1.9.7','2014-06-16 09:50:00.043007 1.9','2014-06-16 09:50:00.043007 1',7,'2014-06-16 01:50:30','2014-06-16 01:50:40','3','1'),('2014-06-16 09:50:00.043007 1.9.8','2014-06-16 09:50:00.043007 1.9','2014-06-16 09:50:00.043007 1',8,'2014-06-16 01:50:30','2014-06-16 01:50:40','3','1'),('2014-06-16 09:51:00.041106 1.1.1','2014-06-16 09:51:00.041106 1.1','2014-06-16 09:51:00.041106 1',1,'2014-06-16 01:51:00','2014-06-16 01:51:10','3','1'),('2014-06-16 09:51:00.041106 1.1.2','2014-06-16 09:51:00.041106 1.1','2014-06-16 09:51:00.041106 1',2,'2014-06-16 01:51:00','2014-06-16 01:51:00','4','1'),('2014-06-16 09:51:00.041106 1.10.20','2014-06-16 09:51:00.041106 1.10','2014-06-16 09:51:00.041106 1',20,'2014-06-16 01:51:40','2014-06-16 01:51:50','3','1'),('2014-06-16 09:51:00.041106 1.2.3','2014-06-16 09:51:00.041106 1.2','2014-06-16 09:51:00.041106 1',3,'2014-06-16 01:51:10','2014-06-16 01:51:20','3','1'),('2014-06-16 09:51:00.041106 1.2.4','2014-06-16 09:51:00.041106 1.2','2014-06-16 09:51:00.041106 1',4,'2014-06-16 01:51:10','2014-06-16 01:51:20','3','1'),('2014-06-16 09:51:00.041106 1.2.5','2014-06-16 09:51:00.041106 1.2','2014-06-16 09:51:00.041106 1',5,'2014-06-16 01:51:00','2014-06-16 01:51:10','3','1'),('2014-06-16 09:51:00.041106 1.3.6','2014-06-16 09:51:00.041106 1.3','2014-06-16 09:51:00.041106 1',6,'2014-06-16 01:51:20','2014-06-16 01:51:30','3','1'),('2014-06-16 09:51:00.041106 1.9.7','2014-06-16 09:51:00.041106 1.9','2014-06-16 09:51:00.041106 1',7,'2014-06-16 01:51:30','2014-06-16 01:51:40','3','1'),('2014-06-16 09:51:00.041106 1.9.8','2014-06-16 09:51:00.041106 1.9','2014-06-16 09:51:00.041106 1',8,'2014-06-16 01:51:30','2014-06-16 01:51:40','3','1');
/*!40000 ALTER TABLE `scd_task_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_task_param`
--

DROP TABLE IF EXISTS `scd_task_param`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_task_param` (
  `scd_param_id` bigint(20) NOT NULL COMMENT '自增id',
  `task_id` bigint(20) NOT NULL COMMENT '任务id',
  `scd_param_name` varchar(256) NOT NULL COMMENT '参数名称',
  `scd_param_value` varchar(256) DEFAULT NULL COMMENT '参数值',
  `create_user_id` varchar(30) NOT NULL COMMENT '创建人',
  `create_time` date NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`scd_param_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='调度参数信息：\n           调度部分，记录调度参数信息。任务id为空时表示为公共参数，即等同与调度';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_task_param`
--

LOCK TABLES `scd_task_param` WRITE;
/*!40000 ALTER TABLE `scd_task_param` DISABLE KEYS */;
INSERT INTO `scd_task_param` VALUES (1,2,'time','-t 5','1','2014-06-09'),(2,2,'ip','localhost','1','2014-06-09');
/*!40000 ALTER TABLE `scd_task_param` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_task_rel`
--

DROP TABLE IF EXISTS `scd_task_rel`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_task_rel` (
  `task_rel_id` bigint(20) NOT NULL COMMENT '自增id',
  `task_id` bigint(20) NOT NULL COMMENT '任务id',
  `rel_task_id` bigint(20) NOT NULL COMMENT '依赖的任务id',
  `create_user_id` varchar(30) NOT NULL COMMENT '创建人',
  `create_time` date NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`task_rel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务依赖关系表：\n           记录任务之间依赖关系，也就是本作业中准备执行的任务与上级作业中任务的';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_task_rel`
--

LOCK TABLES `scd_task_rel` WRITE;
/*!40000 ALTER TABLE `scd_task_rel` DISABLE KEYS */;
INSERT INTO `scd_task_rel` VALUES (1,3,1,'1','2014-05-29'),(2,3,2,'1','2014-05-29'),(3,4,1,'1','2014-05-29'),(4,4,2,'1','2014-05-29'),(5,6,3,'1','2014-05-29'),(6,6,5,'1','2014-05-29'),(7,7,4,'1','2014-05-29'),(8,7,6,'1','2014-05-29'),(9,8,6,'1','2014-05-29'),(10,20,7,'1','2014-05-29'),(11,20,8,'1','2014-05-29'),(12,12,9,'1','2014-05-29'),(13,12,10,'1','2014-05-29'),(14,13,12,'1','2014-05-29'),(15,14,13,'1','2014-05-29'),(16,15,13,'1','2014-05-29'),(17,16,11,'1','2014-05-29'),(18,16,13,'1','2014-05-29'),(19,17,14,'1','2014-05-29'),(20,17,15,'1','2014-05-29'),(21,17,16,'1','2014-05-29'),(22,18,15,'1','2014-05-29'),(23,18,16,'1','2014-05-29'),(24,19,15,'1','2014-05-29');
/*!40000 ALTER TABLE `scd_task_rel` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_task_type`
--

DROP TABLE IF EXISTS `scd_task_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_task_type` (
  `task_type_id` bigint(20) NOT NULL COMMENT '自增id',
  `task_type_name` varchar(256) NOT NULL COMMENT '任务类型名称',
  `task_type_desc` varchar(500) DEFAULT NULL COMMENT '任务类型描述',
  `create_time` date NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`task_type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务类型：\n           任务部分，记录具体任务的类型。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_task_type`
--

LOCK TABLES `scd_task_type` WRITE;
/*!40000 ALTER TABLE `scd_task_type` DISABLE KEYS */;
/*!40000 ALTER TABLE `scd_task_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_user`
--

DROP TABLE IF EXISTS `scd_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_user` (
  `user_id` varchar(30) NOT NULL COMMENT 'hr用户编码',
  `user_name` varchar(128) NOT NULL COMMENT '用户名称',
  `user_mail` varchar(256) NOT NULL COMMENT '用户邮箱',
  `user_password` varchar(64) DEFAULT NULL COMMENT '用户密码',
  `user_phone` varchar(64) DEFAULT NULL COMMENT '用户手机号码',
  `create_user_id` varchar(30) NOT NULL COMMENT '创建人',
  `create_time` date NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户信息表：\n           用户部分，记录用户信息。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_user`
--

LOCK TABLES `scd_user` WRITE;
/*!40000 ALTER TABLE `scd_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `scd_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_user_job`
--

DROP TABLE IF EXISTS `scd_user_job`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_user_job` (
  `user_job_id` bigint(20) NOT NULL COMMENT '自增id',
  `job_id` bigint(20) NOT NULL COMMENT '作业id',
  `user_id` varchar(30) NOT NULL COMMENT 'hr用户编码',
  `user_permission` varchar(1) NOT NULL COMMENT '权限类别 0 无权限 1.所有者  2 查看权限  ',
  `create_user_id` varchar(30) NOT NULL COMMENT '创建人',
  `create_time` date NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`user_job_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户作业权限表：\n           用户部分，记录用户作业信息。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_user_job`
--

LOCK TABLES `scd_user_job` WRITE;
/*!40000 ALTER TABLE `scd_user_job` DISABLE KEYS */;
/*!40000 ALTER TABLE `scd_user_job` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scd_user_schedule`
--

DROP TABLE IF EXISTS `scd_user_schedule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scd_user_schedule` (
  `user_schedule_id` bigint(20) NOT NULL COMMENT '自增id',
  `scd_id` bigint(20) NOT NULL COMMENT '调度id',
  `user_id` varchar(30) NOT NULL COMMENT 'hr用户编码',
  `user_permission` varchar(1) NOT NULL COMMENT '权限类别 0 无权限 1.所有者  2 查看权限  ',
  `create_user_id` varchar(30) NOT NULL COMMENT '创建人',
  `create_time` date NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`user_schedule_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户调度权限表：\n           用户部分，记录用户调度信息。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scd_user_schedule`
--

LOCK TABLES `scd_user_schedule` WRITE;
/*!40000 ALTER TABLE `scd_user_schedule` DISABLE KEYS */;
/*!40000 ALTER TABLE `scd_user_schedule` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2014-06-16  9:53:47
