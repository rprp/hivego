package main

import (
	"bytes"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/rpc"
	"runtime/debug"
	"sync"
	"time"
)

//调度执行信息结构
type ExecSchedule struct {
	lock           sync.Mutex
	batchId        string              //批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)
	schedule       *Schedule           //调度
	startTime      time.Time           //开始时间
	endTime        time.Time           //结束时间
	state          int8                //状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.意外中止
	result         float32             //结果,调度中执行成功任务的百分比
	execType       int8                //执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行
	execJob        *ExecJob            //作业执行信息
	execTasks      map[int64]*ExecTask //任务执行信息
	jobCnt         int64               //调度中作业数量
	taskCnt        int64               //调度中任务数量
	successTaskCnt int64               //执行成功任务数量
	failTaskCnt    int64               //执行失败任务数量
}

//ExecSchedule.Run()方法执行调度任务。
//过程中会维护一个Chan *ExecTask类型变量staskChan，用来传递执行完成的Task。
//通过遍历Schedule下的全部Task，找出可执行的Task(依赖列表为空的Task)，启动线程执行task.Run
//方法，并将staskChan传给它。当Task执行结束后会把自己放入staskChan中，处理的另一部分监控着
//staskChan，从其中取出执行完毕的task后，会从其它任务的依赖列表中将已执行完毕的task删除，
//并重新找出依赖列表为空的task，启动线程运行它的Run方法。
//全部执行结束后，设置Schedule的下次启动时间。
func (s *ExecSchedule) Run() { // {{{

	//taskChan用来传递完成的任务。
	//当一个作业完成后会将自己放入taskChan变量中
	staskChan := make(chan *ExecTask)

	s.startTime = time.Now()
	s.state = 1
	s.Log()

	l.Infoln("schedule ", s.schedule.name, " is start ", " batchId=", s.batchId)

	//启动独立的任务
	for _, execTask := range s.execTasks {
		//依赖任务列表为空，任务可以执行
		if len(execTask.relExecTasks) == 0 {
			//任务所属作业开始时间为空，设置作业启动信息
			if execTask.execJob.startTime.IsZero() {
				execTask.execJob.startTime = time.Now()
				execTask.execJob.state = 1
				execTask.execJob.Log()

				l.Infoln("job ", execTask.execJob.job.name, " is start ", " batchJobId=", execTask.execJob.batchJobId)

			}

			//执行任务，完成后任务会放入taskChan中
			go execTask.Run(staskChan)
		}
	}

	//不断轮询taskChan中的信息，直到最后一个任务完成
	//调用执行结构的Timer方法，并退出线程。
	for {
		select {
		case t := <-staskChan:
			s.taskCnt--

			//计算任务完成百分比
			s.result = float32(s.schedule.taskCnt-s.taskCnt) / float32(s.schedule.taskCnt)

			if t.state == 3 {
				s.successTaskCnt++

				//设置作业信息
				j := t.execJob
				j.taskCnt--
				//计算任务完成百分比
				j.result = float32(j.job.taskCnt-j.taskCnt) / float32(j.job.taskCnt)
				if j.taskCnt == 0 {
					j.endTime = time.Now()
					j.state = 3
					j.Log()

					l.Infoln("job ", j.job.name, " is end ", " batchJobId=", j.batchJobId, " result=", j.result)
				}

				//任务成功执行，将该任务从其它任务的依赖列表中删除。
				//若删除后依赖列表为空，则启动那个任务。
				for _, nextask := range t.nextExecTasks {
					delete(nextask.relExecTasks, t.task.Id)
					if len(nextask.relExecTasks) == 0 {
						//任务所属作业开始时间为空，设置作业启动信息
						if nextask.execJob.startTime.IsZero() {
							nextask.execJob.startTime = time.Now()
							nextask.execJob.state = 1
							nextask.execJob.Log()
							l.Infoln("job ", nextask.execJob.job.name, " is start ", " batchJobId=", nextask.execJob.batchJobId)
						}
						go nextask.Run(staskChan)
					}

				}
			} else {
				s.failTaskCnt++
				//任务失败，处理下游依赖任务链
				n := clearFailTask(t) - 1
				s.taskCnt -= n

				l.Infoln("task ", t.task.Name, " is fail ", " batchTaskId=", t.batchTaskId, " state=", t.state)

			}

			if s.taskCnt == 0 {
				//全部完成后，写入日志存储至数据库，设置下次启动时间
				s.endTime = time.Now()
				s.state = 3
				s.Log()

				l.Infoln("schedule ", s.schedule.name, " is end ", " batchId=", s.batchId,
					" success=", s.successTaskCnt, " fail=", s.failTaskCnt, " result=", s.result)

				//自动调度执行，完成后设置下次执行时间
				if s.execType == 1 {
					//设置下次执行时间
					go s.schedule.Timer()
				}
				return
			}

		}
	}

} // }}}

