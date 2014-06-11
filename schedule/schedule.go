//调度模块，负责从元数据库读取并解析调度信息。
//将需要执行的任务发送给执行模块，并读取返回信息。
//package schedule
package main

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
)

//全局变量定义
var (
	//全局log对象
	l = logrus.New()
	p = l.WithFields

	gPort   string  // 监听端口号
	gDbConn *sql.DB //数据库链接

	gScdList *ScheduleList //全局调度列表

	gExecScdChan chan ExecSchedule //执行的调度结构

	gTasks map[int64]*Task

	gExecTasks map[int64]*ExecTask
)

//初始化工作
func init() { // {{{
	runtime.GOMAXPROCS(16)

	//设置log模块的默认格式
	l.Formatter = new(logrus.TextFormatter) // default
	l.Level = logrus.Info

	//从配置文件中获取数据库连接、服务端口号等信息
	gPort = ":8123"

	gExecTasks = make(map[int64]*ExecTask)
	gTasks = make(map[int64]*Task)

} // }}}

//调度列表
type ScheduleList struct {
	schedules map[int64]*Schedule //调度列表
	tasks     map[int64]*Task     //任务列表
	jobs      map[int64]*Job      //作业列表
}

//从元数据库获取Job列表
func (sl *ScheduleList) setJobs() (err error) { // {{{
	sl.jobs, err = getAllJobs()
	return err
} // }}}

//从元数据库获取Task列表
func (sl *ScheduleList) setTasks() (err error) { // {{{
	sl.tasks, err = getAllTasks()
	return err
} // }}}

//从元数据库获取Schedule列表
func (sl *ScheduleList) setSchedules() (err error) { // {{{
	sl.schedules, err = getAllSchedules()
	return err
} // }}}

//执行调度会调用Schedule的Timer方法。
//Timer方法会根据调度周期及启动时间，按时启动，随后会依据Schedule信息构建执行结构
//并送入chan中。
func (sl *ScheduleList) StartSchedule() { // {{{

	for _, scd := range sl.schedules {
		go scd.Timer()
	}

} // }}}

//InitSchedules方法，初始化调度列表
//获取调度信息，在内存中构建Schedule结构。
func (sl *ScheduleList) InitSchedules() (err error) { // {{{

	//从元数据库读取调度信息
	sl.setSchedules()
	sl.setJobs()
	sl.setTasks()

	reltasks, err := getRelTasks() //获取Task的依赖链

	jobtask, err := getJobTaskid() //获取Job的Task列表

	//设置job中的task列表
	//由于框架规定一个task只能在一个job中，N:1关系
	//只需遍历一遍task与job对应关系结构，从jobs的map中找出job设置它的task即可
	for taskid, jobid := range jobtask {
		sl.jobs[jobid].tasks[taskid] = sl.tasks[taskid]
		sl.jobs[jobid].taskCnt++
	}
	l.Infoln("set task in job")

	//设置task的依赖链
	for _, maptask := range reltasks {
		sl.tasks[maptask.taskId].RelTasks[maptask.reltaskId] = sl.tasks[maptask.reltaskId]
		sl.tasks[maptask.taskId].RelTaskCnt++
	}
	l.Infoln("set task relation")

	//构建调度链信息
	for _, scd := range sl.schedules {
		var ok bool

		if scd.job, ok = sl.jobs[scd.jobId]; !ok {
			continue
		}
		//设置调度中的job
		scd.jobCnt++
		scd.taskCnt = scd.job.taskCnt

		//设置job链
		for j := scd.job; j.nextJobId != 0; {
			j.nextJob = sl.jobs[j.nextJobId]
			j.preJob = sl.jobs[j.preJobId]
			j = j.nextJob
			scd.jobCnt++
			scd.taskCnt += j.taskCnt
			l.Infoln(scd.name, "-", j.name, " was created")

		}

		l.Infoln(scd.name, " was created", " jobcnt=", scd.jobCnt, " taskcnt=", scd.taskCnt)
	}

	return nil
} // }}}

//StartSchedule函数是调度模块的入口函数。
func StartSchedule() error { // {{{
	// 连接数据库
	cnn, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/hive?charset=utf8")
	checkErr(err)
	gDbConn = cnn

	defer gDbConn.Close()

	//创建并初始化调度列表
	sLst := &ScheduleList{}
	sLst.InitSchedules()

	//printSchedule(sLst.schedules)

	//执行调度
	sLst.StartSchedule()

	s := make(chan int64)
	<-s

	return nil
} // }}}

func main() {
	StartSchedule()
}
