//调度模块，负责从元数据库读取并解析调度信息。
//将需要执行的任务发送给执行模块，并读取返回信息。
package schedule

import (
	"database/sql"
	"fmt"
	"github.com/Sirupsen/logrus"
	"sort"
	"time"
)

type GlobalConfigStruct struct {
	L           *logrus.Logger
	HiveConn    *sql.DB
	LogConn     *sql.DB
	Port        string
	ExecScdChan chan *ExecSchedule
	Tasks       map[int64]*Task
	ExecTasks   map[int64]*ExecTask
	Schedules   *ScheduleList
}

func DefaultGlobal() *GlobalConfigStruct {
	sc := &GlobalConfigStruct{}
	sc.L = logrus.New()
	sc.L.Formatter = new(logrus.TextFormatter) // default
	sc.L.Level = logrus.Info
	sc.Port = ":3128"
	sc.ExecScdChan = make(chan *ExecSchedule)
	sc.ExecTasks = make(map[int64]*ExecTask)
	sc.Tasks = make(map[int64]*Task)
	sc.Schedules = &ScheduleList{}
	return sc
}

//全局变量定义
var (
	g *GlobalConfigStruct
)

//ScheduleList 调度列表结构，它包含了全部的调度信息，并有两个方法来初始化和启动其中的调度。
type ScheduleList struct {
	ScheduleList map[int64]*Schedule //调度列表
}

//从元数据库获取Schedule列表
//StartSchedule方法，会遍历列表中的Schedule并启动goroutine调用它的Timer方法。
func (sl *ScheduleList) StartSchedule() { // {{{

	//从元数据库读取调度信息,初始化调度列表
	sl.ScheduleList = getAllSchedules()

	for _, scd := range sl.ScheduleList {
		//Timer方法会根据调度周期及启动时间，按时启动，随后会依据Schedule信息构建执行结构
		go scd.Timer()
	}

} // }}}

//StartSchedule函数是调度模块的入口函数。
func StartSchedule(global *GlobalConfigStruct) { // {{{
	g = global

	//执行调度
	g.Schedules.StartSchedule()

	return
} // }}}

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
	g.L.Println(s.id, s.name, "will start at", s.nextStart)
	select {
	case <-time.After(countDown):
		//刷新调度
		s.refreshSchedule()

		g.L.Println("schedule", s.id, s.name, "is start")
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
	g.L.Println("refresh schedule", s.name)
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
	g.L.Println("schedule refreshed", s)
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
	sql := `INSERT INTO scd_schedule
            (scd_id, scd_name, scd_num, scd_cyc,
             scd_timeout, scd_job_id, scd_desc, create_user_id,
             create_time, modify_user_id, modify_time)
		VALUES      (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = g.HiveConn.Exec(sql, &s.id, &s.name, &s.count, &s.cyc,
		&s.timeOut, &s.jobId, &s.desc, &s.createUserId, &s.createTime, &s.modifyUserId, &s.modifyTime)
	g.L.Debugln("schedule", s.name, " was added.")

	return err
} // }}}

//Update方法将Schedule对象更新到元数据库。
func (s *Schedule) Update() (err error) { // {{{
	sql := `UPDATE scd_schedule 
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
	_, err = g.HiveConn.Exec(sql, &s.name, &s.count, &s.cyc,
		&s.timeOut, &s.jobId, &s.desc, &s.createUserId, &s.createTime, &s.modifyUserId, &s.modifyTime, &s.id)
	g.L.Debugln("schedule", s.name, " was updated.")

	return err
} // }}}

//Delete方法，删除元数据库中的调度信息
func (s *Schedule) Delete() error { // {{{
	sql := `Delete scd_schedule WHERE scd_id=?`
	_, err := g.HiveConn.Exec(sql, &s.id)
	g.L.Debugln("schedule", s.name, " was deleted.")

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
			FROM scd_schedule scd`
	rows, err := g.HiveConn.Query(sql)
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
			FROM scd_start s
			WHERE s.scd_id=?`
	rows, err := g.HiveConn.Query(sql, id)
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
			FROM scd_schedule scd
			WHERE scd.scd_id=?`
	rows, err := g.HiveConn.Query(sql, id)
	CheckErr("getSchedule run Sql "+sql, err)

	scd = &Schedule{}
	scd.startSecond = make([]time.Duration, 0)
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&scd.id, &scd.name, &scd.count, &scd.cyc,
			&scd.timeOut, &scd.jobId, &scd.desc)
		PrintErr("get schedule info", err)
		scd.startSecond = getStart(scd.id)
		g.L.Debugln("get Schedule", scd)

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
			FROM scd_schedule scd`
	rows, err := g.HiveConn.Query(sql)
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
		g.L.Debugln("get Schedule", scd)
	}

	return scds
} // }}}