//Pause暂停调度执行
func (s *ExecSchedule) Pause() {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, t := range s.execTasks {
		t.state = 2
	}

}

//处理下游依赖任务链
//clearFailTask会处理失败的任务，失败的任务被当做参数传递进来后，会将这个任务从其依赖任务
//的下级列表中删除，若该任务还有下级任务则进行递归调用。
//返回删掉的任务数量。
func clearFailTask(t *ExecTask) (n int64) { // {{{

	if len(t.nextExecTasks) != 0 {
		for _, nextaks := range t.nextExecTasks {
			n += clearFailTask(nextaks)
		}
	}

	for _, reltask := range t.relExecTasks {
		delete(reltask.nextExecTasks, t.task.Id)
	}

	return n + 1
} // }}}

//保存执行日志
func (s *ExecSchedule) Log() (err error) { // {{{

	if s.state == 0 {
		sql := `INSERT INTO hive.scd_schedule_log
						(batch_id,
						 scd_id,
						 start_time,
						 end_time,
						 state,
						 result,
						 batch_type)
			VALUES      (?,
						 ?,
						 ?,
						 ?,
						 ?,
						 ?,
						 ?)`
		_, err = gDbConn.Exec(sql, &s.batchId, &s.schedule.id, &s.startTime, &s.endTime, &s.state, &s.result, &s.execType)
	} else {
		sql := `UPDATE hive.scd_schedule_log
						 set start_time=?,
						 end_time=?,
						 state=?,
						 result=?
				WHERE batch_id=?`
		_, err = gDbConn.Exec(sql, &s.startTime, &s.endTime, &s.state, &s.result, &s.batchId)
	}

	return err
} // }}}

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
	execTasks  []*ExecTask //任务执行信息
	taskCnt    int64       //作业中任务数量
}

//保存执行日志
func (j *ExecJob) Log() (err error) { // {{{
	if j.state == 0 {
		sql := `INSERT INTO hive.scd_job_log
						(batch_job_id,batch_id,
						 job_id,
						 start_time,
						 end_time,
						 state,
						 result,
						 batch_type)
			VALUES      (?,
						 ?,
						 ?,
						 ?,
						 ?,
						 ?,
						 ?,
						 ?)`
		_, err = gDbConn.Exec(sql, &j.batchJobId, &j.batchId, &j.job.id, &j.startTime, &j.endTime, &j.state, &j.result, &j.execType)
	} else {
		sql := `UPDATE hive.scd_job_log
						 set start_time=?,
						 end_time=?,
						 state=?,
						 result=?
				WHERE batch_job_id=?`
		_, err = gDbConn.Exec(sql, &j.startTime, &j.endTime, &j.state, &j.result, &j.batchJobId)
	}

	return err
} // }}}

//任务执行信息结构
type ExecTask struct {
	batchTaskId   string              //任务批次ID，作业批次ID + 任务ID
	batchJobId    string              //作业批次ID，批次ID + 作业ID
	batchId       string              //批次ID，规则scheduleId + 周期开始时间(不含周期内启动时间)
	task          *Task               //任务
	startTime     time.Time           //开始时间
	endTime       time.Time           //结束时间
	state         int8                //状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.失败
	execType      int8                //执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行
	execJob       *ExecJob            //任务所属作业
	output        string              //任务输出
	nextExecTasks map[int64]*ExecTask //下级任务执行信息
	relExecTasks  map[int64]*ExecTask //依赖的任务
}

