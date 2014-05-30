//调度模块的数据结构
//package schedule
package main

import (
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//调度信息结构
type Schedule struct {
	id          int64             //调度ID
	name        string            //调度名称
	count       int8              //调度次数
	cyc         string            //调度周期
	param       map[string]string //调度参数
	startSecond time.Duration     //周期内启动时间
	nextStart   time.Time         //周期内启动时间
	timeOut     int64             //最大执行时间
	jobId       int64             //作业ID
	job         *Job              //作业
	desc        string            //调度说明
}

//根据调度的周期及启动时间，按时将调度传至执行列表执行。
func (s *Schedule) Timer() {
	//获取距启动的时间（秒）
	countDown, err := getCountDown(s.cyc, s.startSecond)
	checkErr(err)

	s.nextStart = time.Now().Add(countDown)
	select {
	case <-time.After(countDown):
		//调度信息，存入chan中
		gchScd <- s
	}
	return
}

//作业信息结构
type Job struct {
	id        int64           //作业ID
	name      string          //作业名称
	timeOut   int64           //最大执行时间
	desc      string          //作业说明
	preJobId  int64           //上级作业ID
	preJob    *Job            //上级作业
	nextJobId int64           //下级作业ID
	nextJob   *Job            //下级作业
	tasks     map[int64]*Task //作业中的任务
}

// 任务信息结构
type Task struct {
	id          int64             // 任务的ID
	address     string            // 任务的执行地址
	name        string            // 任务名称
	jobType     string            // 任务类型
	cyc         string            //调度周期
	startSecond int64             //周期内启动时间
	cmd         string            // 任务执行的命令或脚本、函数名等。
	timeOut     int64             // 设定超时时间，0表示不做超时限制。单位秒
	param       map[string]string // 任务的参数信息
	jobId       int64             //所属作业ID
	relTasks    map[int64]*Task   //依赖的任务
}

// 任务依赖结构
type RelTask struct {
	taskId    int64 //任务ID
	reltaskId int64 //依赖任务ID
}

//调度执行信息结构
type ExecSchedule struct {
	batchId   string     //批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)
	schedule  *Schedule  //调度
	startTime time.Time  //开始时间
	endTime   time.Time  //结束时间
	state     string     //状态
	result    int8       //结果
	execType  string     //执行类型
	execJob   []*ExecJob //作业执行信息
}

//作业执行信息结构
type ExecJob struct {
	batchId   string    //批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)
	job       *Job      //作业
	startTime time.Time //开始时间
	endTime   time.Time //结束时间
	state     string    //状态
	result    int8      //结果
	execType  string    //执行类型
	execTask  []*Task   //任务执行信息
}

//任务执行信息结构
type ExecTask struct {
	batchId   string    //批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)
	task      *Task     //任务
	startTime time.Time //开始时间
	endTime   time.Time //结束时间
	state     string    //状态
	result    int8      //结果
	execType  string    //执行类型
}

//从元数据库获取Schedule列表。
func getAllSchedules() (scds []*Schedule, err error) {
	var stime int64

	//查询全部schedule列表
	sql := `SELECT scd.scd_id,
				scd.scd_name,
				scd.scd_num,
				scd.scd_cyc,
				scd.scd_start,
				scd.scd_timeout,
				scd.scd_job_id,
				scd.scd_desc
			FROM hive.scd_schedule scd`

	rows, err := gDbConn.Query(sql)
	checkErr(err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		scd := &Schedule{}
		err = rows.Scan(&scd.id, &scd.name, &scd.count, &scd.cyc, &stime,
			&scd.timeOut, &scd.jobId, &scd.desc)
		scd.startSecond = time.Duration(stime) * time.Second

		//初始化param的内存
		scd.param = make(map[string]string)

		scds = append(scds, scd)
	}

	return scds, err
}

//从元数据库获取Schedule下的Job列表。
func getAllJobs() (jobs map[int64]*Job, err error) {

	jobs = make(map[int64]*Job)

	//查询全部Job列表
	sql := `SELECT job.job_id,
			   job.job_name,
			   job.job_timeout,
			   job.job_desc,
			   job.prev_job_id,
			   job.next_job_id
			FROM hive.scd_job job`

	rows, err := gDbConn.Query(sql)
	checkErr(err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		job := &Job{}
		err = rows.Scan(&job.id, &job.name, &job.timeOut, &job.desc, &job.preJobId, &job.nextJobId)

		//初始化Task内存
		job.tasks = make(map[int64]*Task)
		jobs[job.id] = job
	}

	return jobs, err
}

//从元数据库获取Job下的Task列表。
func getAllTasks() (tasks map[int64]*Task, err error) {

	tasks = make(map[int64]*Task)

	//查询全部Task列表
	sql := `SELECT task.task_id,
			   task.task_address,
			   task.task_name,
			   task.task_type_id,
			   task.task_cyc,
			   task.task_start,
			   task.task_cmd
			FROM hive.scd_task task`

	rows, err := gDbConn.Query(sql)
	checkErr(err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		task := &Task{}
		err = rows.Scan(&task.id, &task.address, &task.name, &task.jobType, &task.cyc, &task.startSecond, &task.cmd)
		//初始化relTask、param的内存
		task.relTasks = make(map[int64]*Task)
		task.param = make(map[string]string)

		tasks[task.id] = task
	}
	return tasks, err
}

//从元数据库获取Job下的Task列表。
func getJobTask() (jobtask map[int64]int64, err error) {

	jobtask = make(map[int64]int64)

	//查询Job中全部Task列表
	sql := `SELECT jt.job_id,
				jt.task_id
			FROM hive.scd_job_task jt`

	rows, err := gDbConn.Query(sql)
	checkErr(err)

	//循环读取记录
	for rows.Next() {
		var jobid, taskid int64
		err = rows.Scan(&jobid, &taskid)
		jobtask[taskid] = jobid
	}
	return jobtask, err
}

//从元数据库获取Task的依赖列表。
func getRelTasks() (relTasks []*RelTask, err error) {

	relTasks = make([]*RelTask, 0)

	//查询Task的依赖列表
	sql := `SELECT tr.task_id,
				tr.rel_task_id
			FROM hive.scd_task_rel tr
			ORDER BY tr.task_id`

	rows, err := gDbConn.Query(sql)
	checkErr(err)

	//循环读取记录
	for rows.Next() {
		var taskid, reltaskid int64
		err = rows.Scan(&taskid, &reltaskid)
		relTasks = append(relTasks, &RelTask{taskId: taskid, reltaskId: reltaskid})
	}
	return relTasks, err
}

//获取距启动的时间（秒）
func getCountDown(cyc string, ss time.Duration) (countDown time.Duration, err error) {
	now := GetNow()
	var startTime time.Time
	//解析周期并取得距下一周期的时间
	switch {
	case cyc == "ss":
		//按秒取整
		s := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(),
			now.Second(), 0, time.Local).Add(time.Second)
		startTime = s.Add(ss)
	case cyc == "mi":
		//按分钟取整
		mi := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0,
			0, time.Local).Add(time.Minute)
		startTime = mi.Add(ss)
	case cyc == "h":
		//按小时取整
		h := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0,
			time.Local).Add(time.Hour)
		startTime = h.Add(ss)
	case cyc == "d":
		//按日取整
		d := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0,
			time.Local).AddDate(0, 0, 1)
		startTime = d.Add(ss)
	case cyc == "m":
		//按月取整
		m := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
		startTime = m.Add(ss)
	case cyc == "w":
		//按周取整
		w := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()), 0, 0, 0, 0, time.Local).AddDate(0, 0, 7)
		startTime = w.Add(ss)
	case cyc == "q":
		//回头再处理
	case cyc == "y":
		//按年取整
		y := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local).AddDate(1, 0, 0)
		startTime = y.Add(ss)
	}

	countDown = startTime.Sub(time.Now())

	return countDown, nil

}

//获取当前时间
func GetNow() time.Time {
	return time.Now().Local()
}
