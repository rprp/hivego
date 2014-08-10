package schedule

import (
	"fmt"
	"time"
)

// 任务信息结构
type Task struct { // {{{
	Id           int64             // 任务的ID
	Address      string            // 任务的执行地址
	Name         string            // 任务名称
	TaskType     int64             // 任务类型
	ScheduleCyc  string            //调度周期
	TaskCyc      string            //调度周期
	StartSecond  time.Duration     //周期内启动时间
	Cmd          string            // 任务执行的命令或脚本、函数名等。
	Desc         string            //任务说明
	TimeOut      int64             // 设定超时时间，0表示不做超时限制。单位秒
	Param        []string          // 任务的参数信息
	Attr         map[string]string // 任务的属性信息
	JobId        int64             //所属作业ID
	RelTasksId   []int64           //依赖的任务Id
	RelTasks     map[string]*Task  //`json:"-"` //依赖的任务
	RelTaskCnt   int64             //依赖的任务数量
	CreateUserId int64             //创建人
	CreateTime   time.Time         //创人
	ModifyUserId int64             //修改人
	ModifyTime   time.Time         //修改时间
} // }}}

//refreshTask方法用来从元数据库刷新Task的信息
func (t *Task) refreshTask(jobid int64) { // {{{
	g.L.Println("refresh task", t.Name)
	tt := getTask(t.Id)
	t.Address = tt.Address
	t.Name = tt.Name
	t.TimeOut = tt.TimeOut
	t.TaskType = tt.TaskType
	t.TaskCyc = tt.TaskCyc
	t.StartSecond = tt.StartSecond
	t.Cmd = tt.Cmd
	t.Param = tt.Param
	t.Desc = tt.Desc
	t.JobId = jobid
	t.Attr = tt.Attr
	t.RelTasksId = make([]int64, 0)
	t.RelTasks = make(map[string]*Task)
	t.RelTaskCnt = 0

	reltask := getRelTaskId(t.Id)
	for _, rtid := range reltask {
		t.RelTasksId = append(t.RelTasksId, rtid)
		t.RelTasks[string(rtid)] = g.Tasks[string(rtid)]
		t.RelTaskCnt++
	}

	g.L.Println("task refreshed", t)

} // }}}

//打印task结构信息
func (t *Task) String() string { // {{{

	tn := make([]string, 1)
	for _, rt := range t.RelTasks {
		tn = append(tn, rt.Name)
	}

	return fmt.Sprintf("{Id=%d"+
		" Address=%s"+
		" Name=%s"+
		" TaskCyc=%s"+
		" StartSecond=%v"+
		" Cmd=%s"+
		" Desc=%s"+
		" TimeOut=%d"+
		" Param=%v"+
		" RelTasks=%v"+
		" RelTaskCnt =%d"+
		" CreateTime=%v"+
		" ModifyTime=%v}\n",
		t.Id,
		t.Address,
		t.Name,
		t.TaskCyc,
		t.StartSecond,
		t.Cmd,
		t.Desc,
		t.TimeOut,
		t.Param,
		tn,
		t.RelTaskCnt,
		t.CreateTime,
		t.ModifyTime)

} // }}}

//UpdateTask方法调用Update方法将Task信息更新到元数据库中，
//同时更新Task的参数信息
func (t *Task) UpdateTask() (err error) { // {{{
	if err = t.Update(); err != nil {
		return err
	}

	if err = t.DelParam(); err != nil {
		return err
	}

	for _, p := range t.Param {
		fmt.Println("t.Param:", p)
		if err = t.AddParam(p); err != nil {
			return err
		}
	}

	return err
} // }}}

//DelParam方法从元数据库删除Task的Param信息
func (t *Task) DelParam() (err error) {
	sql := `DELETE FROM scd_task_param
			WHERE task_id=?`
	_, err = g.HiveConn.Exec(sql, &t.Id)

	return err
}

//AddTask方法持久化当前的Task信息。
//先调用Add方法将Task基本信息持久化，成功则依次持久化
//关联的作业信息、Task依赖关系、参数列表。
func (t *Task) AddTask() (err error) { // {{{
	if err = t.Add(); err == nil {
		if err = t.AddRelJob(); err != nil {
			fmt.Println("addRelJob", err)
			return err
		}
		for _, rt := range t.RelTasks {
			if err = t.addRelTask(rt.Id); err != nil {
				fmt.Println("addRelTask", err)
				return err
			}
		}
		for _, p := range t.Param {
			if err = t.AddParam(p); err != nil {
				fmt.Println("addParam", err)
				return err
			}
		}

	}

	return err
} // }}}

