//调度模块，负责从元数据库读取并解析调度信息。
//将需要执行的任务发送给执行模块，并读取返回信息。
package schedule

import (
	"database/sql"
	"fmt"
	"github.com/Sirupsen/logrus"
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
	Schedules   *ScheduleManager
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
	sc.Schedules = &ScheduleManager{Global: sc}
	return sc
}

//全局变量定义
var (
	g *GlobalConfigStruct
)

//ScheduleList 调度列表结构，它包含了全部的调度信息，并有两个方法来初始化和启动其中的调度。
type ScheduleManager struct {
	ScheduleList []*Schedule //调度列表
	Global       *GlobalConfigStruct
}

//从元数据库获取Schedule列表
//StartSchedule方法，会遍历列表中的Schedule并启动goroutine调用它的Timer方法。
func (sl *ScheduleManager) StartSchedule() { // {{{

	g = sl.Global
	//从元数据库读取调度信息,初始化调度列表
	sl.ScheduleList = getAllSchedules()

	for _, scd := range sl.ScheduleList {
		//Timer方法会根据调度周期及启动时间，按时启动，随后会依据Schedule信息构建执行结构
		go scd.Timer()
	}

} // }}}

//调度信息结构
type Schedule struct { // {{{
	Id           int64           //调度ID
	Name         string          //调度名称
	Count        int8            //调度次数
	Cyc          string          //调度周期
	StartSecond  []time.Duration //周期内启动时间
	StartMonth   []int           //周期内启动月份
	NextStart    time.Time       //下次启动时间
	TimeOut      int64           //最大执行时间
	JobId        int64           //作业ID
	Job          *Job            //作业
	Desc         string          //调度说明
	JobCnt       int64           //调度中作业数量
	TaskCnt      int64           //调度中任务数量
	CreateUserId int64           //创建人
	CreateTime   time.Time       //创人
	ModifyUserId int64           //修改人
	ModifyTime   time.Time       //修改时间
}

