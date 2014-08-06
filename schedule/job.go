package schedule

import (
	"errors"
	"fmt"
	"time"
)

//作业信息结构
type Job struct { // {{{
	Id           int64            //作业ID
	ScheduleId   int64            //调度ID
	ScheduleCyc  string           //调度周期
	Name         string           //作业名称
	Desc         string           //作业说明
	PreJobId     int64            //上级作业ID
	PreJob       *Job             `json:"-"` //上级作业
	NextJobId    int64            //下级作业ID
	NextJob      *Job             `json:"-"` //下级作业
	Tasks        map[string]*Task //作业中的任务
	TaskCnt      int64            //调度中任务数量
	CreateUserId int64            //创建人
	CreateTime   time.Time        //创人
	ModifyUserId int64            //修改人
	ModifyTime   time.Time        //修改时间
} // }}}

//refreshJob方法用来从元数据库刷新作业信息
func (j *Job) refreshJob() { // {{{
	g.L.Println("refresh job", j.Name)
	if tj := getJob(j.Id); tj != nil {
		j.Name = tj.Name
		j.Desc = tj.Desc
		j.PreJobId = tj.PreJobId
		j.NextJobId = tj.NextJobId
		j.NextJob = tj.NextJob
		j.Tasks = make(map[string]*Task)
		j.TaskCnt = 0

		pj := getJob(j.PreJobId)
		j.PreJob = pj

		t := getTasks(j.Id)
		j.Tasks = t
		for _, tt := range t {
			tt.ScheduleCyc = j.ScheduleCyc
			j.TaskCnt++
			g.L.Infoln("create task", tt.Name)
			tt.refreshTask(j.Id)
		}

		//获取下级任务
		if nj := getJob(j.NextJobId); nj != nil && nj.Id != 0 {
			nj.ScheduleId = j.ScheduleId
			nj.ScheduleCyc = j.ScheduleCyc
			nj.refreshJob()
			j.NextJob = nj
		}
	}
	g.L.Println("job refreshed", j)
} // }}}

//打印job结构信息
func (j *Job) String() string { // {{{
	var preName, nextName string
	if j.PreJob != nil {
		preName = j.PreJob.Name
	}

	if j.NextJob != nil {
		nextName = j.NextJob.Name
	}

	t1 := make([]string, 1)
	for _, t := range j.Tasks {
		t1 = append(t1, t.Name)
	}

	return fmt.Sprintf("{id=%d"+
		" name=%s"+
		" desc=%s"+
		" preJobname=%s"+
		" nextJobname=%s"+
		" taskList=%v"+
		" taskCnt=%d"+
		" createTime=%v"+
		" modifyTime=%v}\n",
		j.Id,
		j.Name,
		j.Desc,
		preName,
		nextName,
		t1,
		j.TaskCnt,
		j.CreateTime,
		j.ModifyTime)

} // }}}

