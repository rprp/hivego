
CREATE TABLE scd_job (
 job_id integer NOT NULL ,/* '调度id',*/
  job_name varchar(128) NOT NULL ,/* '作业名称',*/
  job_desc varchar(500) DEFAULT NULL ,/* '作业说明',*/
  prev_job_id integer NOT NULL ,/* '上级作业id',*/
  next_job_id integer NOT NULL ,/* '下级作业id',*/
  create_user_id varchar(30) DEFAULT '' ,/* '创建人',*/
  create_time timestamp NULL DEFAULT NULL ,/* '创建时间',*/
  modify_user_id varchar(30) DEFAULT NULL ,/* '修改人',*/
  modify_time timestamp NULL DEFAULT NULL ,/* '修改时间',*/
  PRIMARY KEY (job_id)
);/*作业信息：\n           调度部分，记录调度作业信息。';*/



CREATE TABLE scd_job_log (
  batch_job_id varchar(128) NOT NULL ,/* '作业批次id，规则 批次id+作业id',*/
  batch_id varchar(128) NOT NULL ,/* '批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)',*/
  job_id integer NOT NULL ,/* '作业id',*/
  start_time timestamp NOT NULL ,/* '开始时间',*/
  end_time timestamp NOT NULL  ,/* '结束时间',*/
  state varchar(1) DEFAULT NULL ,/* '状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.意外中止',*/
  result real DEFAULT NULL ,/* '结果,作业中执行成功任务的百分比',*/
  batch_type varchar(1) NOT NULL ,/* '执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行',*/
  PRIMARY KEY (batch_job_id,job_id,start_time)
);/*='作业执行信息表：\n           日志部分，记录作业执行情况。';*/



CREATE TABLE scd_job_task (
  job_task_id integer NOT NULL ,/* '自增id',*/
  job_id integer NOT NULL ,/* '调度id',*/
  task_id integer NOT NULL ,/* '任务id',*/
  job_task_no integer NOT NULL ,/* '序号',*/
  create_user_id varchar(30) NOT NULL ,/* '创建人',*/
  create_time timestamp NOT NULL  ,/* '创建时间',*/
  PRIMARY KEY (job_task_id)
);/*='作业任务映射表：\n           调度部分，记录作业与任务映射信息。';*/



CREATE TABLE scd_schedule (
  scd_id integer NOT NULL ,/* '调度id',*/
  scd_name varchar(128) NOT NULL ,/* '调度名称',*/
  scd_num integer NOT NULL ,/* '调度次数 0.不限次数 ',*/
  scd_cyc varchar(2) NOT NULL ,/* '调度周期 ss 秒 mi 分钟 h 小时 d 日 m 月 w 周 q 季度 y 年',*/
  scd_timeout integer DEFAULT NULL ,/* '最大执行时间，单位 秒',*/
  scd_job_id integer DEFAULT NULL ,/* '作业id',*/
  scd_desc varchar(500) DEFAULT NULL ,/* '调度说明',*/
  create_user_id varchar(30) NOT NULL ,/* '创建人',*/
  create_time timestamp NOT NULL  ,/* '创建时间',*/
  modify_user_id varchar(30) DEFAULT NULL ,/* '修改人',*/
  modify_time timestamp NULL DEFAULT NULL ,/* '修改时间',*/
  PRIMARY KEY (scd_id)
);/*='调度信息：\n           调度部分，记录调度信息。';*/



CREATE TABLE scd_schedule_log (
  batch_id varchar(128) NOT NULL ,/* '批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)',*/
  scd_id integer NOT NULL ,/* '调度id',*/
  start_time timestamp NOT NULL ,/* '开始时间',*/
  end_time timestamp NOT NULL ,/* '结束时间',*/
  state varchar(1) DEFAULT NULL ,/* '状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.失败',*/
  result real DEFAULT NULL ,/* '结果,调度中执行成功任务的百分比',*/
  batch_type varchar(1) NOT NULL ,/* '执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行',*/
  PRIMARY KEY (batch_id,scd_id,start_time)
);/*='用户调度权限表：\n           日志部分，记录调度执行情况。';*/



CREATE TABLE scd_start (
  scd_id integer NOT NULL ,/* '调度id',*/
  scd_start integer NOT NULL ,/* '周期内启动时间单位秒',*/
  create_user_id integer NOT NULL ,/* '创建人',*/
  create_time timestamp NOT NULL  /* '创建时间'*/
);/*='异常处理设置信息：\n           调度部分，记录异常的处理信息。';*/



CREATE TABLE scd_task (
  task_id integer NOT NULL ,/* '任务id',*/
  task_address varchar(128) NOT NULL ,/* '任务地址',*/
  task_name varchar(128) NOT NULL ,/* '任务名称',*/
  task_cyc varchar(2) DEFAULT '' ,/* '调度周期 ss 秒 mi 分钟 h 小时 d 日 m 月 w 周 q 季度 y 年',*/
  task_time_out integer DEFAULT '0' ,/* '超时时间',*/
  task_start integer DEFAULT NULL ,/* '周期内启动时间，格式 mm-dd hh24:mi:ss，最大单位小于调度周期',*/
  task_type_id integer DEFAULT NULL ,/* '任务类型ID',*/
  task_cmd varchar(500) NOT NULL ,/* '任务命令行',*/
  task_desc varchar(500) DEFAULT NULL ,/* '任务说明',*/
  create_user_id varchar(30) DEFAULT '' ,/* '创建人',*/
  create_time timestamp NULL DEFAULT NULL ,/* '创建时间',*/
  modify_user_id varchar(30) DEFAULT NULL ,/* '修改人',*/
  modify_time timestamp NULL DEFAULT NULL ,/* '修改时间',*/
  PRIMARY KEY (task_id)
);/*='任务信息：\r           任务部分，任务信息记录需要执行的具体任务，以及执行方式。由用户录入。';*/



