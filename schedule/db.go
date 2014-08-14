package schedule

import (
	"errors"
	"fmt"
	"time"
)

//Add方法会将Schedule对象增加到元数据库中。
func (s *Schedule) add() (err error) { // {{{
	if err = s.setNewId(); err != nil {
		return err
	}

	sql := `INSERT INTO scd_schedule
            (scd_id, scd_name, scd_num, scd_cyc,
             scd_timeout, scd_job_id, scd_desc, create_user_id,
             create_time, modify_user_id, modify_time)
		VALUES      (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = g.HiveConn.Exec(sql, &s.Id, &s.Name, &s.Count, &s.Cyc,
		&s.TimeOut, &s.JobId, &s.Desc, &s.CreateUserId, &s.CreateTime, &s.ModifyUserId, &s.ModifyTime)
	g.L.Debugln("[s.add] schedule", s, "\nsql=", sql)

	return err
} // }}}

//Update方法将Schedule对象更新到元数据库。
func (s *Schedule) update() error { // {{{
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
	_, err := g.HiveConn.Exec(sql, &s.Name, &s.Count, &s.Cyc,
		&s.TimeOut, &s.JobId, &s.Desc, &s.CreateUserId, &s.CreateTime, &s.ModifyUserId, &s.ModifyTime, &s.Id)
	g.L.Debugln("[s.update] schedule", s, "\nsql=", sql)

	return err
} // }}}

//Delete方法，删除元数据库中的调度信息
func (s *Schedule) deleteSchedule() error { // {{{
	sql := `Delete FROM scd_schedule WHERE scd_id=?`
	_, err := g.HiveConn.Exec(sql, &s.Id)
	g.L.Debugln("[s.deleteSchedule] schedule", s, "\nsql=", sql)

	return err
} // }}}

//setNewId方法，检索元数据库返回新的Schedule Id
func (s *Schedule) setNewId() error { // {{{
	var id int64

	//查询全部schedule列表
	sql := `SELECT max(scd.scd_id) as scd_id
			FROM scd_schedule scd`
	rows, err := g.HiveConn.Query(sql)
	if err != nil {
		e := fmt.Sprintf("[s.setNewid] Query sql [%s] error %s.\n", sql, err.Error())
		return errors.New(e)
	}

	for rows.Next() {
		err = rows.Scan(&id)
	}
	s.Id = id + 1

	return nil
} // }}}

func (s *Schedule) addStart(t time.Duration, m int) error { // {{{
	sql := `INSERT INTO scd_start 
            (scd_id, scd_start, scd_start_month,
            create_user_id, create_time)
         VALUES  (?, ?, ?, ?, ?)`
	_, err := g.HiveConn.Exec(sql, &s.Id, &t, &m, &s.ModifyUserId, &s.ModifyTime)
	if err != nil {
		e := fmt.Sprintf("[s.addStart] Exec sql [%s] error %s.\n", sql, err.Error())
		return errors.New(e)
	}
	g.L.Debugln("[s.addStart] ", "\nsql=", sql)
	return nil
} // }}}

//delStart删除该Schedule的所有启动时间列表
func (s *Schedule) delStart() error { // {{{
	sql := `DELETE FROM scd_start WHERE scd_id=?`
	_, err := g.HiveConn.Exec(sql, &s.Id)
	if err != nil {
		e := fmt.Sprintf("[s.delStart] Exec sql [%s] error %s.\n", sql, err.Error())
		return errors.New(e)
	}
	g.L.Debugln("[s.delStart] ", "\nsql=", sql)

	return nil
} // }}}

//getStart，从元数据库获取指定Schedule的启动时间。
func (s *Schedule) setStart() error { // {{{

	s.StartSecond = make([]time.Duration, 0)
	s.StartMonth = make([]int, 0)

	//查询全部schedule启动时间列表
	sql := `SELECT s.scd_start,s.scd_start_month
			FROM scd_start s
			WHERE s.scd_id=?`
	rows, err := g.HiveConn.Query(sql, s.Id)
	if err != nil {
		e := fmt.Sprintf("[s.setStart] Exec sql [%s] error %s.\n", sql, err.Error())
		return errors.New(e)
	}
	g.L.Debugln("[s.setStart] ", "\nsql=", sql)

	for rows.Next() {
		var td int64
		var tm int
		err = rows.Scan(&td, &tm)
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
	return nil
} // }}}

//getSchedule，从元数据库获取指定的Schedule信息。
func getSchedule(id int64) (*Schedule, error) { // {{{
	scd := &Schedule{}

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
	if err != nil {
		e := fmt.Sprintf("getSchedule run Sql error %s %s\n", sql, err.Error())
		return scd, errors.New(e)
	}
	g.L.Debugln("[s.getSchedule] ", "\nsql=", sql)

	scd.StartSecond = make([]time.Duration, 0)
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&scd.Id, &scd.Name, &scd.Count, &scd.Cyc,
			&scd.TimeOut, &scd.JobId, &scd.Desc)
		scd.setStart()
		if err != nil {
			e := fmt.Sprintf("getSchedule error %s\n", err.Error())
			return scd, errors.New(e)
		}

	}

	if scd.Id == -1 {
		e := fmt.Sprintf("not found schedule [%d] from db.\n", scd.Id)
		err = errors.New(e)
	}

	return scd, err
} // }}}

//从元数据库获取Schedule列表。
func getAllSchedules() ([]*Schedule, error) { // {{{
	scds := make([]*Schedule, 0)

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
	if err != nil {
		e := fmt.Sprintf("[getAllSchedule] run Sql error %s %s\n", sql, err.Error())
		return scds, errors.New(e)
	}
	g.L.Debugln("[getAllSchedule] ", "\nsql=", sql)

	for rows.Next() {
		scd := &Schedule{}
		scd.StartSecond = make([]time.Duration, 0)
		err = rows.Scan(&scd.Id, &scd.Name, &scd.Count, &scd.Cyc, &scd.TimeOut,
			&scd.JobId, &scd.Desc, &scd.CreateUserId, &scd.CreateTime, &scd.ModifyUserId,
			&scd.ModifyTime)
		scd.setStart()

		scds = append(scds, scd)
	}

	return scds, err
} // }}}

//从元数据库获取Job信息。
func getJob(id int64) (*Job, error) { // {{{
	j := &Job{}
	//查询全部Job列表
	sql := `SELECT job.job_id,
			   job.job_name,
			   job.job_desc,
			   job.prev_job_id,
			   job.next_job_id
			FROM scd_job job
			WHERE job.job_id=?`
	rows, err := g.HiveConn.Query(sql, id)
	if err != nil {
		e := fmt.Sprintf("[getJob] run Sql error %s %s\n", sql, err.Error())
		return nil, errors.New(e)
	}
	g.L.Debugln("[getJob] ", "\nsql=", sql)

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&j.Id, &j.Name, &j.Desc, &j.PreJobId, &j.NextJobId)
		//初始化Task内存
		j.Tasks = make(map[string]*Task)
	}

	if j.Id == 0 {
		j = nil
		e := fmt.Sprintf("[getJob] job [%d] not found \n", id)
		err = errors.New(e)
	}

	return j, err
} // }}}

//增加作业信息至元数据库
func (j *Job) add() (err error) { // {{{
	j.setNewId()
	j.Tasks = make(map[string]*Task)
	j.CreateTime, j.ModifyTime = time.Now(), time.Now()
	sql := `INSERT INTO scd_job
            (job_id, job_name, job_desc, prev_job_id,
             next_job_id, create_user_id, create_time,
             modify_user_id, modify_time)
		VALUES      (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = g.HiveConn.Exec(sql, &j.Id, &j.Name, &j.Desc, &j.PreJobId, &j.NextJobId, &j.CreateUserId, &j.CreateTime, &j.ModifyUserId, &j.ModifyTime)
	if err != nil {
		e := fmt.Sprintf("[j.add] run Sql error %s %s\n", sql, err.Error())
		return errors.New(e)
	}
	g.L.Debugln("[j.add] ", "\nsql=", sql)
	return err
} // }}}

//从元数据库获取Job下的Task列表。
func (j *Job) getTasksId() ([]int64, error) { // {{{
	tasksid := make([]int64, 0)

	//查询Job中全部Task列表
	sql := `SELECT jt.task_id
			FROM scd_job_task jt
            WHERE jt.job_id=?`
	rows, err := g.HiveConn.Query(sql, &j.Id)
	if err != nil {
		e := fmt.Sprintf("[j.getTasksId] Query sql [%s] error %s.\n", sql, err.Error())
		return tasksid, errors.New(e)
	}
	g.L.Debugln("[j.getTasksId] ", "\nsql=", sql)

	//循环读取记录
	for rows.Next() {
		var tid int64
		err = rows.Scan(&tid)
		tasksid = append(tasksid, tid)
	}
	return tasksid, err
} // }}}

//获取新Id
func (j *Job) setNewId() (err error) { // {{{
	var id int64

	//查询全部schedule列表
	sql := `SELECT max(job.job_id) as job_id
			FROM scd_job job`
	rows, err := g.HiveConn.Query(sql)
	if err != nil {
		e := fmt.Sprintf("[j.setNewId] Query sql [%s] error %s.\n", sql, err.Error())
		return errors.New(e)
	}

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
	}
	j.Id = id + 1

	return err

} // }}}

