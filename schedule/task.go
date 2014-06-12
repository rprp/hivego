//调度模块的数据结构
//package schedule
package main

// 任务信息结构
type Task struct {
	Id          int64             // 任务的ID
	Address     string            // 任务的执行地址
	Name        string            // 任务名称
	JobType     string            // 任务类型
	Cyc         string            //调度周期
	StartSecond int64             //周期内启动时间
	Cmd         string            // 任务执行的命令或脚本、函数名等。
	TimeOut     int64             // 设定超时时间，0表示不做超时限制。单位秒
	Param       map[string]string // 任务的参数信息
	Attr        map[string]string // 任务的属性信息
	JobId       int64             //所属作业ID
	RelTasks    map[int64]*Task   //依赖的任务
	RelTaskCnt  int64             //依赖的任务数量
}

//refreshTask方法用来从元数据库刷新Task的信息
func (t *Task) refreshTask(jobid int64) { // {{{
	if tt, ok := getTask(t.Id); ok {
		t.Address = tt.Address
		t.Name = tt.Name
		t.TimeOut = tt.TimeOut
		t.JobType = tt.JobType
		t.Cyc = tt.Cyc
		t.StartSecond = tt.StartSecond
		t.Cmd = tt.Cmd
		t.Param = tt.Param
		t.JobId = jobid
		t.Attr = tt.Attr
		t.RelTasks = make(map[int64]*Task)
		t.RelTaskCnt = 0

		if reltask, ok := getRelTaskId(t.Id); ok {
			for _, rtid := range reltask {
				t.RelTasks[rtid] = gTasks[rtid]
				t.RelTaskCnt++
			}
		}

	}

} // }}}

// 任务依赖结构
type RelTask struct {
	taskId    int64 //任务ID
	reltaskId int64 //依赖任务ID
}

//从元数据库获取任务参数信息
func getTaskParam(id int64) (taskParam map[string]string, err error) { // {{{

	taskParam = make(map[string]string)

	//查询指定的Task属性列表
	sql := `SELECT pm.task_id,
				   pm.scd_param_name,
				   pm.scd_param_value
			FROM   hive.scd_task_param pm
			WHERE pm.task_id=?`

	rows, err := gDbConn.Query(sql, id)
	checkErr(err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		var id int64
		var name, value string
		err = rows.Scan(&id, &name, &value)
		taskParam[name] = value
	}
	return taskParam, err
} // }}}

//从元数据库获取Job下的Task列表。
func getTaskAttr(id int64) (taskAttr map[string]string, err error) { // {{{

	taskAttr = make(map[string]string)

	//查询指定的Task属性列表
	sql := `SELECT ta.task_attr_name,
			   ta.task_attr_value
			FROM   hive.scd_task_attr ta
			WHERE  task_id = ?`

	rows, err := gDbConn.Query(sql, id)
	checkErr(err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		var name, value string
		err = rows.Scan(&name, &value)
		taskAttr[name] = value
	}
	return taskAttr, err
} // }}}

//从元数据库获取Task信息。
func getTask(id int64) (task *Task, ok bool) { // {{{

	//查询全部Task列表
	sql := `SELECT task.task_id,
			   task.task_address,
			   task.task_name,
			   task.task_time_out,
			   task.task_type_id,
			   task.task_cyc,
			   task.task_start,
			   task.task_cmd
			FROM hive.scd_task task
			WHERE task.task_id=?`

	rows, err := gDbConn.Query(sql, id)
	checkErr(err)

	task = &Task{}
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&task.Id, &task.Address, &task.Name, &task.TimeOut, &task.JobType, &task.Cyc, &task.StartSecond, &task.Cmd)
		if err == nil {
			ok = true
		}
		//初始化relTask、param的内存
		task.RelTasks = make(map[int64]*Task)
		task.Param = make(map[string]string)
		task.Attr = make(map[string]string)
		task.Attr, err = getTaskAttr(task.Id)
		task.Param, err = getTaskParam(task.Id)
		checkErr(err)

	}
	return task, ok
} // }}}

