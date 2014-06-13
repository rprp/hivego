//调度模块的数据结构
//package schedule
package main

//作业信息结构
type Job struct {
	id        int64           //作业ID
	name      string          //作业名称
	desc      string          //作业说明
	preJobId  int64           //上级作业ID
	preJob    *Job            //上级作业
	nextJobId int64           //下级作业ID
	nextJob   *Job            //下级作业
	tasks     map[int64]*Task //作业中的任务
	taskCnt   int64           //调度中任务数量
}

//refreshJob方法用来从元数据库刷新作业信息
func (j *Job) refreshJob() { // {{{

	if tj, ok := getJob(j.id); ok {
		j.name = tj.name
		j.desc = tj.desc
		j.preJobId = tj.preJobId
		j.nextJobId = tj.nextJobId
		j.nextJob = tj.nextJob
		j.tasks = make(map[int64]*Task)
		j.taskCnt = 0

		l.Infoln("create job", j.name)
		if pj, ok := getJob(j.preJobId); ok {
			j.preJob = pj
		}

		if t, ok := getTasks(j.id); ok {
			j.tasks = t
			for _, tt := range t {
				j.taskCnt++
				l.Infoln("create task", tt.Name)
				tt.refreshTask(j.id)
			}
		}

		if nj, ok := getJob(j.nextJobId); ok {
			nj.refreshJob()
			j.nextJob = nj
		}
	}
} // }}}

//从元数据库获取Job信息。
func getJob(id int64) (job *Job, ok bool) { // {{{

	//查询全部Job列表
	sql := `SELECT job.job_id,
			   job.job_name,
			   job.job_desc,
			   job.prev_job_id,
			   job.next_job_id
			FROM hive.scd_job job
			WHERE job.job_id=?`

	rows, err := gDbConn.Query(sql, id)
	checkErr(err)

	job = &Job{}
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&job.id, &job.name, &job.desc, &job.preJobId, &job.nextJobId)
		if err == nil {
			ok = true
		}
		//初始化Task内存
		job.tasks = make(map[int64]*Task)
	}

	return job, ok
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
			FROM hive.scd_job job`

	rows, err := gDbConn.Query(sql)
	checkErr(err)

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
