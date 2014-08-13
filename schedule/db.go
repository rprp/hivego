package schedule

import (
	"errors"
	"fmt"
	"time"
)

//Add方法会将Schedule对象增加到元数据库中。
func (s *Schedule) add() (err error) { // {{{
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
func (s *Schedule) update() (err error) { // {{{
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
func (s *Schedule) deleteSchedule() error { // {{{
	sql := `Delete FROM scd_schedule WHERE scd_id=?`
	_, err := g.HiveConn.Exec(sql, &s.Id)
	g.L.Debugln("schedule", s.Name, " was deleted.")

	return err
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

//addStart将Schedule的启动列表持久化到数据库
//添加前先调用delStart方法将Schedule中的原有启动列表清空
//需要注意的是：内存中的启动列表单位为纳秒，存储前需要转成秒
//若成功则开始添加，失败返回err信息
func (s *Schedule) addStart() (err error) {
	if err = s.delStart(); err == nil {
		for i, st := range s.StartSecond {
			t := time.Duration(st) / time.Second
			sm := s.StartMonth[i]
			sql := `INSERT INTO scd_start 
            (scd_id, scd_start, scd_start_month,
            create_user_id, create_time)
         VALUES  (?, ?, ?, ?, ?)`
			_, err = g.HiveConn.Exec(sql, &s.Id, &t, &sm, &s.ModifyUserId, &s.ModifyTime)
			CheckErr("setStart run Sql "+sql, err)
		}
	}
	return err
}

//delStart删除该Schedule的所有启动时间列表
func (s *Schedule) delStart() (err error) {
	sql := `DELETE FROM scd_start WHERE scd_id=?`
	_, err = g.HiveConn.Exec(sql, &s.Id)
	return err

}

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
		g.L.Warningln(e)
		return scd, errors.New(e)
	}

	scd.StartSecond = make([]time.Duration, 0)
	//循环读取记录，格式化后存入变量ｂ
	for rows.Next() {
		err = rows.Scan(&scd.Id, &scd.Name, &scd.Count, &scd.Cyc,
			&scd.TimeOut, &scd.JobId, &scd.Desc)
		scd.setStart()
		if err != nil {
			e := fmt.Sprintf("getSchedule error %s\n", err.Error())
			g.L.Warningln(e)
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
func getAllSchedules() (scds []*Schedule, err error) { // {{{
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

	return scds, nil
} // }}}
