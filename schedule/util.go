package main

import (
	"fmt"
	"time"
)

//获取距启动的时间（秒）
func getCountDown(cyc string, ss time.Duration) (countDown time.Duration, err error) { // {{{
	now := GetNow()
	var startTime time.Time
	//解析周期并取得距下一周期的时间
	switch {
	case cyc == "ss":
		//按秒取整
		s := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(),
			now.Second(), 0, time.Local).Add(time.Second)
		startTime = s.Add(ss)
	case cyc == "mi":
		//按分钟取整
		mi := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0,
			0, time.Local).Add(time.Minute)
		startTime = mi.Add(ss)
	case cyc == "h":
		//按小时取整
		h := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0,
			time.Local).Add(time.Hour)
		startTime = h.Add(ss)
	case cyc == "d":
		//按日取整
		d := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0,
			time.Local).AddDate(0, 0, 1)
		startTime = d.Add(ss)
	case cyc == "m":
		//按月取整
		m := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
		startTime = m.Add(ss)
	case cyc == "w":
		//按周取整
		w := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()), 0, 0, 0, 0, time.Local).AddDate(0, 0, 7)
		startTime = w.Add(ss)
	case cyc == "q":
		//回头再处理
	case cyc == "y":
		//按年取整
		y := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local).AddDate(1, 0, 0)
		startTime = y.Add(ss)
	}

	countDown = startTime.Sub(time.Now())

	return countDown, nil

} // }}}

//获取当前时间
func GetNow() time.Time { // {{{
	return time.Now().Local()
} // }}}

func checkErr(err error) { // {{{
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
} // }}}

//打印调度信息
func printSchedule(scds map[int64]*Schedule) { // {{{
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