//AddRelJob将Task与Job的关系持久化。
func (t *Task) AddRelJob() (err error) {
	var id int64
	if id, err = t.GetRelJobId(); err == nil {
		sql := `INSERT INTO scd_job_task
            (job_task_id,job_id,task_id,job_task_no,
            create_user_id,create_time)
            VALUES    (?, ?, ?, ?, ?, ?)`
		_, err = g.HiveConn.Exec(sql, &id, &t.JobId, &t.Id, &t.Id, &t.CreateUserId, &t.CreateTime)
	}
	return err
}

func (t *Task) DeleteRelTask(relid int64) (err error) {
	t.RelTasksId = append(t.RelTasksId, relid)
	var i int
	for k, v := range t.RelTasksId {
		if v == relid {
			i = k
		}
	}
	t.RelTasksId = append(t.RelTasksId[0:i], t.RelTasksId[i+1:]...)
	t.RelTaskCnt--
	delete(t.RelTasks, string(relid))

	return t.deleteRelTask(relid)
}

func (t *Task) AddRelTask(rt *Task) (err error) {
	t.RelTasksId = append(t.RelTasksId, rt.Id)
	t.RelTaskCnt++
	t.RelTasks[string(rt.Id)] = rt

	return t.addRelTask(rt.Id)
}

//GetRelJobId获取最大的Id
func (t *Task) GetRelJobId() (id int64, err error) {

	//查询全部schedule列表
	sql := `SELECT max(t.job_task_id) as job_task_id
			FROM scd_job_task t`
	rows, err := g.HiveConn.Query(sql)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
	}

	return id + 1, err

}

