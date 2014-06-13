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

	gPort    string  // 监听端口号
	dbString string  //数据库连接串
	gDbConn  *sql.DB //数据库链接

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

	dbString = "root:@tcp(127.0.0.1:3306)/hive?charset=utf8"
} // }}}

//调度列表
type ScheduleList struct {
	Schedules map[int64]*Schedule //调度列表
}

//从元数据库获取Schedule列表
//执行调度会调用Schedule的Timer方法。
//Timer方法会根据调度周期及启动时间，按时启动，随后会依据Schedule信息构建执行结构
//并送入chan中。
func (sl *ScheduleList) StartSchedule() { // {{{

	for _, scd := range sl.Schedules {
		go scd.Timer()
	}

} // }}}

//InitSchedules方法，初始化调度列表
//获取调度信息，在内存中构建Schedule结构。
func (sl *ScheduleList) InitSchedules() (err error) { // {{{

	//从元数据库读取调度信息
	sl.Schedules, err = getAllSchedules()

	//构建调度链信息
	for _, scd := range sl.Schedules {
		scd.refreshSchedule()
		l.Infoln(scd.name, " was created", " jobcnt=", scd.jobCnt, " taskcnt=", scd.taskCnt)
	}

	return nil
} // }}}

//StartSchedule函数是调度模块的入口函数。
func StartSchedule() error { // {{{
	// 连接数据库
	cnn, err := sql.Open("mysql", dbString)
	checkErr(err)
	gDbConn = cnn

	defer gDbConn.Close()

	//创建并初始化调度列表
	sLst := &ScheduleList{}
	sLst.InitSchedules()

	//执行调度
	sLst.StartSchedule()

	s := make(chan int64)
	<-s

	return nil
} // }}}

func main() {
	StartSchedule()
}
