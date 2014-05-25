//package worker
package main

import (
	"bytes"
	_ "fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"net/http"
	"net/rpc"
	"os/exec"
	_ "runtime"
	"runtime/debug"
	"time"
)

var (
	// 监听端口号
	gPort string = ":8123"
	//全局log对象
	log = logrus.New()
)

func init() {
	//设置log模块的默认格式
	log.Formatter = new(logrus.TextFormatter) // default
}

// 任务信息结构
type Job struct {
	Id      int64             // job的id
	Name    string            // 名称
	Type    string            // 类型
	Cmd     string            // job执行的命令或脚本、函数名等。
	TimeOut int64             // 设定超时时间，0表示不做超时限制。单位秒
	Param   map[string]string // Job的参数信息
}

//RPC返回的消息
type Reply struct {
	Stdout string //命令的标准输出
}

// 任务接口
//type Executer interface {
// 根据传入的任务信息执行指定任务
// 发生错误或超时则返回error
// 返回的数据集，在interface类型的变量中
//Execute(job Job) (interface{}, error)
//}

//RPC结构
//服务端处理部分，接受client端发送的指令。
type Executer struct {
}

//Run调用相应的模块，完成对Job的执行
//参数job，需要执行的任务信息。
//返回任务的输出
func (this *Executer) Run(job *Job, reply *Reply) error {
	defer func() {
		if err := recover(); err != nil {
			var buf bytes.Buffer
			buf.Write(debug.Stack())
			log.WithFields(logrus.Fields{
				"panic": buf.String(),
			}).Warn("worker.Executer.Run()")
		}
	}()

	re := make(chan string, 0)

	go func() {
		for {
			select {
			case msg := <-re:
				//设置返回值
				reply.Stdout += msg
				log.WithFields(logrus.Fields{
					"cmdlog": msg,
				}).Warn("worker.Executer.Run()")
			}
		}
	}()

	//执行job任务
	err := runCmd(job, re)

	return err
}

//runCmd用来执行参数cmd中指定的命令，并返回执行时间和错误信息。
func runCmd(job *Job, reply chan string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			var buf bytes.Buffer
			buf.Write(debug.Stack())
			log.WithFields(logrus.Fields{
				"panic": buf.String(),
			}).Warn("worker.runCmd(job *Job)")
		}
	}()

	var (
		c *exec.Cmd

		//执行的命令行参数
		cmdArgs []string
	)

	//记录开始执行的时间
	startTime := time.Now().Format("2006-01-02 15:04:05")

	//命令成功执行标志
	ok := make(chan bool, 1)
	chErr := make(chan bool, 1)

	//从job结构中获取并组合命令参数
	for _, v := range job.Param {
		cmdArgs = append(cmdArgs, v)
	}

	//启动一个goroutine执行任务，超时则直接返回，
	//正常结束则设置成功执行标志ok
	go func() {
		c = exec.Command(job.Cmd, cmdArgs...)
		stdout, err := c.StdoutPipe()
		if err != nil {
			log.WithFields(logrus.Fields{
				"StdoutPipe": err.Error(),
			}).Warn("worker.runCmd(job *Job)")
			chErr <- true
		}

		stderr, err := c.StderrPipe()
		if err != nil {
			log.WithFields(logrus.Fields{
				"StderrPipe": err.Error(),
			}).Warn("worker.runCmd(job *Job)")
			chErr <- true
		}

		r := io.MultiReader(stdout, stderr)
		if err := c.Start(); err != nil {
			log.WithFields(logrus.Fields{
				"Start": err.Error(),
			}).Warn("worker.runCmd(job *Job)")
			chErr <- true
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
			}).Warn("worker.runCmd(job *Job)")
			chErr <- true
		}

		ok <- true
	}()

	//监听通道，超时则kill掉进程
	select {
	case <-time.After(time.Duration(job.TimeOut) * time.Millisecond):
		log.WithFields(logrus.Fields{
			"StartTime": startTime,
			"EndTime":   time.Now().Format("2006-01-02 15:04:05"),
			"JobId":     job.Id,
			"JobName":   job.Name,
			"JobCmd":    job.Cmd,
			"JobArg":    cmdArgs,
		}).Warn("worker.runCmd(job *Job) is timeout")
		c.Process.Kill()
		return
	case <-chErr:
		//异常退出
		log.WithFields(logrus.Fields{
			"StartTime": startTime,
			"EndTime":   time.Now().Format("2006-01-02 15:04:05"),
			"JobId":     job.Id,
			"JobName":   job.Name,
			"JobCmd":    job.Cmd,
			"JobArg":    cmdArgs,
		}).Info("worker.runCmd(job *Job) is error")
	case <-ok:
		//正常退出
		log.WithFields(logrus.Fields{
			"StartTime": startTime,
			"EndTime":   time.Now().Format("2006-01-02 15:04:05"),
			"JobId":     job.Id,
			"JobName":   job.Name,
			"JobCmd":    job.Cmd,
			"JobArg":    cmdArgs,
		}).Info("worker.runCmd(job *Job) is ok")
		return
	}

	return
}

func main() {
	executer := new(Executer)
	rpc.Register(executer)
	rpc.HandleHTTP()

	log.WithFields(logrus.Fields{
		"Port": gPort,
	}).Info("Server is running")

	err := http.ListenAndServe(gPort, nil)

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("error!")
	}

}
