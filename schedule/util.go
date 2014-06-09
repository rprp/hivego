package main

import (
	"fmt"
)

func checkErr(err error) { // {{{
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
} // }}}

//打印调度信息
func printSchedule(scds []*Schedule) { // {{{
	for _, scd := range scds {
		fmt.Println(scd.name, "\tjobs=", scd.jobCnt, " tasks=", scd.taskCnt)
		//打印调度中的作业信息
		for j := scd.job; j != nil; {
			fmt.Println("\t--------------------------------------")
			fmt.Println("\t", j.name)
			//打印作业中的任务信息
			for _, t := range j.tasks {
				fmt.Println("\t\t", t.Name)

				fmt.Print("\t\t\t[")
				//打印任务依赖链
				for _, r := range t.RelTasks {
					fmt.Print(r.Name, ",")

				}
				fmt.Print("]\n")
			}
			fmt.Print("\n")
			j = j.nextJob

		}

	}
} // }}}