//从元数据库获取Job下的Task列表。
func getAllTasks() (tasks map[int64]*Task, err error) { // {{{

	tasks = make(map[int64]*Task)

	//查询全部Task列表
	sql := `SELECT task.task_id,
			   task.task_address,
			   task.task_name,
			   task.task_time_out,
			   task.task_type_id,
			   task.task_cyc,
			   task.task_start,
			   task.task_cmd
			FROM hive.scd_task task`

	rows, err := gDbConn.Query(sql)
	checkErr(err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		task := &Task{}
		err = rows.Scan(&task.Id, &task.Address, &task.Name, &task.TimeOut, &task.JobType, &task.Cyc, &task.StartSecond, &task.Cmd)
		//初始化relTask、param的内存
		task.RelTasks = make(map[int64]*Task)
		task.Param = make(map[string]string)
		task.Attr = make(map[string]string)
		task.Attr, err = getTaskAttr(task.Id)
		task.Param, err = getTaskParam(task.Id)
		checkErr(err)

		tasks[task.Id] = task
	}
	return tasks, err
} // }}}

//从元数据库获取Job下的Task列表。
func getTasks(jobId int64) (tasks map[int64]*Task, ok bool) { // {{{

	tasks = make(map[int64]*Task)

	//查询Job中全部Task列表
	sql := `SELECT jt.task_id
			FROM hive.scd_job_task jt
			WHERE jt.job_id=?`

	rows, err := gDbConn.Query(sql, jobId)
	checkErr(err)

	//循环读取记录
	for rows.Next() {
		var taskid int64
		err = rows.Scan(&taskid)
		if err == nil {
			ok = true
			if task, ok := getTask(taskid); ok {
				tasks[taskid] = task
				gTasks[taskid] = task
			}
		}
	}
	return tasks, ok
} // }}}

//从元数据库获取Job下的Task列表。
func getJobTaskid() (jobtask map[int64]int64, err error) { // {{{

	jobtask = make(map[int64]int64)

	//查询Job中全部Task列表
	sql := `SELECT jt.job_id,
				jt.task_id
			FROM hive.scd_job_task jt`

	rows, err := gDbConn.Query(sql)
	checkErr(err)

	//循环读取记录
	for rows.Next() {
		var jobid, taskid int64
		err = rows.Scan(&jobid, &taskid)
		jobtask[taskid] = jobid
	}
	return jobtask, err
} // }}}

//从元数据库获取Task的依赖列表。
func getRelTaskId(id int64) (relTaskId []int64, ok bool) { // {{{

	relTaskId = make([]int64, 0)

	//查询Task的依赖列表
	sql := `SELECT tr.rel_task_id
			FROM hive.scd_task_rel tr
			Where tr.task_id=?`

	rows, err := gDbConn.Query(sql, id)
	checkErr(err)

	//循环读取记录
	for rows.Next() {
		var rtid int64
		err = rows.Scan(&rtid)
		if err == nil {
			ok = true
		}
		relTaskId = append(relTaskId, rtid)
	}
	return relTaskId, ok
} // }}}

//从元数据库获取Task的依赖列表。
func getRelTasks() (relTasks []*RelTask, err error) { // {{{

	relTasks = make([]*RelTask, 0)

	//查询Task的依赖列表
	sql := `SELECT tr.task_id,
				tr.rel_task_id
			FROM hive.scd_task_rel tr
			ORDER BY tr.task_id`

	rows, err := gDbConn.Query(sql)
	checkErr(err)

	//循环读取记录
	for rows.Next() {
		var taskid, reltaskid int64
		err = rows.Scan(&taskid, &reltaskid)
		relTasks = append(relTasks, &RelTask{taskId: taskid, reltaskId: reltaskid})
	}
	return relTasks, err
} // }}}