//保存执行日志
func (t *ExecTask) Log() (err error) { // {{{
	if t.state == 0 {
		sql := `INSERT INTO hive.scd_task_log
						(batch_task_id,batch_job_id,batch_id,
						 task_id,
						 start_time,
						 end_time,
						 state,
						 batch_type)
			VALUES      (?,
						 ?,
						 ?,
						 ?,
						 ?,
						 ?,
						 ?,
						 ?)`
		_, err = gDbConn.Exec(sql, &t.batchTaskId, &t.batchJobId, &t.batchId, &t.task.Id, &t.startTime, &t.endTime, &t.state, &t.execType)
	} else {
		sql := `UPDATE hive.scd_task_log
						 set start_time=?,
						 end_time=?,
						 state=?
				WHERE batch_task_id=?`
		_, err = gDbConn.Exec(sql, &t.startTime, &t.endTime, &t.state, &t.batchTaskId)
	}

	return err
} // }}}

type Reply struct {
	Err    error  //错误信息
	Stdout string //标准输出
}

//Run方法负责执行任务。
//首先会判断是否符合执行条件，符合则执行
//执行时会从任务执行结构中取出需要执行的信息，通过RPC发送给执行模块执行。
//完成后更新执行信息，并将任务置入taskChan变量中，供后续处理。
func (t *ExecTask) Run(taskChan chan *ExecTask) { // {{{
	rl := new(Reply)
	defer func() {
		if err := recover(); err != nil {
			var buf bytes.Buffer
			buf.Write(debug.Stack())
			t.endTime = time.Now()
			t.state = 4
			l.Warningln("task run error", "batchTaskId=", t.batchTaskId, "TaskName=",
				t.task.Name, "output=", rl.Stdout, "err=", err, " stack=", buf.String())
			t.Log()

			taskChan <- t
			return
		}
	}()

	//若任务为暂停状态则不执行直接退出
	if t.state == 2 {
		l.Infoln("task ", t.task.Name, " is ignore batchTaskId=", t.batchTaskId)
		t.Log()
		taskChan <- t
		return
	}

	//判断是否在执行周期内,若是则直接执行，否则跳过返回执行完成的状态，并继续下一步骤
	//TO-DO 暂时搁着，以后再完善

	t.startTime = time.Now()
	t.state = 1

	t.Log()
	l.Infoln("task ", t.task.Name, " is start batchTaskId=", t.batchTaskId)

	//执行任务
	address := t.task.Address
	task := t.task

	t.state = 3

	if client, err := rpc.DialHTTP("tcp", address+gPort); err == nil {

		if err := client.Call("CmdExecuter.Run", task, &rl); err == nil {
			if rl.Err != nil {
				t.output = rl.Err.Error()
				t.state = 4
			}
		} else {
			panic(err.Error())
		}
	} else {
		panic("unexpected HTTP ")
	}

	t.output = t.output + rl.Stdout
	t.endTime = time.Now()

	t.Log()
	l.Infoln("task", t.task.Name, "is end batchTaskId =", t.batchTaskId, "state =",
		t.state, "output =", rl.Stdout)

	taskChan <- t

} // }}}

func NewExecSchedule(rscd *Schedule) (execScd *ExecSchedule, err error) { // {{{
	//批次ID
	bid := fmt.Sprintf("%s %d", time.Now().Format("2006-01-02 15:04:05.000000"), rscd.id)
	return NewExecScheduleById(bid, rscd)

} // }}}

//NewExecSchedule会根据传入的Schedule参数来构建一个调度的执行结构。
//执行结构包含完整的执行链，构造完成后会调用ExecSchedule的Run方法来开始执行。
func NewExecScheduleById(bid string, rscd *Schedule) (execScd *ExecSchedule, err error) { // {{{

	execScd = &ExecSchedule{
		batchId:   bid,  //批次ID
		schedule:  rscd, //调度
		state:     0,    //状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.意外中止
		result:    0,    //结果,调度中执行成功任务的百分比
		execType:  1,    //执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行
		jobCnt:    rscd.jobCnt,
		taskCnt:   rscd.taskCnt,
		execTasks: make(map[int64]*ExecTask), //设置任务列表
	}

	err = execScd.Log()

	//构建调度中的作业执行结构
	execScd.execJob, err = NewExecJob(bid, rscd.job)

	//生成调度中的任务列表
	for j := execScd.execJob; j != nil; {
		for _, t := range j.execTasks {
			execScd.execTasks[t.task.Id] = t
		}
		j = j.nextJob
	}

	l.Infoln("ExecSchedule ", execScd.schedule.name, " is create batchId=", bid)

	return execScd, err
} // }}}

