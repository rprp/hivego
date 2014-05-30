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
	}

	//设置task的依赖链
	for _, maptask := range reltasks {
		tasks[maptask.taskId].relTasks[maptask.reltaskId] = tasks[maptask.reltaskId]
	}

	//构建调度链信息
	for _, scd := range schedules {
		//设置调度中的job
		scd.job = jobs[scd.jobId]

		//设置job链
		for j := scd.job; j.nextJobId != 0; {
			j.nextJob = jobs[j.nextJobId]
			j.preJob = jobs[j.preJobId]
			j = j.nextJob

		}

		//当构建完成一个调度后，调用它的Timer方法。
		go scd.Timer()

	}

	//从chan中得到需要执行的调度，启动一个线程执行
	for {
		select {
		case rscd := <-gchScd:
			fmt.Println(time.Now(), "\t", rscd.name, "is start")
			//启动一个线程开始构建执行结构链
			exScd, err := constructSchedule(rscd)

			//构建完成后，启动线程执行调度任务
			//执行体exScd中包含了Schedule信息，
			//当全部执行结束后，会设置Schedule的下次启动时间。
			go exScd.Run()

			//完成后，设置下次启动时间
			//go rscd.Timer()
		}

	}

	//打印调度信息
	printSchedule(schedules)

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
func constructSchedule(rscd *Schedule) (exScd ExecSchedule, err error) {

	return exScd, nil
}

//打印调度信息
func printSchedule(scds []*Schedule) {
	for _, scd := range scds {
		fmt.Println(scd.name)
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
