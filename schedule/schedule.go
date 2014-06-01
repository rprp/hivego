//调度模块，负责从元数据库读取并解析调度信息。
//将需要执行的任务发送给执行模块，并读取返回信息。
//package schedule
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"time"
)

//全局变量定义
var (
	gPort   string  // 监听端口号
	gDbConn *sql.DB //数据库链接

	gScds []Schedule //全局调度列表

	gchScd   chan *Schedule    //执行的调度结构
	gchExScd chan ExecSchedule //执行的调度结构
)

//初始化工作
func init() {
	runtime.GOMAXPROCS(16)

	//从配置文件中获取数据库连接、服务端口号等信息
	gPort = ":8123"

	gchScd = make(chan *Schedule)

}

//StartSchedule函数是调度模块的入口函数。程序初始化完成后，它负责连接元数据库，
//获取调度信息，在内存中构建Schedule结构。完成后，会调用Schedule的Timer方法。
//Timer方法会根据调度周期及启动时间，按时启动，随后会依据Schedule信息构建执行结构
//并送入chan中。
//模块的另一部分在不断的检测chan中的内容，将取到的执行结构体后创建新的goroutine
//执行。
func StartSchedule() error {
	// 连接数据库
	cnn, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/hive?charset=utf8")
	checkErr(err)
	gDbConn = cnn

	defer gDbConn.Close()

	//连接元数据库，初始化调度信息至内存
	reltasks, err := getRelTasks() //获取Task的依赖链
	checkErr(err)

	jobtask, err := getJobTask() //获取Job的Task列表
	checkErr(err)

	tasks, err := getAllTasks() //获取Task列表
	checkErr(err)

	jobs, err := getAllJobs() //获取Job列表
	checkErr(err)

	schedules, err := getAllSchedules() //获取Schedule列表
	checkErr(err)

	//设置job中的task列表
	//由于框架规定一个task只能在一个job中，N:1关系
	//只需遍历一遍task与job对应关系结构，从jobs的map中找出job设置它的task即可
	for taskid, jobid := range jobtask {
		jobs[jobid].tasks[taskid] = tasks[taskid]
		jobs[jobid].taskCnt++
	}

	//设置task的依赖链
	for _, maptask := range reltasks {
		tasks[maptask.taskId].relTasks[maptask.reltaskId] = tasks[maptask.reltaskId]
		tasks[maptask.taskId].relTaskCnt++
	}

	//构建调度链信息
	for _, scd := range schedules {
		var ok bool

		if scd.job, ok = jobs[scd.jobId]; !ok {
			continue
		}
		//设置调度中的job
		scd.jobCnt++
		scd.taskCnt = scd.job.taskCnt

		//设置job链
		for j := scd.job; j.nextJobId != 0; {
			j.nextJob = jobs[j.nextJobId]
			j.preJob = jobs[j.preJobId]
			j = j.nextJob
			scd.jobCnt++
			scd.taskCnt += j.taskCnt

		}

		//当构建完成一个调度后，调用它的Timer方法。
		go scd.Timer()

	}

	//打印调度信息
	printSchedule(schedules)

	//从chan中得到需要执行的调度，启动一个线程执行
	for {
		select {
		case rscd := <-gchScd:
			fmt.Println(time.Now(), "\t", rscd.name, "is start")
			//启动一个线程开始构建执行结构链
			exScd, err := constructSchedule(rscd)

			fmt.Println(exScd, err)
			//构建完成后，启动线程执行调度任务
			//执行体exScd中包含了Schedule信息，
			//当全部执行结束后，会设置Schedule的下次启动时间。
			go exScd.Run()

		}

	}

	return nil
}

func main() {
	StartSchedule()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

//启动一个线程开始构建执行结构链
func constructSchedule(rscd *Schedule) (exScd *ExecSchedule, err error) {
	exScd = new(ExecSchedule)

	bid := fmt.Sprintf("%d-%s", rscd.id, time.Now().Format("2006-01-02 15:04:05.000000"))
	exScd.batchId = bid                                        //批次ID
	exScd.schedule = rscd                                      //调度
	exScd.startTime = time.Now()                               //开始时间
	exScd.state = "0"                                          //状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.意外中止
	exScd.result = 0                                           //结果,调度中执行成功任务的百分比
	exScd.execType = "1"                                       //执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行
	exScd.execJob, err = constructJob(exScd.batchId, rscd.job) //作业执行信息
	exScd.jobCnt = rscd.jobCnt
	exScd.taskCnt = rscd.taskCnt

	fmt.Println(exScd.batchId)
	return exScd, err
}

//构建作业执行结构
func constructJob(batchId string, job *Job) (exJob *ExecJob, err error) {
	exJob = new(ExecJob)

	bjd := fmt.Sprint("%s-%d", batchId, job.id)
	exJob.batchJobId = bjd  //作业批次ID，批次ID+作业ID
	exJob.batchId = batchId //批次ID
	exJob.job = job         //作业
	exJob.state = "0"       //状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.意外中止
	exJob.result = 0        //结果,作业中执行成功任务的百分比
	exJob.execType = "1"    //执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行

	//作业中的任务执行结构
	tasks := make([]*ExecTask, 0)
	for _, t := range exJob.job.tasks {
		exTask := new(ExecTask)
		bjtd := fmt.Sprint("%s-%d", bjd, t.id)
		exTask.batchTaskId = bjtd
		exTask.batchJobId = bjd
		exTask.batchId = batchId
		exTask.task = t
		exTask.state = "0"    //状态 0.不满足条件未执行 1. 执行中 2. 暂停 3. 完成 4.意外中止
		exTask.result = 0     //结果，0.成功 1.失败
		exTask.execType = "1" //执行类型 1. 自动定时调度 2.手动人工调度 3.修复执行
		tasks = append(tasks, exTask)

	}
	exJob.execTask = tasks

	if job.nextJob != nil {
		//递归调用，直到最后一个作业
		exJob.nextJob, err = constructJob(batchId, job.nextJob)
	}
	return exJob, err

}

//打印调度信息
func printSchedule(scds []*Schedule) {
	for _, scd := range scds {
		fmt.Println(scd.name, "\tjobs=", scd.jobCnt, " tasks=", scd.taskCnt)
		//打印调度中的作业信息
		for j := scd.job; j != nil; {
			fmt.Println("\t--------------------------------------")
			fmt.Println("\t", j.name)
			//打印作业中的任务信息
			for _, t := range j.tasks {
				fmt.Println("\t\t", t.name)

				fmt.Print("\t\t\t[")
				//打印任务依赖链
				for _, r := range t.relTasks {
					fmt.Print(r.name, ",")

				}
				fmt.Print("]\n")
			}
			fmt.Print("\n")
			j = j.nextJob

		}

	}
}