//NewExecJob根据输入的job和batchId构建作业执行链，并返回。
func NewExecJob(batchId string, job *Job) (execJob *ExecJob, err error) { // {{{
	bjd := fmt.Sprintf("%s.%d", batchId, job.id)

	execJob = &ExecJob{
		batchJobId: bjd,     //作业批次ID，批次ID+作业ID
		batchId:    batchId, //批次ID
		job:        job,     //作业
		state:      0,       //状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.意外中止
		result:     0,       //结果,作业中执行成功任务的百分比
		execType:   1,       //执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行
		execTasks:  make([]*ExecTask, 0),
	}

	err = execJob.Log()

	//构建当前作业中的任务执行结构
	for _, t := range execJob.job.tasks {
		execTask := new(ExecTask)
		bjtd := fmt.Sprintf("%s.%d", bjd, t.Id)
		execTask.batchTaskId = bjtd
		execTask.batchJobId = bjd
		execTask.batchId = batchId
		execTask.task = t
		execTask.state = 0         //状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.意外中止
		execTask.execType = 1      //执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行
		execTask.execJob = execJob //任务所属的job

		//将任务执行结构存入全局的gExecTasks变量，以便后面获取依赖任务执行信息时使用
		gExecTasks[t.Id] = execTask

		//依赖任务分配内存
		execTask.relExecTasks = make(map[int64]*ExecTask)
		execTask.nextExecTasks = make(map[int64]*ExecTask)

		//设置依赖任务执行信息
		//先获取依赖任务列表，通过每一个依赖任务的id从全局gExecTasks中获取到依赖任务
		for _, relTask := range t.RelTasks {
			retask := gExecTasks[relTask.Id]
			execTask.relExecTasks[relTask.Id] = retask

			//将execTask设置为依赖任务的下级任务
			retask.nextExecTasks[execTask.task.Id] = execTask

		}

		err = execTask.Log()

		execJob.execTasks = append(execJob.execTasks, execTask)

		l.Infoln("ExecTask ", execTask.task.Name, " is create batchTaskId=", execTask.batchTaskId)
	}

	l.Infoln("ExecJob ", execJob.job.name, " is create batchJobId=", execJob.batchJobId)

	//继续构建作业的下级作业
	if job.nextJob != nil {
		//递归调用，直到最后一个作业
		execJob.nextJob, err = NewExecJob(batchId, job.nextJob)
	}
	return execJob, err

} // }}}

//ExecSchedule.Restore(batchId string)方法修复执行指定的调度。
//根据传入的batchId，构建调度执行结构，并调用Run方法执行其中的任务
func Restore(batchId string, scdId int64) (err error) { // {{{

	l.Infoln("Restore schedule by ", " batchid=", batchId, " scdId=", scdId)

	//获取执行成功的Task
	successTaskId := getSuccessTaskId(batchId)

	//创建ExecSchedule结构
	execSchedule, err := NewExecScheduleById(batchId, gScdList.schedules[scdId])
	execSchedule.execType = 3
	execSchedule.state = 1

	//删除成功的任务
	for _, tId := range successTaskId {
		t := execSchedule.execTasks[tId]
		for _, nextaks := range t.nextExecTasks {
			delete(nextaks.relExecTasks, t.task.Id)
		}
		delete(execSchedule.execTasks, tId)
		execSchedule.taskCnt--
	}

	//设置作业、任务的初始状态
	for _, t := range execSchedule.execTasks {
		t.execType = 3
		t.state = 1
		t.execJob.execType = 3
		t.execJob.state = 1
	}

	l.Infoln("schedule will restore")

	//执行
	execSchedule.Run()

	l.Infoln("schedule was restored")

	return nil

} // }}}

//getSuccessTaskId会根据传入的batchId从元数据库查找出执行成功的task
func getSuccessTaskId(batchId string) []int64 { // {{{

	sql := `SELECT task_id
			FROM   hive.scd_task_log
			WHERE  state = 3
			   AND batch_id =?`

	rows, err := gDbConn.Query(sql, batchId)
	checkErr(err)

	taskIds := make([]int64, 0)
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		taskIds = append(taskIds, id)
	}

	return taskIds
} // }}}