//修改作业信息至元数据库
func (j *Job) update() (err error) { // {{{
	sql := `UPDATE scd_job
		SET job_name=?, 
			job_desc=?,
			prev_job_id=?,
            next_job_id=?, 
            modify_user_id=?, 
			modify_time=?
	    WHERE job_id=?`
	_, err = g.HiveConn.Exec(sql, &j.Name, &j.Desc, &j.PreJobId, &j.NextJobId, &j.ModifyUserId, &j.ModifyTime, &j.Id)
	if err != nil {
		e := fmt.Sprintf("[j.update] Query sql [%s] error %s.\n", sql, err.Error())
		err = errors.New(e)
	}
	return err
} // }}}

//删除作业信息至元数据库
func (j *Job) deleteJob() (err error) { // {{{
	sql := `DELETE FROM scd_job WHERE job_id=?`
	_, err = g.HiveConn.Exec(sql, &j.Id)
	if err != nil {
		e := fmt.Sprintf("[j.setNewId] Query sql [%s] error %s.\n", sql, err.Error())
		err = errors.New(e)
	}
	return err
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
	if err != nil {
		e := fmt.Sprintf("[j.setNewId] Query sql [%s] error %s.\n", sql, err.Error())
		err = errors.New(e)
	}

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

//从元数据库获取Task信息。
func (t *Task) getTask() error { // {{{
	var td, id int64
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
	rows, err := g.HiveConn.Query(sql, t.Id)
	if err != nil {
		e := fmt.Sprintf("\n[t.getTask] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id, &t.Address, &t.Name, &t.TimeOut, &t.TaskType, &t.TaskCyc, &t.Desc, &td, &t.Cmd)
		if err != nil {
			e := fmt.Sprintf("\n[t.getTask] %s.", err.Error())
			return errors.New(e)
		}

		t.StartSecond = time.Duration(td) * time.Second
		//初始化relTask、param的内存
		t.RelTasksId = make([]int64, 0)
		t.RelTasks = make(map[string]*Task)
		t.Param = make([]string, 0)
		t.Attr = make(map[string]string)
	}

	if id == 0 {
		e := fmt.Sprintf("\n[t.getTask] task [%d] not found.", t.Id)
		err = errors.New(e)
	}

	return err
} // }}}

//从元数据库获取任务参数信息
func (t *Task) getTaskParam() error { // {{{
	//查询指定的Task属性列表
	sql := `SELECT pm.scd_param_name,
				   pm.scd_param_value
			FROM   scd_task_param pm
			WHERE pm.task_id=?`

	rows, err := g.HiveConn.Query(sql, t.Id)
	if err != nil {
		e := fmt.Sprintf("\n[t.getTaskParam] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		var name, value string
		err = rows.Scan(&name, &value)
		if err != nil {
			e := fmt.Sprintf("\n[t.getTaskParam] %s.", err.Error())
			return errors.New(e)
		}
		t.Param = append(t.Param, value)
	}
	return err
} // }}}

//从元数据库获取Job下的Task列表。
func (t *Task) getTaskAttr() error { // {{{

	//查询指定的Task属性列表
	sql := `SELECT ta.task_attr_name,
			   ta.task_attr_value
			FROM   scd_task_attr ta
			WHERE  task_id = ?`
	rows, err := g.HiveConn.Query(sql, t.Id)
	if err != nil {
		e := fmt.Sprintf("\n[t.getTaskAttr] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		var name, value string
		err = rows.Scan(&name, &value)
		if err != nil {
			e := fmt.Sprintf("\n[t.getTaskAttr] %s.", err.Error())
			return errors.New(e)
		}
		t.Attr[name] = value
	}
	return err
} // }}}

//从元数据库获取Task的依赖列表。
func (t *Task) getRelTaskId() error { // {{{
	//查询Task的依赖列表
	sql := `SELECT tr.rel_task_id
			FROM scd_task_rel tr
			Where tr.task_id=?`
	rows, err := g.HiveConn.Query(sql, t.Id)
	if err != nil {
		e := fmt.Sprintf("\n[t.getRelTaskId] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}

	//循环读取记录
	for rows.Next() {
		var rtid int64
		err = rows.Scan(&rtid)
		if err != nil {
			e := fmt.Sprintf("\n[t.getRelTaskId] %s.", err.Error())
			return errors.New(e)
		}
		t.RelTasksId = append(t.RelTasksId, rtid)
	}
	return err
} // }}}

//更新任务至元数据库
func (t *Task) update() error { // {{{
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
	_, err := g.HiveConn.Exec(sql, &t.Address, &t.Name, &t.TaskCyc, &t.TimeOut, &t.StartSecond, &t.TaskType, &t.Cmd, &t.Desc, &t.ModifyUserId, &t.ModifyTime, &t.Id)
	if err != nil {
		e := fmt.Sprintf("\n[t.update] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}
	return err
} // }}}

//DelParam方法从元数据库删除Task的Param信息
func (t *Task) delParam() error { // {{{
	sql := `DELETE FROM scd_task_param
			WHERE task_id=?`
	_, err := g.HiveConn.Exec(sql, &t.Id)
	if err != nil {
		e := fmt.Sprintf("\n[t.delParam] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}

	return err
} // }}}

//增加作业参数信息至元数据库
func (t *Task) addParam(pvalue string) error { // {{{
	pid, _ := t.getNewParamTaskId()
	sql := `INSERT INTO scd_task_param
            (scd_param_id,task_id, scd_param_name, scd_param_value,
             create_user_id, create_time)
			VALUES      (?, ?, ?, ?, ?, ?)`
	_, err := g.HiveConn.Exec(sql, &pid, &t.Id, "0", &pvalue, &t.CreateUserId, &t.CreateTime)
	if err != nil {
		e := fmt.Sprintf("\n[t.addParam] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}

	return err
} // }}}

//获取新TaskParamId
func (t *Task) getNewParamTaskId() (int64, error) { // {{{

	//查询全部schedule列表
	sql := `SELECT max(p.scd_param_id) as scd_param_id
			FROM scd_task_param p`

	rows, err := g.HiveConn.Query(sql)
	if err != nil {
		e := fmt.Sprintf("\n[t.getNewParamTaskId] sql %s error %s.", sql, err.Error())
		return -1, errors.New(e)
	}

	var id int64
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			e := fmt.Sprintf("[t.getNewParamTaskId] %s.\n", err.Error())
			return -1, errors.New(e)
		}

	}

	return id + 1, err
} // }}}

//获取新JobTaskId
func (t *Task) getNewRelTaskId() (int64, error) { // {{{

	//查询全部schedule列表
	sql := `SELECT max(rt.task_rel_id) as task_rel_id
			FROM scd_task_rel rt`

	rows, err := g.HiveConn.Query(sql)
	if err != nil {
		e := fmt.Sprintf("\n[t.getNewRelTaskId] sql %s error %s.", sql, err.Error())
		return -1, errors.New(e)
	}

	var id int64
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			e := fmt.Sprintf("[t.getNewRelTaskId] %s.\n", err.Error())
			return -1, errors.New(e)
		}
	}

	return id + 1, err
} // }}}

//获取新Id
func (t *Task) setNewId() (err error) { // {{{
	var id int64

	//查询全部schedule列表
	sql := `SELECT max(t.task_id) as task_id
			FROM scd_task t`
	rows, err := g.HiveConn.Query(sql)
	if err != nil {
		e := fmt.Sprintf("\n[t.setNewId] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}

	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			e := fmt.Sprintf("[t.setNewId] %s.\n", err.Error())
			return errors.New(e)
		}
	}
	t.Id = id + 1

	return err

} // }}}

//增加作业信息至元数据库
func (t *Task) add() (err error) { // {{{
	err = t.setNewId()
	if err != nil {
		e := fmt.Sprintf("[t.add] %s.\n", err.Error())
		return errors.New(e)
	}

	sql := `INSERT INTO scd_task
            (task_id, task_address, task_name, task_cyc,
             task_time_out, task_start, task_type_id,
             task_cmd, task_desc, create_user_id, create_time,
             modify_user_id, modify_time)
			VALUES      (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = g.HiveConn.Exec(sql, &t.Id, &t.Address, &t.Name, &t.TaskCyc, &t.TimeOut, &t.StartSecond, &t.TaskType, &t.Cmd, &t.Desc, &t.CreateUserId, &t.CreateTime, &t.ModifyUserId, &t.ModifyTime)
	if err != nil {
		e := fmt.Sprintf("\n[t.add] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}

	t.RelTasksId = make([]int64, 0)
	t.RelTasks = make(map[string]*Task)
	t.Attr = make(map[string]string)
	t.Param = make([]string, 0)
	return err
} // }}}

//增加依赖任务至元数据库
func (t *Task) addRelTask(id int64) error { // {{{
	relid, _ := t.getNewRelTaskId()
	sql := `INSERT INTO scd_task_rel
            (task_rel_id, task_id, rel_task_id, create_user_id, create_time)
			VALUES      (?, ?, ?, ?, ? )`
	_, err := g.HiveConn.Exec(sql, &relid, &t.Id, &id, &t.CreateUserId, &t.CreateTime)
	if err != nil {
		e := fmt.Sprintf("\n[t.addRelTask] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}

	return err
} // }}}

//GetRelJobId获取最大的Id
func (t *Task) getRelJobId() (int64, error) { // {{{

	//查询全部schedule列表
	sql := `SELECT max(t.job_task_id) as job_task_id
			FROM scd_job_task t`
	rows, err := g.HiveConn.Query(sql)
	if err != nil {
		e := fmt.Sprintf("\n[t.getRelJobId] sql %s error %s.", sql, err.Error())
		return -1, errors.New(e)
	}

	var id int64
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			e := fmt.Sprintf("[t.getRelJobId] %s.\n", err.Error())
			return -1, errors.New(e)
		}
	}

	return id + 1, err
} // }}}

//AddRelJob将Task与Job的关系持久化。
func (t *Task) addRelJob() (err error) { // {{{
	var id int64
	if id, err = t.getRelJobId(); err == nil {
		sql := `INSERT INTO scd_job_task
            (job_task_id,job_id,task_id,job_task_no,
            create_user_id,create_time)
            VALUES    (?, ?, ?, ?, ?, ?)`
		_, err = g.HiveConn.Exec(sql, &id, &t.JobId, &t.Id, &t.Id, &t.CreateUserId, &t.CreateTime)
	}
	return err
} // }}}

//删除依赖任务至元数据库
func (t *Task) deleteRelTask(id int64) error { // {{{
	sql := `DELETE FROM scd_task_rel WHERE task_id=? and rel_task_id=?`
	_, err := g.HiveConn.Exec(sql, &t.Id, &id)
	if err != nil {
		e := fmt.Sprintf("\n[t.deleteRelTask] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}

	return err
} // }}}

func (t *Task) deleteJobTaskRel() (err error) { // {{{
	sql := `DELETE FROM scd_job_task WHERE job_id=? and task_id=?`
	_, err = g.HiveConn.Exec(sql, &t.JobId, &t.Id)
	if err != nil {
		e := fmt.Sprintf("[t.deleteJobTaskRel] Query sql [%s] error %s.\n", sql, err.Error())
		err = errors.New(e)
	}

	return err
} // }}}

//删除任务至元数据库
func (t *Task) deleteTask() error { // {{{
	sql := `DELETE FROM scd_task WHERE task_id=?`
	_, err := g.HiveConn.Exec(sql, &t.Id)
	if err != nil {
		e := fmt.Sprintf("\n[t.deleteTask] sql %s error %s.", sql, err.Error())
		return errors.New(e)
	}

	return err
} // }}}
