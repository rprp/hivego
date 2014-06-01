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
	startSecond time.Duration     //周期内启动时间
	nextStart   time.Time         //周期内启动时间
	timeOut     int64             //最大执行时间
	jobId       int64             //作业ID
	job         *Job              //作业
	desc        string            //调度说明
	jobCnt      int32             //调度中作业数量
	taskCnt     int32             //调度中任务数量
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
	taskCnt   int32           //调度中任务数量
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
	relTaskCnt  int32             //依赖的任务数量
}

// 任务依赖结构
type RelTask struct {
	taskId    int64 //任务ID
	reltaskId int64 //依赖任务ID
}

//调度执行信息结构
type ExecSchedule struct {
	batchId        string    //批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)
	schedule       *Schedule //调度
	startTime      time.Time //开始时间
	endTime        time.Time //结束时间
	state          int8      //状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.意外中止
	result         float32   //结果,调度中执行成功任务的百分比
	execType       int8      //执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行
	execJob        *ExecJob  //作业执行信息
	jobCnt         int32     //调度中作业数量
	taskCnt        int32     //调度中任务数量
	successTaskCnt int32     //执行成功任务数量
	failTaskCnt    int32     //执行失败任务数量
}

//启动线程执行调度任务
//全部执行结束后，设置Schedule的下次启动时间。
func (s *ExecSchedule) Run() {
	//taskChan用来传递完成任务的状态。
	//当一个作业完成后会成功将true放入taskChan变量中，失败放入false
	taskChan := make(chan bool)
	s.state = 1

	//执行调度中作业的Run方法
	//该方法会启动线程执行Job中的Task，并会递归调用。直到最后一个Job
	go s.execJob.Run(taskChan)

	//不断轮询ok中的信息，直到最后一个任务完成
	//调用执行结构的Timer方法，并退出线程。
	for {
		select {
		case tc := <-taskChan:
			s.taskCnt--

			//计算任务完成百分比
			s.result = float32(s.schedule.taskCnt-s.taskCnt) / float32(s.schedule.taskCnt)

			if tc {
				s.successTaskCnt++
			} else {
				s.failTaskCnt++
			}

			if s.taskCnt == 0 {
				//全部完成后，写入日志存储至数据库，设置下次启动时间
				s.endTime = time.Now()
				s.state = 3
				err := s.Save()
				go s.schedule.Timer()
				return
			}

		}
	}

}

//保存执行日志
func (s *ExecSchedule) Save() error {
	return nil
}

//作业执行信息结构
type ExecJob struct {
	batchJobId string      //作业批次ID，批次ID + 作业ID
	batchId    string      //批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)
	job        *Job        //作业
	startTime  time.Time   //开始时间
	endTime    time.Time   //结束时间
	state      int8        //状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.意外中止
	result     float32     //结果执行成功任务的百分比
	nextJob    *ExecJob    //下一个作业
	execType   int8        //执行类型1. 自动定时调度 2.手动人工调度 3.修复执行
	execTask   []*ExecTask //任务执行信息
	taskCnt    int32       //作业中任务数量
}

//启动线程执行作业，并将状态标志传递给执行的任务线程
func (j *ExecJob) Run(taskChan chan bool) {
	jobChan := make(chan bool)
	j.startTime = time.Now()
	j.state = 1

	//对于作业中的每个任务都启动一个线程去执行
	for _, execTask := range j.execTask {
		go execTask.Run(taskChan, jobChan)

	}

	if j.nextJob != nil {
		go j.nextJob.Run(taskChan)
	}

	for {
		select {
		case jc := <-jobChan:
			j.taskCnt--

			//计算任务完成百分比
			j.result = float32(j.job.taskCnt-j.taskCnt) / float32(j.job.taskCnt)
			if j.taskCnt == 0 {
				j.endTime = time.Now()
				j.state = 3
				return
			}

		}

	}

}

//任务执行信息结构
type ExecTask struct {
	batchTaskId string      //任务批次ID，作业批次ID + 任务ID
	batchJobId  string      //作业批次ID，批次ID + 作业ID
	batchId     string      //批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)
	task        *Task       //任务
	startTime   time.Time   //开始时间
	endTime     time.Time   //结束时间
	state       int8        //状态
	result      int8        //结果
	execType    string      //执行类型
	relExecTask []*ExecTask //依赖的任务
}

//任务执行
func (t *ExecTask) Run(taskChan chan bool, jobChan chan bool) {
	//判断是否在执行周期内

	//判断是否依赖任务都执行完毕

	t.startTime = time.Now()

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
