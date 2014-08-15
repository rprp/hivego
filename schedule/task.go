package schedule

import (
	"errors"
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

//根据Task.Id从元数据库获取信息初始化Task结构，包含以下动作
//初始化Task基本信息
//      Task属性信息
//      Task的参数信息
//      依赖的Task列表
//失败返回错误信息。
func (t *Task) InitTask() error { // {{{
	err := t.getTask()
	if err != nil {
		e := fmt.Sprintf("\n[t.InitTask] %s.", err.Error())
		return errors.New(e)
	}

	err = t.getTaskAttr()
	if err != nil {
		e := fmt.Sprintf("\n[t.InitTask] %s.", err.Error())
		return errors.New(e)
	}

	err = t.getTaskParam()
	if err != nil {
		e := fmt.Sprintf("\n[t.InitTask] %s.", err.Error())
		return errors.New(e)
	}

	t.RelTasksId = make([]int64, 0)
	t.RelTasks = make(map[string]*Task)
	t.RelTaskCnt = 0

	err = t.getRelTaskId()
	for _, rtid := range t.RelTasksId {
		ok := false
		t.RelTasks[string(rtid)], ok = g.Tasks[string(rtid)]
		if !ok {
			e := fmt.Sprintf("[t.InitTask] Task [%d] not found RelTask [%d] .\n", t.Id, rtid)
			g.L.Warningln(e)
			continue
		}
		t.RelTaskCnt++

	}

	return nil
} // }}}

//更新Task信息到元数据库。
//更新基本信息后，更新参数信息
func (t *Task) UpdateTask() error { // {{{
	err := t.update()
	if err != nil {
		e := fmt.Sprintf("\n[t.UpdateTask] %s.", err.Error())
		return errors.New(e)
	}

	err = t.delParam()
	if err != nil {
		e := fmt.Sprintf("\n[t.UpdateTask] %s.", err.Error())
		return errors.New(e)
	}

	for _, p := range t.Param {
		err = t.addParam(p)
		if err != nil {
			e := fmt.Sprintf("\n[t.UpdateTask] %s.", err.Error())
			return errors.New(e)
		}
	}

	return err
} // }}}

//AddTask方法持久化当前的Task信息。
//调用add方法将Task基本信息持久化。
//完成后处理作业关联信息、Task依赖关系、参数列表。
func (t *Task) AddTask() (err error) { // {{{
	err = t.add()
	if err != nil {
		e := fmt.Sprintf("\n[t.AddTask] %s.", err.Error())
		return errors.New(e)
	}

	err = t.addRelJob()
	if err != nil {
		e := fmt.Sprintf("\n[t.AddTask] %s.", err.Error())
		return errors.New(e)
	}

	for _, rt := range t.RelTasks {
		err = t.addRelTask(rt.Id)
		if err != nil {
			e := fmt.Sprintf("\n[t.AddTask] %s.", err.Error())
			return errors.New(e)
		}
	}

	for _, p := range t.Param {
		err = t.addParam(p)
		if err != nil {
			e := fmt.Sprintf("\n[t.AddTask] %s.", err.Error())
			return errors.New(e)
		}
	}

	return err
} // }}}

//删除依赖的任务关系
func (t *Task) DeleteRelTask(relid int64) error { // {{{
	var i int
	for k, v := range t.RelTasksId {
		if v == relid {
			i = k
		}
	}
	t.RelTasksId = append(t.RelTasksId[0:i], t.RelTasksId[i+1:]...)
	t.RelTaskCnt--
	delete(t.RelTasks, string(relid))

	err := t.deleteRelTask(relid)
	if err != nil {
		e := fmt.Sprintf("\n[t.DeleteRelTask] %s.", err.Error())
		return errors.New(e)
	}

	return err
} // }}}

//增加依赖的任务
func (t *Task) AddRelTask(rt *Task) (err error) { // {{{
	t.RelTasksId = append(t.RelTasksId, rt.Id)
	t.RelTaskCnt++
	t.RelTasks[string(rt.Id)] = rt

	err = t.addRelTask(rt.Id)
	if err != nil {
		e := fmt.Sprintf("\n[t.AddRelTask] error %s.", err.Error())
		return errors.New(e)
	}
	return err
} // }}}

//删除Task
func (t *Task) Delete() (err error) { // {{{
	err = t.delParam()
	if err != nil {
		e := fmt.Sprintf("\n[t.Delete] error %s.", err.Error())
		return errors.New(e)
	}

	for _, rid := range t.RelTasksId {
		err = t.DeleteRelTask(rid)
		if err != nil {
			e := fmt.Sprintf("\n[t.Delete] %s.", err.Error())
			return errors.New(e)
		}

	}

	err = t.deleteJobTaskRel()
	if err != nil {
		e := fmt.Sprintf("\n[t.Delete] error %s.", err.Error())
		return errors.New(e)
	}

	err = t.deleteTask()
	if err != nil {
		e := fmt.Sprintf("\n[t.Delete] error %s.", err.Error())
		return errors.New(e)
	}
	return err

} // }}}