//增加作业信息至元数据库
func (t *Task) Add() (err error) { // {{{
	t.SetNewId()
	sql := `INSERT INTO scd_task
            (task_id, task_address, task_name, task_cyc,
             task_time_out, task_start, task_type_id,
             task_cmd, task_desc, create_user_id, create_time,
             modify_user_id, modify_time)
			VALUES      (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = g.HiveConn.Exec(sql, &t.Id, &t.Address, &t.Name, &t.TaskCyc, &t.TimeOut, &t.StartSecond, &t.TaskType, &t.Cmd, &t.Desc, &t.CreateUserId, &t.CreateTime, &t.ModifyUserId, &t.ModifyTime)

	return err
} // }}}

//增加作业参数信息至元数据库
func (t *Task) AddParam(pvalue string) (err error) { // {{{
	pid, _ := t.GetNewParamTaskId()
	sql := `INSERT INTO scd_task_param
            (scd_param_id,task_id, scd_param_name, scd_param_value,
             create_user_id, create_time)
			VALUES      (?, ?, ?, ?, ?, ?)`
	_, err = g.HiveConn.Exec(sql, &pid, &t.Id, "0", &pvalue, &t.CreateUserId, &t.CreateTime)

	return err
} // }}}

//删除依赖任务至元数据库
func (t *Task) deleteRelTask(id int64) (err error) { // {{{
	sql := `DELETE FROM scd_task_rel WHERE task_id=? and rel_task_id=?`
	_, err = g.HiveConn.Exec(sql, &t.Id, &id)

	return err
} // }}}

//增加依赖任务至元数据库
func (t *Task) addRelTask(id int64) (err error) { // {{{
	relid, _ := t.GetNewRelTaskId()
	sql := `INSERT INTO scd_task_rel
            (task_rel_id, task_id, rel_task_id, create_user_id, create_time)
			VALUES      (?, ?, ?, ?, ? )`
	_, err = g.HiveConn.Exec(sql, &relid, &t.Id, &id, &t.CreateUserId, &t.CreateTime)

	return err
} // }}}

//更新任务至元数据库
func (t *Task) Update() (err error) { // {{{
	sql := `UPDATE scd_task
			SET task_address=?,
				task_name=?,
				task_cyc=?,
				task_time_out=?,
				task_start=?,
				task_type_id=?,
				task_cmd=?,
				task_desc=?,
				modify_user_id=?,
				modify_time=?
			WHERE task_id=?`
	_, err = g.HiveConn.Exec(sql, &t.Address, &t.Name, &t.TaskCyc, &t.TimeOut, &t.StartSecond, &t.TaskType, &t.Cmd, &t.Desc, &t.ModifyUserId, &t.ModifyTime, &t.Id)

	return err
} // }}}

//删除任务至元数据库
func (t *Task) Delete() (err error) { // {{{

	sql := `DELETE scd_task_param WHERE task_id=?`
	_, err = g.HiveConn.Exec(sql, &t.Id)

	sql = `DELETE scd_task_rel WHERE task_id=?`
	_, err = g.HiveConn.Exec(sql, &t.Id)

	sql = `DELETE scd_task WHERE task_id=?`
	_, err = g.HiveConn.Exec(sql, &t.Id)

	return err
} // }}}

//获取新TaskParamId
func (t *Task) GetNewParamTaskId() (id int64, err error) { // {{{

	//查询全部schedule列表
	sql := `SELECT max(p.scd_param_id) as scd_param_id
			FROM scd_task_param p`

	rows, err := g.HiveConn.Query(sql)
	CheckErr("GetNewParamTaskId run sql"+sql, err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
	}

	return id + 1, err

} // }}}

//获取新JobTaskId
func (t *Task) GetNewRelTaskId() (id int64, err error) { // {{{

	//查询全部schedule列表
	sql := `SELECT max(rt.task_rel_id) as task_rel_id
			FROM scd_task_rel rt`

	rows, err := g.HiveConn.Query(sql)
	CheckErr("GetNewRelTaskId run sql"+sql, err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
	}

	return id + 1, err

} // }}}

//获取新Id
func (t *Task) SetNewId() (err error) { // {{{
	var id int64

	//查询全部schedule列表
	sql := `SELECT max(t.task_id) as task_id
			FROM scd_task t`
	rows, err := g.HiveConn.Query(sql)
	CheckErr("SetNewId run sql"+sql, err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
	}
	t.Id = id + 1

	return err

} // }}}

// 任务依赖结构
type RelTask struct { // {{{
	taskId    int64 //任务ID
	reltaskId int64 //依赖任务ID
} // }}}

//从元数据库获取任务参数信息
func getTaskParam(id int64) (taskParam []string, err error) { // {{{

	taskParam = make([]string, 0)

	//查询指定的Task属性列表
	sql := `SELECT pm.task_id,
				   pm.scd_param_name,
				   pm.scd_param_value
			FROM   scd_task_param pm
			WHERE pm.task_id=?`

	rows, err := g.HiveConn.Query(sql, id)
	CheckErr("getTaskParam run sql"+sql, err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		var id int64
		var name, value string
		err = rows.Scan(&id, &name, &value)
		taskParam = append(taskParam, value)
	}
	return taskParam, err
} // }}}

//从元数据库获取Job下的Task列表。
func getTaskAttr(id int64) (taskAttr map[string]string, err error) { // {{{

	taskAttr = make(map[string]string)

	//查询指定的Task属性列表
	sql := `SELECT ta.task_attr_name,
			   ta.task_attr_value
			FROM   scd_task_attr ta
			WHERE  task_id = ?`
	rows, err := g.HiveConn.Query(sql, id)
	CheckErr("getTaskAttr run sql"+sql, err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		var name, value string
		err = rows.Scan(&name, &value)
		taskAttr[name] = value
	}
	return taskAttr, err
} // }}}

//从元数据库获取Task信息。
func getTask(id int64) (task *Task) { // {{{

	var td int64
	//查询全部Task列表
	sql := `SELECT task.task_id,
			   task.task_address,
			   task.task_name,
			   task.task_time_out,
			   task.task_type_id,
			   task.task_cyc,
			   task.task_desc,
			   task.task_start,
			   task.task_cmd
			FROM scd_task task
			WHERE task.task_id=?`
	rows, err := g.HiveConn.Query(sql, id)
	CheckErr("getTask run sql"+sql, err)

	task = &Task{}

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&task.Id, &task.Address, &task.Name, &task.TimeOut, &task.TaskType, &task.TaskCyc, &task.Desc, &td, &task.Cmd)
		//初始化relTask、param的内存
		task.RelTasksId = make([]int64, 0)
		task.RelTasks = make(map[string]*Task)
		task.Param = make([]string, 0)
		task.Attr = make(map[string]string)
		task.Attr, err = getTaskAttr(task.Id)
		task.Param, err = getTaskParam(task.Id)
		task.StartSecond = time.Duration(td) * time.Second
		CheckErr("getTask", err)
		g.L.Debugln("get task", task)
	}
	return task
} // }}}

//从元数据库获取Job下的Task列表。
func getAllTasks() (tasks map[string]*Task, err error) { // {{{

	tasks = make(map[string]*Task)

	//查询全部Task列表
	sql := `SELECT task.task_id,
			   task.task_address,
			   task.task_name,
			   task.task_time_out,
			   task.task_type_id,
			   task.task_cyc,
			   task.task_desc,
			   task.task_start,
			   task.task_cmd
			FROM scd_task task`
	rows, err := g.HiveConn.Query(sql)
	CheckErr("getAllTasks run sql"+sql, err)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		task := &Task{}
		var td int64
		err = rows.Scan(&task.Id, &task.Address, &task.Name, &task.TimeOut, &task.TaskType, &task.TaskCyc, &task.Desc, &td, &task.Cmd)
		//初始化relTask、param的内存
		task.RelTasksId = make([]int64, 0)
		task.RelTasks = make(map[string]*Task)
		task.Param = make([]string, 0)
		task.Attr = make(map[string]string)
		task.Attr, err = getTaskAttr(task.Id)
		task.Param, err = getTaskParam(task.Id)
		task.StartSecond = time.Duration(td) * time.Second
		CheckErr("getAllTask", err)

		tasks[string(task.Id)] = task
	}
	return tasks, err
} // }}}

//从元数据库获取Job下的Task列表。
func getTasks(jobId int64) (tasks map[string]*Task) { // {{{

	tasks = make(map[string]*Task)

	//查询Job中全部Task列表
	sql := `SELECT jt.task_id
			FROM scd_job_task jt
			WHERE jt.job_id=?`
	rows, err := g.HiveConn.Query(sql, jobId)
	CheckErr("getTasks run sql"+sql, err)

	//循环读取记录
	for rows.Next() {
		var taskid int64
		err = rows.Scan(&taskid)
		CheckErr("getTasks", err)
		if task := getTask(taskid); task.Id != 0 {
			tasks[string(taskid)] = task
			g.Tasks[string(taskid)] = task
		}
	}
	g.L.Debugln("get task", tasks)
	return tasks
} // }}}

//从元数据库获取Job下的Task列表。
func getJobTaskid() (jobtask map[string]int64, err error) { // {{{

	jobtask = make(map[string]int64)

	//查询Job中全部Task列表
	sql := `SELECT jt.job_id,
				jt.task_id
			FROM scd_job_task jt`
	rows, err := g.HiveConn.Query(sql)
	CheckErr("getJobTaskid run sql"+sql, err)

	//循环读取记录
	for rows.Next() {
		var jobid, taskid int64
		err = rows.Scan(&jobid, &taskid)
		jobtask[string(taskid)] = jobid
	}
	return jobtask, err
} // }}}

//从元数据库获取Task的依赖列表。
func getRelTaskId(id int64) (relTaskId []int64) { // {{{

	relTaskId = make([]int64, 0)

	//查询Task的依赖列表
	sql := `SELECT tr.rel_task_id
			FROM scd_task_rel tr
			Where tr.task_id=?`
	rows, err := g.HiveConn.Query(sql, id)
	CheckErr("getRelTaskId run sql"+sql, err)

	//循环读取记录
	for rows.Next() {
		var rtid int64
		err = rows.Scan(&rtid)
		CheckErr("getRelTaskId", err)
		relTaskId = append(relTaskId, rtid)
	}
	return relTaskId
} // }}}

//从元数据库获取Task的依赖列表。
func getRelTasks() (relTasks []*RelTask, err error) { // {{{

	relTasks = make([]*RelTask, 0)

	//查询Task的依赖列表
	sql := `SELECT tr.task_id,
				tr.rel_task_id
			FROM scd_task_rel tr
			ORDER BY tr.task_id`
	rows, err := g.HiveConn.Query(sql)
	CheckErr("getRelTasks run sql"+sql, err)

	//循环读取记录
	for rows.Next() {
		var taskid, reltaskid int64
		err = rows.Scan(&taskid, &reltaskid)
		relTasks = append(relTasks, &RelTask{taskId: taskid, reltaskId: reltaskid})
	}
	return relTasks, err
} // }}}
