//worker执行模块worker负责在本地执行调度模块发送的命令，并将输出信息返回给调度模块。
//worker执行时会启动http服务监听8123端口，提供RPC调用接口CmdExecuter.Run()方法。
package worker

import (
	"bytes"
	"errors"
	"github.com/Sirupsen/logrus"
	"io"
	"net"
	"net/rpc"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var (

	//全局log对象
	l = logrus.New()
	p = l.WithFields
)

func init() { // {{{
	//设置log模块的默认格式
	l.Formatter = new(logrus.TextFormatter) // default
	l.Level = logrus.Info
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
	Param       []string          // 任务的参数信息
	Attr        map[string]string // 任务的属性信息
	JobId       int64             //所属作业ID
	RelTasks    map[string]*Task  //依赖的任务
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

	//执行task任务
	err := runCmd(task, reply)

	return err
} // }}}

//runCmd用来执行参数cmd中指定的命令，并返回执行时间和错误信息。
func runCmd(task *Task, reply *Reply) error { // {{{
	defer func() {
		if err := recover(); err != nil {
			var buf bytes.Buffer
			buf.Write(debug.Stack())
			l.Warnln("panic=", buf.String())
			reply.Err = errors.New("panic")
			return
		}
	}()
	var c *exec.Cmd
	var cmdArgs []string //执行的命令行参数

	//从task结构中获取并组合命令参数
	for _, v := range task.Param {
		cmdArgs = append(cmdArgs, v)
	}

	//命令成功执行标志
	ok := make(chan bool, 1)
	chErr := make(chan error, 1)

	cmd := strings.TrimSpace(task.Cmd)
	c = exec.Command(cmd, cmdArgs...)
	//启动一个goroutine执行任务，超时则直接返回，
	//正常结束则设置成功执行标志ok
	go func() {

		stdout, err := c.StdoutPipe() //挂载标准输出
		if err != nil {
			l.Warnln("StdoutPipe=", err.Error())
			chErr <- err
			return
		}

		stderr, err := c.StderrPipe() //挂着错误输出
		if err != nil {
			l.Warnln("StderrPipe=", err.Error())
			chErr <- err
			return
		}

		r := io.MultiReader(stdout, stderr)
		if err := c.Start(); err != nil {
			l.Warnln("Start=", err.Error())
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
				l.Infoln("cmdStdout=", string(bf))
				reply.Stdout += string(bf)
			}
		}

		if err := c.Wait(); err != nil {
			l.Warnln("Wait=", err.Error())
			chErr <- err
			return
		}

		ok <- true
	}()

	//监听通道，超时则kill掉进程
	if task.TimeOut > 0 {
		select {
		case <-time.After(time.Duration(task.TimeOut) * 1000 * time.Millisecond):
			c.Process.Kill()
			l.Warnln("runCmd is time out TaskName=", task.Name, "TaskCmd=", task.Cmd, "TaskArg=",
				cmdArgs, "Error=", "time out")
			reply.Err = errors.New("time out")
			return errors.New("time out")
		case e := <-chErr:
			//异常退出
			l.Warnln("runCmd is err TaskName=", task.Name, "TaskCmd=", task.Cmd, "TaskArg=",
				cmdArgs, "Error=", e.Error())
			reply.Err = e
			return e
		case <-ok:
			//正常退出
			l.Infoln("runCmd is ok TaskName=", task.Name, "TaskCmd=", task.Cmd, "TaskArg=",
				cmdArgs)
			return nil
		}
	} else {

		select {
		case e := <-chErr:
			//异常退出
			l.Warnln("runCmd is err TaskName=", task.Name, "TaskCmd=", task.Cmd, "TaskArg=",
				cmdArgs, "Error=", e.Error())
			reply.Err = e
			return e
		case <-ok:
			//正常退出
			l.Infoln("runCmd is ok TaskName=", task.Name, "TaskCmd=", task.Cmd, "TaskArg=",
				cmdArgs)
			return nil
		}
	}

	return nil
} // }}}

//启动HTTP服务监控指定端口
func ListenAndServer(port string) { // {{{
	executer := new(CmdExecuter)
	rpc.Register(executer)

	l.Infoln("Worker is running Port:", port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	checkErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go func() {
				rpc.ServeConn(conn)
			}()
		}
	}()

} // }}}

func checkErr(err error) { // {{{
	if err != nil {
		l.Infoln("error", err.Error())
		panic(err)
	}
} // }}}
