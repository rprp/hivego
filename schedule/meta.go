//调度模块的数据结构
//package schedule
package main

import (
	"time"
)

//调度信息结构
type Schedule struct {
	id           int64           //调度ID
	name         string          //调度名称
	count        int8            //调度次数
	cyc          string          //调度周期
	startSecond  []time.Duration //周期内启动时间
	nextStart    time.Time       //周期内启动时间
	timeOut      int64           //最大执行时间
	jobId        int64           //作业ID
	job          *Job            //作业
	desc         string          //调度说明
	jobCnt       int64           //调度中作业数量
	taskCnt      int64           //调度中任务数量
	createUserId int64           //创建人
	createTime   time.Time       //创人
	modifyUserId int64           //修改人
	modifyTime   time.Time       //修改时间
}

//根据调度的周期及启动时间，按时将调度传至执行列表执行。
func (s *Schedule) Timer() { // {{{

	//获取距启动的时间（秒）
	countDown, err := getCountDown(s.cyc, s.startSecond)
	checkErr(err)

	s.nextStart = time.Now().Add(countDown)
	l.Infoln(s.name, "will start at", s.nextStart)
	select {
	case <-time.After(countDown):
		//刷新调度
		s.refreshSchedule()

		l.Infoln(s.name, "is start")
		//启动一个线程开始构建执行结构链
		es, err := NewExecSchedule(s)
		checkErr(err)
		//启动线程执行调度任务
		go es.Run()
	}
	return
} // }}}

//refreshSchedule方法用来从元数据库刷新调度信息
func (s *Schedule) refreshSchedule() { // {{{
	if ts, ok := getSchedule(s.id); ok {
		s.name = ts.name
		s.count = ts.count
		s.cyc = ts.cyc
		s.startSecond = ts.startSecond
		s.nextStart = ts.nextStart
		s.timeOut = ts.timeOut
		s.jobId = ts.jobId
		s.desc = ts.desc

		if tj, ok := getJob(s.jobId); ok {
			tj.refreshJob()
			s.job = tj
		}

		s.jobCnt = 0
		s.taskCnt = 0
		for j := s.job; j != nil; {
			s.jobCnt++
			s.taskCnt += j.taskCnt
			j = j.nextJob
		}
	}

} // }}}

//增加调度信息至元数据库
func (s *Schedule) Add() (err error) { // {{{
	s.SetNewId()
	sql := `INSERT INTO hive.scd_schedule
            (scd_id, scd_name, scd_num, scd_cyc,
             scd_timeout, scd_job_id, scd_desc, create_user_id,
             create_time, modify_user_id, modify_time)
		VALUES      (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = gDbConn.Exec(sql, &s.id, &s.name, &s.count, &s.cyc,
		&s.timeOut, &s.jobId, &s.desc, &s.createUserId, &s.createTime, &s.modifyUserId, &s.modifyTime)

	return err
} // }}}

//修改调度信息至元数据库
func (s *Schedule) Update() (err error) { // {{{
	sql := `UPDATE hive.scd_schedule 
		SET  scd_name=?,
             scd_num=?,
             scd_cyc=?,
             scd_timeout=?,
             scd_job_id=?,
             scd_desc=?,
             create_user_id=?,
             create_time=?,
             modify_user_id=?,
             modify_time=?
		 WHERE scd_id=?`
	_, err = gDbConn.Exec(sql, &s.name, &s.count, &s.cyc,
		&s.timeOut, &s.jobId, &s.desc, &s.createUserId, &s.createTime, &s.modifyUserId, &s.modifyTime, &s.id)

	return err
} // }}}

//删除元数据库中的调度信息
func (s *Schedule) Delete() (err error) { // {{{
	sql := `Delete hive.scd_schedule WHERE scd_id=?`
	_, err = gDbConn.Exec(sql, &s.id)

	return err
} // }}}

//设置调度下的Job
func (s *Schedule) SetJob(jobid int64) (err error) { // {{{
	s.jobId = jobid
	if j, ok := getJob(jobid); ok {
		s.job = j
	}

	return nil
} // }}}

//获取新Id
func (s *Schedule) SetNewId() (err error) { // {{{
	var id int64

	//查询全部schedule列表
	sql := `SELECT max(scd.scd_id) as scd_id
			FROM hive.scd_schedule scd`

	rows, err := gDbConn.Query(sql)
	checkErr(err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
	}
	s.id = id + 1

	return err

} // }}}

//从元数据库获取指定Schedule的启动时间。
func getStart(id int64) (st []time.Duration, ok bool) { // {{{

	st = make([]time.Duration, 0)

	//查询全部schedule启动时间列表
	sql := `SELECT s.scd_start
			FROM hive.scd_start s
			WHERE s.scd_id=?`

	rows, err := gDbConn.Query(sql, id)
	checkErr(err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		var td int64
		err = rows.Scan(&td)
		if err == nil {
			ok = true
		}
		st = append(st, time.Duration(td)*time.Second)

	}

	if len(st) == 0 {
		st = append(st, time.Duration(0))
	}

	return st, ok
} // }}}

//从元数据库获取指定的Schedule。
func getSchedule(id int64) (scd *Schedule, ok bool) { // {{{

	//查询全部schedule列表
	sql := `SELECT scd.scd_id,
				scd.scd_name,
				scd.scd_num,
				scd.scd_cyc,
				scd.scd_timeout,
				scd.scd_job_id,
				scd.scd_desc
			FROM hive.scd_schedule scd
			WHERE scd.scd_id=?`

	rows, err := gDbConn.Query(sql, id)
	checkErr(err)

	scd = &Schedule{}
	scd.startSecond = make([]time.Duration, 0)
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&scd.id, &scd.name, &scd.count, &scd.cyc,
			&scd.timeOut, &scd.jobId, &scd.desc)
		if err == nil {
			ok = true
		}
		scd.startSecond, _ = getStart(scd.id)

	}

	return scd, ok
} // }}}

//从元数据库获取Schedule列表。
func getAllSchedules() (scds map[int64]*Schedule, err error) { // {{{
	scds = make(map[int64]*Schedule)

	//查询全部schedule列表
	sql := `SELECT scd.scd_id,
				scd.scd_name,
				scd.scd_num,
				scd.scd_cyc,
				scd.scd_timeout,
				scd.scd_job_id,
				scd.scd_desc
			FROM hive.scd_schedule scd`

	rows, err := gDbConn.Query(sql)
	checkErr(err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		scd := &Schedule{}
		scd.startSecond = make([]time.Duration, 0)
		err = rows.Scan(&scd.id, &scd.name, &scd.count, &scd.cyc,
			&scd.timeOut, &scd.jobId, &scd.desc)
		scd.startSecond, _ = getStart(scd.id)

		scds[scd.id] = scd
	}

	return scds, err
} // }}}
