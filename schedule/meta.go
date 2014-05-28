//调度模块的数据结构
//package schedule
package main

import (
	_ "database/sql"
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
	startSecond int64             //周期内启动时间
	timeOut     int64             //最大执行时间
	jobId       int64             //作业ID
	job         *Job              //作业
	desc        string            //调度说明
}

////根据调度的周期及启动时间，按时将调度传至执行列表执行。
//func (s *Schedule) Timer() (execSchedule ExecSchedule, err error) {

//return
////构建执行结构，存入chan中
//}

//作业信息结构
type Job struct {
	id        int64            //作业ID
	name      string           //作业名称
	timeOut   int64            //最大执行时间
	desc      string           //作业说明
	preJobId  int64            //上级作业ID
	preJob    *Job             //上级作业
	nextJobId int64            //下级作业ID
	nextJob   *Job             //下级作业
	task      map[string]*Task //作业中的任务
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
	relTask     map[string]*Task  //依赖的任务
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
		err = rows.Scan(&scd.id, &scd.name, &scd.count, &scd.cyc, &scd.startSecond,
			&scd.timeOut, &scd.jobId, &scd.desc)
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
		tasks[task.id] = task
	}
	return tasks, err
}
