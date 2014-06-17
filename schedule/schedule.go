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
	"sort"
	"time"
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
	hcfg := LoadHiveConfig("hive.toml")

	runtime.GOMAXPROCS(hcfg.Maxprocs)

	//设置log模块的默认格式
	l.Formatter = new(logrus.TextFormatter) // default
	l.Level = logrus.Level(hcfg.Loglevel)

	//从配置文件中获取数据库连接、服务端口号等信息
	gPort = hcfg.Port

	gExecTasks = make(map[int64]*ExecTask)
	gTasks = make(map[int64]*Task)

	dbString = hcfg.Conn
} // }}}

//ScheduleList 调度列表结构，它包含了全部的调度信息，并有两个方法来初始化和启动其中的调度。
type ScheduleList struct {
	Schedules map[int64]*Schedule //调度列表
}

//从元数据库获取Schedule列表
//StartSchedule方法，会遍历列表中的Schedule并启动goroutine调用它的Timer方法。
func (sl *ScheduleList) StartSchedule() { // {{{

	//从元数据库读取调度信息,初始化调度列表
	sl.Schedules = getAllSchedules()

	for _, scd := range sl.Schedules {
		//Timer方法会根据调度周期及启动时间，按时启动，随后会依据Schedule信息构建执行结构
		go scd.Timer()
	}

} // }}}

//StartSchedule函数是调度模块的入口函数。
func StartSchedule() error { // {{{
	// 连接数据库
	cnn, err := sql.Open("mysql", dbString)
	CheckErr("StartSchedule ", err)
	gDbConn = cnn
	defer gDbConn.Close()

	//创建并初始化调度列表
	sLst := &ScheduleList{}

	//执行调度
	sLst.StartSchedule()

	s := make(chan int64)
	<-s

	return nil
} // }}}

func main() {
	StartSchedule()
}

//调度信息结构
type Schedule struct { // {{{
	id           int64           //调度ID
	name         string          //调度名称
	count        int8            //调度次数
	cyc          string          //调度周期
	startSecond  []time.Duration //周期内启动时间
	nextStart    time.Time       //周期内启动时间
	timeOut      int64           //最大执行时间
	jobId        int64           //作业ID
	job          *Job            //作业
	desc         string          //调度说明
	jobCnt       int64           //调度中作业数量
	taskCnt      int64           //调度中任务数量
	createUserId int64           //创建人
	createTime   time.Time       //创人
	modifyUserId int64           //修改人
	modifyTime   time.Time       //修改时间
}

//根据调度的周期及启动时间，按时将调度传至执行列表执行。
func (s *Schedule) Timer() { // {{{

	//获取距启动的时间（秒）
	countDown, err := getCountDown(s.cyc, s.startSecond)
	CheckErr("getCountDown", err)

	s.nextStart = time.Now().Add(countDown)
	l.Println(s.id, s.name, "will start at", s.nextStart)
	select {
	case <-time.After(countDown):
		//刷新调度
		s.refreshSchedule()

		l.Println("schedule", s.id, s.name, "is start")
		//启动一个线程开始构建执行结构链
		es, err := NewExecSchedule(s)
		CheckErr("New ExecSchedule", err)
		//启动线程执行调度任务
		go es.Run()
	}
	return
} // }}}

//refreshSchedule方法用来从元数据库刷新调度信息
func (s *Schedule) refreshSchedule() { // {{{
	l.Println("refresh schedule", s.name)
	ts := getSchedule(s.id)
	s.name = ts.name
	s.count = ts.count
	s.cyc = ts.cyc
	s.startSecond = ts.startSecond
	s.timeOut = ts.timeOut
	s.jobId = ts.jobId
	s.desc = ts.desc

	tj := getJob(s.jobId)
	tj.scheduleId = s.id
	tj.scheduleCyc = s.cyc
	tj.refreshJob()
	s.job = tj

	s.jobCnt = 0
	s.taskCnt = 0
	for j := s.job; j != nil; {
		s.jobCnt++
		s.taskCnt += j.taskCnt
		j = j.nextJob
	}
	l.Println("schedule refreshed", s)
} // }}}

//打印Schedule结构信息
func (s *Schedule) String() string { // {{{
	return fmt.Sprintf("{id=%d"+
		" name=%s"+
		" cyc=%s"+
		" startSecond=%v"+
		" timeout=%d"+
		" jobCnt=%d"+
		" taskCnt=%d"+
		" nextStart=%v"+
		" createTime=%v"+
		" desc=%s}\n",
		s.id, s.name, s.cyc, s.startSecond,
		s.timeOut, s.jobCnt, s.taskCnt, s.nextStart, s.createTime, s.desc)
} // }}}