//根据调度的周期及启动时间，按时将调度传至执行列表执行。
func (s *Schedule) Timer() { // {{{

	//获取距启动的时间（秒）
	countDown, err := getCountDown(s.Cyc, s.StartMonth, s.StartSecond)
	CheckErr("getCountDown", err)

	s.NextStart = time.Now().Add(countDown)
	g.L.Println(s.Id, s.Name, "will start at", s.NextStart)
	select {
	case <-time.After(countDown):
		//刷新调度
		s.refreshSchedule()

		g.L.Println("schedule", s.Id, s.Name, "is start")
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
	g.L.Println("refresh schedule", s.Name)
	ts := getSchedule(s.Id)
	s.Name = ts.Name
	s.Count = ts.Count
	s.Cyc = ts.Cyc
	s.StartSecond = ts.StartSecond
	s.TimeOut = ts.TimeOut
	s.JobId = ts.JobId
	s.Desc = ts.Desc

	tj := getJob(s.JobId)
	tj.ScheduleId = s.Id
	tj.ScheduleCyc = s.Cyc
	tj.refreshJob()
	s.Job = tj

	s.JobCnt = 0
	s.TaskCnt = 0
	for j := s.Job; j != nil; {
		s.JobCnt++
		s.TaskCnt += j.TaskCnt
		j = j.NextJob
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
		s.Id, s.Name, s.Cyc, s.StartSecond,
		s.TimeOut, s.JobCnt, s.TaskCnt, s.NextStart, s.CreateTime, s.Desc)
} // }}}

//Add方法会将Schedule对象增加到元数据库中。
func (s *Schedule) Add() (err error) { // {{{
	s.SetNewId()
	sql := `INSERT INTO scd_schedule
            (scd_id, scd_name, scd_num, scd_cyc,
             scd_timeout, scd_job_id, scd_desc, create_user_id,
             create_time, modify_user_id, modify_time)
		VALUES      (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = g.HiveConn.Exec(sql, &s.Id, &s.Name, &s.Count, &s.Cyc,
		&s.TimeOut, &s.JobId, &s.Desc, &s.CreateUserId, &s.CreateTime, &s.ModifyUserId, &s.ModifyTime)
	g.L.Debugln("schedule", s.Name, " was added.")

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
	_, err = g.HiveConn.Exec(sql, &s.Name, &s.Count, &s.Cyc,
		&s.TimeOut, &s.JobId, &s.Desc, &s.CreateUserId, &s.CreateTime, &s.ModifyUserId, &s.ModifyTime, &s.Id)
	g.L.Debugln("schedule", s.Name, " was updated.")

	return err
} // }}}

//Delete方法，删除元数据库中的调度信息
func (s *Schedule) Delete() error { // {{{
	sql := `Delete scd_schedule WHERE scd_id=?`
	_, err := g.HiveConn.Exec(sql, &s.Id)
	g.L.Debugln("schedule", s.Name, " was deleted.")

	return err
} // }}}

//SetJob方法，设置调度下的Job
func (s *Schedule) SetJob(jobid int64) { // {{{
	s.JobId = jobid
	s.Job = getJob(jobid)
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
	s.Id = id + 1

	return

} // }}}// }}}

//getStart，从元数据库获取指定Schedule的启动时间。
func (s *Schedule) setStart() { // {{{

	s.StartSecond = make([]time.Duration, 0)
	s.StartMonth = make([]int, 0)

	//查询全部schedule启动时间列表
	sql := `SELECT s.scd_start,s.scd_start_month
			FROM scd_start s
			WHERE s.scd_id=?`
	rows, err := g.HiveConn.Query(sql, s.Id)
	CheckErr("setStart run Sql "+sql, err)

	for rows.Next() {
		var td int64
		var tm int
		err = rows.Scan(&td, &tm)
		PrintErr("get schedule start", err)
		s.StartSecond = append(s.StartSecond, time.Duration(td)*time.Second)
		if tm > 0 {
			//DB中存储的Start_month是指第几月，但后续对年周期进行时间运算时，会从每年1月开始加，所以这里先减去1个月
			tm -= 1
		}
		s.StartMonth = append(s.StartMonth, tm)
	}

	//若没有查到Schedule的启动时间，则赋默认值。
	if len(s.StartSecond) == 0 {
		s.StartSecond = append(s.StartSecond, time.Duration(0))
		s.StartMonth = append(s.StartMonth, int(0))
	}

	//排序时间
	s.sortStart()
} // }}}

//启动时间排序
//算法选择排序
func (s *Schedule) sortStart() {
	var i, j, k int

	for i = 0; i < len(s.StartMonth); i++ {
		k = i

		for j = i + 1; j < len(s.StartMonth); j++ {
			if s.StartMonth[j] < s.StartMonth[k] {
				k = j
			} else if (s.StartMonth[j] == s.StartMonth[k]) && (s.StartSecond[j] < s.StartSecond[k]) {
				k = j
			}
		}

		if k != i {
			s.StartMonth[k], s.StartMonth[i] = s.StartMonth[i], s.StartMonth[k]
			s.StartSecond[k], s.StartSecond[i] = s.StartSecond[i], s.StartSecond[k]
		}

	}

}

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
	scd.StartSecond = make([]time.Duration, 0)
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&scd.Id, &scd.Name, &scd.Count, &scd.Cyc,
			&scd.TimeOut, &scd.JobId, &scd.Desc)
		PrintErr("get schedule info", err)
		scd.setStart()
		g.L.Debugln("get Schedule", scd)

	}

	return scd
} // }}}

//从元数据库获取Schedule列表。
func getAllSchedules() (scds []*Schedule) { // {{{
	scds = make([]*Schedule, 0)

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
		scd.StartSecond = make([]time.Duration, 0)
		err = rows.Scan(&scd.Id, &scd.Name, &scd.Count, &scd.Cyc, &scd.TimeOut,
			&scd.JobId, &scd.Desc, &scd.CreateUserId, &scd.CreateTime, &scd.ModifyUserId,
			&scd.ModifyTime)
		PrintErr("get schedule info", err)
		scd.setStart()

		scds = append(scds, scd)
		g.L.Debugln("get Schedule", scd)
	}

	return scds
} // }}}
