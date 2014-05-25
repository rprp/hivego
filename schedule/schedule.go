//package worker
package main

import (
	"fmt"
	"net/rpc"
	"os"
	_ "runtime"
)

// 监听端口号
var gPort string = ":8123"

// 任务信息结构
type Job struct {
	Id      int64             // job的id
	Name    string            // 名称
	Type    string            // 类型
	Cmd     string            // job执行的命令或脚本、函数名等。
	TimeOut int64             // 设定超时时间，0表示不做超时限制。单位秒
	Param   map[string]string // Job的参数信息
}

type Reply struct {
	Stdout string
}

//RPC结构
//服务端处理部分，接受client端发送的指令。
type Executer struct {
}

func main() {
	param := make(map[string]string)

	fmt.Println(os.Args)
	cmd := os.Args[1]
	for i := 2; i < len(os.Args); i++ {
		param[string(i)] = os.Args[i]
	}

	job := &Job{
		Id:      1234,
		Name:    "first job",
		Type:    "1",
		Cmd:     cmd,
		TimeOut: 10000,
		Param:   param,
	}
	rl := new(Reply)
	client, err := rpc.DialHTTP("tcp", "localhost"+gPort)
	err = client.Call("Executer.Run", *job, &rl)

	fmt.Println(rl.Stdout)

	if err != nil {
		fmt.Println(err.Error())
	}

}