//增加作业信息至元数据库
func (j *Job) Add() (err error) { // {{{
	j.SetNewId()
	j.Tasks = make(map[string]*Task)
	j.CreateTime = time.Now()
	j.ModifyTime = time.Now()
	sql := `INSERT INTO scd_job
            (job_id, job_name, job_desc, prev_job_id,
             next_job_id, create_user_id, create_time,
             modify_user_id, modify_time)
		VALUES      (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = g.HiveConn.Exec(sql, &j.Id, &j.Name, &j.Desc, &j.PreJobId, &j.NextJobId, &j.CreateUserId, &j.CreateTime, &j.ModifyUserId, &j.ModifyTime)

	return err
} // }}}

//增加作业任务映射关系至元数据库
func (j *Job) AddTask(taskid int64, taskno int64) (err error) { // {{{
	jobtaskid, _ := j.GetNewJobTaskId()

	sql := `INSERT INTO scd_job_task
            (job_task_id, job_id, task_id, job_task_no,
             create_user_id, create_time)
		VALUES      (?, ?, ?, ?, ?, ?)`

	_, err = g.HiveConn.Exec(sql, &jobtaskid, &j.Id, &taskid, &taskno, &j.CreateUserId, &j.CreateTime)

	return err
} // }}}

//UpdateTask更新Job中指定Task的信息。
//它会根据参数查找本Job下符合的Task，找到后更新信息
//并调用Task的Add方法进行持久化操作。
func (j *Job) UpdateTask(task *Task) (err error) {
	t, ok := j.Tasks[string(task.Id)]
	if !ok {
		return errors.New(fmt.Sprintf("update error. not found task by id %d", task.Id))
	}
	t.Name = task.Name
	t.Desc = task.Desc
	t.Address = task.Address
	t.TaskType = task.TaskType
	t.TaskCyc = task.TaskCyc
	t.StartSecond = task.StartSecond
	t.Cmd = task.Cmd
	t.TimeOut = task.TimeOut
	t.Param = task.Param
	t.Attr = task.Attr
	t.ModifyUserId = task.ModifyUserId
	t.ModifyTime = time.Now()

	if err := t.UpdateTask(); err != nil {
		return err
	}

	return nil
}

//删除作业任务映射关系至元数据库
func (j *Job) DeleteTask(taskid int64) (err error) { // {{{

	sql := `DELETE scd_job_task WHERE job_id=? and task_id=?`

	_, err = g.HiveConn.Exec(sql, &j.Id, &taskid)

	return err
} // }}}

//修改作业信息至元数据库
func (j *Job) Update() (err error) { // {{{
	sql := `UPDATE scd_job
		SET job_name=?, 
			job_desc=?,
			prev_job_id=?,
            next_job_id=?, 
            modify_user_id=?, 
			modify_time=?
	    WHERE job_id=?`
	_, err = g.HiveConn.Exec(sql, &j.Name, &j.Desc, &j.PreJobId, &j.NextJobId, &j.ModifyUserId, &j.ModifyTime, &j.Id)
	return err
} // }}}

//删除作业信息至元数据库
func (j *Job) Delete() (err error) { // {{{
	sql := `DELETE FROM scd_job WHERE job_id=?`
	_, err = g.HiveConn.Exec(sql, &j.Id)
	return err
} // }}}

//获取新Id
func (j *Job) SetNewId() (err error) { // {{{
	var id int64

	//查询全部schedule列表
	sql := `SELECT max(job.job_id) as job_id
			FROM scd_job job`
	rows, err := g.HiveConn.Query(sql)
	CheckErr("job SetNewId run Sql "+sql, err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
	}
	j.Id = id + 1

	return err

} // }}}

//获取新JobTaskId
func (j *Job) GetNewJobTaskId() (id int64, err error) { // {{{

	//查询全部schedule列表
	sql := `SELECT max(jt.job_task_id) as job_task_id
			FROM scd_job_task jt`
	rows, err := g.HiveConn.Query(sql)
	CheckErr("GetNewJobTaskId run Sql "+sql, err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
	}

	return id + 1, err

} // }}}

//从元数据库获取Job信息。
func getJob(id int64) (job *Job) { // {{{

	//查询全部Job列表
	sql := `SELECT job.job_id,
			   job.job_name,
			   job.job_desc,
			   job.prev_job_id,
			   job.next_job_id
			FROM scd_job job
			WHERE job.job_id=?`
	rows, err := g.HiveConn.Query(sql, id)
	CheckErr("getJob run Sql "+sql, err)

	j := &Job{}
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&j.Id, &j.Name, &j.Desc, &j.PreJobId, &j.NextJobId)
		CheckErr("getJob ", err)
		//初始化Task内存
		j.Tasks = make(map[string]*Task)
		g.L.Debugln("get job", j)
	}

	if j.Id == 0 {
		j = nil
	} else {
		job = j
	}
	return
} // }}}

//从元数据库获取Schedule下的Job列表。
func getAllJobs() (jobs map[string]*Job, err error) { // {{{

	jobs = make(map[string]*Job)

	//查询全部Job列表
	sql := `SELECT job.job_id,
			   job.job_name,
			   job.job_desc,
			   job.prev_job_id,
			   job.next_job_id
			FROM scd_job job`
	rows, err := g.HiveConn.Query(sql)
	CheckErr("getAllJobs run Sql "+sql, err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		job := &Job{}
		err = rows.Scan(&job.Id, &job.Name, &job.Desc, &job.PreJobId, &job.NextJobId)

		//初始化Task内存
		job.Tasks = make(map[string]*Task)
		jobs[string(job.Id)] = job
	}

	return jobs, err
} // }}}