//Add方法会将Schedule对象增加到元数据库中。
func (s *Schedule) Add() (err error) { // {{{
	s.SetNewId()
	sql := `INSERT INTO hive.scd_schedule
            (scd_id, scd_name, scd_num, scd_cyc,
             scd_timeout, scd_job_id, scd_desc, create_user_id,
             create_time, modify_user_id, modify_time)
		VALUES      (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = gDbConn.Exec(sql, &s.id, &s.name, &s.count, &s.cyc,
		&s.timeOut, &s.jobId, &s.desc, &s.createUserId, &s.createTime, &s.modifyUserId, &s.modifyTime)
	l.Debugln("schedule", s.name, " was added.")

	return err
} // }}}

//Update方法将Schedule对象更新到元数据库。
func (s *Schedule) Update() (err error) { // {{{
	sql := `UPDATE hive.scd_schedule 
		SET  scd_name=?,
             scd_num=?,
             scd_cyc=?,
             scd_timeout=?,
             scd_job_id=?,
             scd_desc=?,
             create_user_id=?,
             create_time=?,
             modify_user_id=?,
             modify_time=?
		 WHERE scd_id=?`
	_, err = gDbConn.Exec(sql, &s.name, &s.count, &s.cyc,
		&s.timeOut, &s.jobId, &s.desc, &s.createUserId, &s.createTime, &s.modifyUserId, &s.modifyTime, &s.id)
	l.Debugln("schedule", s.name, " was updated.")

	return err
} // }}}

//Delete方法，删除元数据库中的调度信息
func (s *Schedule) Delete() error { // {{{
	sql := `Delete hive.scd_schedule WHERE scd_id=?`
	_, err := gDbConn.Exec(sql, &s.id)
	l.Debugln("schedule", s.name, " was deleted.")

	return err
} // }}}

//SetJob方法，设置调度下的Job
func (s *Schedule) SetJob(jobid int64) { // {{{
	s.jobId = jobid
	s.job = getJob(jobid)
	return
} // }}}

//SetNewId方法，检索元数据库返回新的Schedule Id
func (s *Schedule) SetNewId() { // {{{
	var id int64

	//查询全部schedule列表
	sql := `SELECT max(scd.scd_id) as scd_id
			FROM hive.scd_schedule scd`
	rows, err := gDbConn.Query(sql)
	CheckErr("SetNewId run Sql "+sql, err)

	for rows.Next() {
		err = rows.Scan(&id)
		CheckErr("get schedule new id", err)
	}
	s.id = id + 1

	return

} // }}}// }}}

//getStart，从元数据库获取指定Schedule的启动时间。
func getStart(id int64) (st []time.Duration) { // {{{

	st = make([]time.Duration, 0)

	//查询全部schedule启动时间列表
	sql := `SELECT s.scd_start
			FROM hive.scd_start s
			WHERE s.scd_id=?`
	rows, err := gDbConn.Query(sql, id)
	CheckErr("getStart run Sql "+sql, err)

	for rows.Next() {
		var td int64
		err = rows.Scan(&td)
		PrintErr("get schedule start", err)
		st = append(st, time.Duration(td)*time.Second)
	}

	//若没有查到Schedule的启动时间，则赋默认值。
	if len(st) == 0 {
		st = append(st, time.Duration(0))
	}

	sort.Sort(timeSort(st))
	return st
} // }}}

//time.Duration列表排序
type timeSort []time.Duration // {{{

func (a timeSort) Len() int           { return len(a) }
func (a timeSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a timeSort) Less(i, j int) bool { return a[i] < a[j] } // }}}

//getSchedule，从元数据库获取指定的Schedule信息。
func getSchedule(id int64) (scd *Schedule) { // {{{

	//查询全部schedule列表
	sql := `SELECT scd.scd_id,
				scd.scd_name,
				scd.scd_num,
				scd.scd_cyc,
				scd.scd_timeout,
				scd.scd_job_id,
				scd.scd_desc
			FROM hive.scd_schedule scd
			WHERE scd.scd_id=?`
	rows, err := gDbConn.Query(sql, id)
	CheckErr("getSchedule run Sql "+sql, err)

	scd = &Schedule{}
	scd.startSecond = make([]time.Duration, 0)
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&scd.id, &scd.name, &scd.count, &scd.cyc,
			&scd.timeOut, &scd.jobId, &scd.desc)
		PrintErr("get schedule info", err)
		scd.startSecond = getStart(scd.id)
		l.Debugln("get Schedule", scd)

	}

	return scd
} // }}}

//从元数据库获取Schedule列表。
func getAllSchedules() (scds map[int64]*Schedule) { // {{{
	scds = make(map[int64]*Schedule)

	//查询全部schedule列表
	sql := `SELECT scd.scd_id,
				scd.scd_name,
				scd.scd_num,
				scd.scd_cyc,
				scd.scd_timeout,
				scd.scd_job_id,
				scd.scd_desc,
				scd.create_user_id,
				scd.create_time,
				scd.modify_user_id,
				scd.modify_time
			FROM hive.scd_schedule scd`
	rows, err := gDbConn.Query(sql)
	CheckErr("getAllSchedules run Sql "+sql, err)

	for rows.Next() {
		scd := &Schedule{}
		scd.startSecond = make([]time.Duration, 0)
		err = rows.Scan(&scd.id, &scd.name, &scd.count, &scd.cyc, &scd.timeOut,
			&scd.jobId, &scd.desc, &scd.createUserId, &scd.createTime, &scd.modifyUserId,
			&scd.modifyTime)
		PrintErr("get schedule info", err)
		scd.startSecond = getStart(scd.id)

		scds[scd.id] = scd
		l.Debugln("get Schedule", scd)
	}

	return scds
} // }}}
