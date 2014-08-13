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

func (j *Job) deleteTask(taskid int64) (err error) { // {{{
	sql := `DELETE FROM scd_job_task WHERE job_id=? and task_id=?`
	_, err = g.HiveConn.Exec(sql, &j.Id, &taskid)
	if err != nil {
		e := fmt.Sprintf("[j.deleteTask] Query sql [%s] error %s.\n", sql, err.Error())
		err = errors.New(e)
	}

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