CREATE TABLE scd_task_attr (
  task_attr_id integer NOT NULL ,/* '自增id',*/
  task_id integer NOT NULL ,/* '任务id',*/
  task_attr_name varchar(500) NOT NULL ,/* '任务属性名称',*/
  task_attr_value text ,/* '任务属性值',*/
  create_time timestamp NOT NULL  ,/* '创建时间',*/
  PRIMARY KEY (task_attr_id)
);/*='任务属性表：\n           任务部分，记录具体任务的属性值。';*/



CREATE TABLE scd_task_log (
  batch_task_id varchar(128) NOT NULL ,/* '任务批次id，规则作业批次id+任务id',*/
  batch_job_id varchar(128) NOT NULL ,/* '作业批次id，规则 批次id+作业id',*/
  batch_id varchar(128) NOT NULL ,/* '批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)',*/
  task_id integer NOT NULL ,/* '任务id',*/
  start_time timestamp NOT NULL  ,/* '开始时间',*/
  end_time timestamp NOT NULL  ,/* '结束时间',*/
  state varchar(1) DEFAULT NULL ,/* '状态 0.初始状态 1. 执行中 2. 暂停 3. 完成 4.忽略 5.意外中止',*/
  batch_type varchar(1) NOT NULL ,/* '执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行',*/
  PRIMARY KEY (batch_task_id,task_id,start_time)
);/*='任务执行信息表：\n           日志部分，记录任务执行情况。';*/



CREATE TABLE scd_task_param (
  scd_param_id integer NOT NULL ,/* '自增id',*/
  task_id integer NOT NULL ,/* '任务id',*/
  scd_param_name varchar(128) NOT NULL ,/* '参数名称',*/
  scd_param_value varchar(254) DEFAULT NULL ,/* '参数值',*/
  create_user_id varchar(30) NOT NULL ,/* '创建人',*/
  create_time timestamp NOT NULL  ,/* '创建时间',*/
  PRIMARY KEY (scd_param_id)
);/*='调度参数信息：\n           调度部分，记录调度参数信息。任务id为空时表示为公共参数，即等同与调度';*/



CREATE TABLE scd_task_rel (
  task_rel_id integer NOT NULL ,/* '自增id',*/
  task_id integer NOT NULL ,/* '任务id',*/
  rel_task_id integer NOT NULL ,/* '依赖的任务id',*/
  create_user_id varchar(30) NOT NULL ,/* '创建人',*/
  create_time timestamp NOT NULL  ,/* '创建时间',*/
  PRIMARY KEY (task_rel_id)
);/*='任务依赖关系表：\n           记录任务之间依赖关系，也就是本作业中准备执行的任务与上级作业中任务的';*/



CREATE TABLE scd_task_type (
  task_type_id integer NOT NULL ,/* '自增id',*/
  task_type_name varchar(128) NOT NULL ,/* '任务类型名称',*/
  task_type_desc varchar(500) DEFAULT NULL ,/* '任务类型描述',*/
  create_time timestamp NOT NULL  ,/* '创建时间',*/
  PRIMARY KEY (task_type_id)
);/*='任务类型：\n           任务部分，记录具体任务的类型。';*/


CREATE TABLE scd_user (
  user_id varchar(30) NOT NULL ,/* 'hr用户编码',*/
  user_name varchar(128) NOT NULL ,/* '用户名称',*/
  user_mail varchar(128) NOT NULL ,/* '用户邮箱',*/
  user_password varchar(64) DEFAULT NULL ,/* '用户密码',*/
  user_phone varchar(64) DEFAULT NULL ,/* '用户手机号码',*/
  create_user_id varchar(30) NOT NULL ,/* '创建人',*/
  create_time timestamp NOT NULL  ,/* '创建时间',*/
  PRIMARY KEY (user_id)
);/*='用户信息表：\n           用户部分，记录用户信息。';*/



CREATE TABLE scd_user_schedule (
  user_schedule_id integer NOT NULL ,/* '自增id',*/
  scd_id integer NOT NULL ,/* '调度id',*/
  user_id varchar(30) NOT NULL ,/* 'hr用户编码',*/
  user_permission varchar(1) NOT NULL ,/* '权限类别 0 无权限 1.所有者  2 查看权限  ',*/
  create_user_id varchar(30) NOT NULL ,/* '创建人',*/
  create_time timestamp NOT NULL  ,/* '创建时间',*/
  PRIMARY KEY (user_schedule_id)
);/*='用户调度权限表：\n           用户部分，记录用户调度信息。';*/
