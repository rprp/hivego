//调度模块，负责从元数据库读取并解析调度信息。
//将需要执行的任务发送给执行模块，并读取返回信息。
//package schedule
package main

import (
	"database/sql"
	"fmt"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"time"
)

//全局变量定义
var (
	//全局log对象
	gLog = logrus.New()
	p    = gLog.WithFields

	gPort   string  // 监听端口号
	gDbConn *sql.DB //数据库链接

	gScdList *ScheduleList //全局调度列表

	gScdChan     chan *Schedule    //执行的调度结构
	gExecScdChan chan ExecSchedule //执行的调度结构

	gExecTasks map[int64]*ExecTask
)

//初始化工作
func init() { // {{{
	runtime.GOMAXPROCS(16)

	//设置log模块的默认格式
	gLog.Formatter = new(logrus.TextFormatter) // default
	gLog.Level = logrus.Debug

	//从配置文件中获取数据库连接、服务端口号等信息
	gPort = ":8123"

	gScdChan = make(chan *Schedule)
	gExecTasks = make(map[int64]*ExecTask)

} // }}}

//StartSchedule函数是调度模块的入口函数。程序初始化完成后，它负责连接元数据库，
//获取调度信息，在内存中构建Schedule结构。完成后，会调用Schedule的Timer方法。
//Timer方法会根据调度周期及启动时间，按时启动，随后会依据Schedule信息构建执行结构
//并送入chan中。
//模块的另一部分在不断的检测chan中的内容，将取到的执行结构体后创建新的goroutine
//执行。
func StartSchedule() error { // {{{
	// 连接数据库
	cnn, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/hive?charset=utf8")
	checkErr(err)
	gDbConn = cnn

	defer gDbConn.Close()

	//调度列表
	sLst := &ScheduleList{}

	sLst.InitSchedules()

	sLst.Run()

	//从chan中得到需要执行的调度，启动一个线程执行
	for {
		select {
		case rscd := <-gScdChan:
			fmt.Println(time.Now(), "\t", rscd.name, "is start")
			//启动一个线程开始构建执行结构链
			s, err := NewExecSchedule(rscd)
			checkErr(err)
			//启动线程执行调度任务
			go s.Run()

		}

	}

	return nil
} // }}}

func main() {
	StartSchedule()
}
