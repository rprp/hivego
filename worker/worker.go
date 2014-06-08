//worker执行模块worker负责在本地执行调度模块发送的命令，并将输出信息返回给调度模块。
//worker执行时会启动http服务监听8123端口，提供RPC调用接口CmdExecuter.Run()方法。
//package worker
package main

import (
	"bytes"
	"errors"
	"github.com/Sirupsen/logrus"
	"io"
	"net/http"
	"net/rpc"
	"os/exec"
	"runtime"
	"runtime/debug"
	"time"
)

var (
	// 监听端口号
	gPort string = ":8123"
	//全局log对象
	log = logrus.New()
)

func init() { // {{{
	//设置log模块的默认格式
	log.Formatter = new(logrus.TextFormatter) // default
	runtime.GOMAXPROCS(16)
} // }}}

// 任务信息结构
type Task struct {
	Id          int64             // 任务的ID
	Address     string            // 任务的执行地址
	Name        string            // 任务名称
	JobType     string            // 任务类型
	Cyc         string            //调度周期
	StartSecond int64             //周期内启动时间
	Cmd         string            // 任务执行的命令或脚本、函数名等。
	TimeOut     int64             // 设定超时时间，0表示不做超时限制。单位秒
	Param       map[string]string // 任务的参数信息
	Attr        map[string]string // 任务的属性信息
	JobId       int64             //所属作业ID
	RelTasks    map[int64]*Task   //依赖的任务
	RelTaskCnt  int64             //依赖的任务数量
}

//返回的消息
type Reply struct {
	Err    error  //错误信息
	Stdout string //标准输出
}

//RPC结构
//服务端处理部分，接受client端发送的指令。
type CmdExecuter struct {
}

//Run调用相应的模块，完成对Task的执行
//参数task，需要执行的任务信息。
//参数reply，任务执行输出的信息。
func (this *CmdExecuter) Run(task *Task, reply *Reply) error { // {{{
	fn := "worker.CmdExecuter.Run"
	defer func() {
		if err := recover(); err != nil {
			var buf bytes.Buffer
			buf.Write(debug.Stack())
			log.WithFields(logrus.Fields{
				"panic": buf.String(),
			}).Warn(fn)
			reply.Err = errors.New("panic")
			return
		}
	}()

	//channel类型变量chRp用来实时传递命令执行过程中的输出信息。
	chRp := make(chan string, 0)
	go func() {
		for {
			select {
			case msg := <-chRp:
				//设置返回值
				reply.Stdout += msg
				log.WithFields(logrus.Fields{
					"cmdlog": msg,
				}).Info(fn)
			}
		}
	}()

	//执行task任务
	err := runCmd(task, chRp)
	reply.Err = err

	return err
} // }}}

//runCmd用来执行参数cmd中指定的命令，并返回执行时间和错误信息。
func runCmd(task *Task, reply chan string) error { // {{{
	var c *exec.Cmd
	var cmdArgs []string //执行的命令行参数

	//从task结构中获取并组合命令参数
	for _, v := range task.Param {
		cmdArgs = append(cmdArgs, v)
	}

	//记录开始执行的时间
	startTime := time.Now().Format("2006-01-02 15:04:05")

	//命令成功执行标志
	ok := make(chan bool, 1)
	chErr := make(chan error, 1)

	//启动一个goroutine执行任务，超时则直接返回，
	//正常结束则设置成功执行标志ok
	go func() {
		c = exec.Command(task.Cmd, cmdArgs...)

		stdout, err := c.StdoutPipe() //挂载标准输出
		if err != nil {
			log.WithFields(logrus.Fields{
				"StdoutPipe": err.Error(),
			}).Warn("worker.runCmd")
			chErr <- err
			return
		}

		stderr, err := c.StderrPipe() //挂着错误输出
		if err != nil {
			log.WithFields(logrus.Fields{
				"StderrPipe": err.Error(),
			}).Warn("worker.runCmd")
			chErr <- err
			return
		}

		r := io.MultiReader(stdout, stderr)
		if err := c.Start(); err != nil {
			log.WithFields(logrus.Fields{
				"Start": err.Error(),
			}).Warn("worker.runCmd")
			chErr <- err
			return
		}

		//读取输出信息，设置到reply通道中
		for {
			bf := make([]byte, 1024)
			count, err := r.Read(bf)
			if err != nil || count == 0 {
				break
			} else {
				reply <- string(bf)
			}
		}

		if err := c.Wait(); err != nil {
			log.WithFields(logrus.Fields{
				"Wait": err.Error(),
			}).Warn("worker.runCmd")
			chErr <- err
			return
		}

		ok <- true
	}()

	//监听通道，超时则kill掉进程
	select {
	case <-time.After(time.Duration(task.TimeOut) * 1000 * time.Millisecond):
		log.WithFields(logrus.Fields{
			"StartTime": startTime,
			"EndTime":   time.Now().Format("2006-01-02 15:04:05"),
			"TaskId":    task.Id,
			"TaskName":  task.Name,
			"TaskCmd":   task.Cmd,
			"TaskArg":   cmdArgs,
		}).Warn("worker.runCmd is timeout")
		c.Process.Kill()
		return errors.New("time out")
	case e := <-chErr:
		//异常退出
		log.WithFields(logrus.Fields{
			"StartTime": startTime,
			"EndTime":   time.Now().Format("2006-01-02 15:04:05"),
			"TaskId":    task.Id,
			"TaskName":  task.Name,
			"TaskCmd":   task.Cmd,
			"TaskArg":   cmdArgs,
		}).Warn("worker.runCmd is error")
		return e
	case <-ok:
		//正常退出
		log.WithFields(logrus.Fields{
			"StartTime": startTime,
			"EndTime":   time.Now().Format("2006-01-02 15:04:05"),
			"TaskId":    task.Id,
			"TaskName":  task.Name,
			"TaskCmd":   task.Cmd,
			"TaskArg":   cmdArgs,
		}).Info("worker.runCmd is ok")
		return nil
	}

	return nil
} // }}}

//启动HTTP服务监控指定端口
func ListenAndServer(port string) { // {{{
	executer := new(CmdExecuter)
	rpc.Register(executer)
	rpc.HandleHTTP()

	log.WithFields(logrus.Fields{
		"Port": port,
	}).Info("Server is running")

	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("error!")
	}

} // }}}

func main() {
	ListenAndServer(gPort)
}
