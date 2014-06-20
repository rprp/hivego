package schedule

import (
	"fmt"
	"time"
)

//作业信息结构
type Job struct {
	id           int64           //作业ID
	scheduleId   int64           //调度ID
	scheduleCyc  string          //调度周期
	name         string          //作业名称
	desc         string          //作业说明
	preJobId     int64           //上级作业ID
	preJob       *Job            //上级作业
	nextJobId    int64           //下级作业ID
	nextJob      *Job            //下级作业
	tasks        map[int64]*Task //作业中的任务
	taskCnt      int64           //调度中任务数量
	createUserId int64           //创建人
	createTime   time.Time       //创人
	modifyUserId int64           //修改人
	modifyTime   time.Time       //修改时间
}

//refreshJob方法用来从元数据库刷新作业信息
func (j *Job) refreshJob() { // {{{
	g.L.Println("refresh job", j.name)
	tj := getJob(j.id)
	j.name = tj.name
	j.desc = tj.desc
	j.preJobId = tj.preJobId
	j.nextJobId = tj.nextJobId
	j.nextJob = tj.nextJob
	j.tasks = make(map[int64]*Task)
	j.taskCnt = 0

	pj := getJob(j.preJobId)
	j.preJob = pj

	t := getTasks(j.id)
	j.tasks = t
	for _, tt := range t {
		tt.ScheduleCyc = j.scheduleCyc
		j.taskCnt++
		g.L.Infoln("create task", tt.Name)
		tt.refreshTask(j.id)
	}

	//获取下级任务
	if nj := getJob(j.nextJobId); nj.id != 0 {
		nj.scheduleId = j.scheduleId
		nj.scheduleCyc = j.scheduleCyc
		nj.refreshJob()
		j.nextJob = nj
	}
	g.L.Println("job refreshed", j)
} // }}}

//打印job结构信息
func (j *Job) String() string { // {{{
	var preName, nextName string
	if j.preJob != nil {
		preName = j.preJob.name
	}

	if j.nextJob != nil {
		nextName = j.nextJob.name
	}

	t1 := make([]string, 1)
	for _, t := range j.tasks {
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
		j.id,
		j.name,
		j.desc,
		preName,
		nextName,
		t1,
		j.taskCnt,
		j.createTime,
		j.modifyTime)

} // }}}

//增加作业信息至元数据库
func (j *Job) Add() (err error) { // {{{
	j.SetNewId()
	sql := `INSERT INTO scd_job
            (job_id, job_name, job_desc, prev_job_id,
             next_job_id, create_user_id, create_time,
             modify_user_id, modify_time)
		VALUES      (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = g.HiveConn.Exec(sql, &j.id, &j.name, &j.desc, &j.preJobId, &j.nextJobId, &j.createUserId, &j.createTime, &j.modifyUserId, &j.modifyTime)
	if err == nil {
		for i, t := range j.tasks {
			j.AddTask(t.Id, i)
		}
	}

	return err

} // }}}

//增加作业任务映射关系至元数据库
func (j *Job) AddTask(taskid int64, taskno int64) (err error) { // {{{
	jobtaskid, _ := j.GetNewJobTaskId()

	sql := `INSERT INTO scd_job_task
            (job_task_id, job_id, task_id, job_task_no,
             create_user_id, create_time)
		VALUES      (?, ?, ?, ?, ?, ?)`

	_, err = g.HiveConn.Exec(sql, &jobtaskid, &j.id, &taskid, &taskno, &j.createUserId, &j.createTime)

	return err
} // }}}

//删除作业任务映射关系至元数据库
func (j *Job) DeleteTask(taskid int64) (err error) { // {{{

	sql := `DELETE scd_job_task WHERE job_id=? and task_id=?`

	_, err = g.HiveConn.Exec(sql, &j.id, &taskid)

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
	_, err = g.HiveConn.Exec(sql, &j.name, &j.desc, &j.preJobId, &j.nextJobId, &j.modifyUserId, &j.modifyTime)
	return err
} // }}}

//删除作业信息至元数据库
func (j *Job) Delete() (err error) { // {{{

	sql := `DELETE scd_job_task WHERE job_id=?`
	_, err = g.HiveConn.Exec(sql, &j.id)

	sql = `DELETE scd_job WHERE job_id=?`
	_, err = g.HiveConn.Exec(sql, &j.id)

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
	j.id = id + 1

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

	job = &Job{}
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&job.id, &job.name, &job.desc, &job.preJobId, &job.nextJobId)
		CheckErr("getJob ", err)
		//初始化Task内存
		job.tasks = make(map[int64]*Task)
		g.L.Debugln("get job", job)
	}

	return job
} // }}}

//从元数据库获取Schedule下的Job列表。
func getAllJobs() (jobs map[int64]*Job, err error) { // {{{

	jobs = make(map[int64]*Job)

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
		err = rows.Scan(&job.id, &job.name, &job.desc, &job.preJobId, &job.nextJobId)

		//初始化Task内存
		job.tasks = make(map[int64]*Task)
		jobs[job.id] = job
	}

	return jobs, err
} // }}}
