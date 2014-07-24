package manager

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/web"
	"github.com/rprp/hive/schedule"
	"log"
	"net/http"
	"strconv"
	"time"
)

func StartManager(sl *schedule.ScheduleManager) {
	m := martini.Classic()
	m.Use(Logger)
	m.Use(martini.Static("web/public"))
	m.Use(web.ContextWithCookieSecret(""))
	m.Use(render.Renderer(render.Options{
		Directory:       "templates",                // Specify what path to load the templates from.
		Extensions:      []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Delims:          render.Delims{"{[{", "}]}"},
		Charset:         "UTF-8",     // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON:      true,        // Output human readable JSON
		IndentXML:       true,        // Output human readable XML
		HTMLContentType: "text/html", // Output XHTML content type instead of default "text/html"
	}))

	m.Map(sl)
	controller(m)

	err := http.ListenAndServe(":3000", m)
	if err != nil {
		log.Fatal("Fail to start server: %v", err)
	}
}

func controller(m *martini.ClassicMartini) {
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Group("/schedules", func(r martini.Router) {
		r.Get("", getSchedules)
		r.Post("", addSchedule)
		r.Get("/:id", getScheduleById)
		r.Put("/:id", updateSchedule)
		r.Delete("/:id", deleteSchedule)

		r.Get("/:sid/jobs", getJobsForSchedule)
		r.Post("/:sid/jobs", binding.Bind(schedule.Job{}), addJob)
		r.Put("/:sid/jobs/:id", binding.Bind(schedule.Job{}), updateJob)
		r.Delete("/:sid/jobs/:id", deleteJob)

	})

}

func getSchedules(ctx *web.Context, r render.Render, res http.ResponseWriter, Ss *schedule.ScheduleManager) { // {{{
	sl := make([]*schedule.Schedule, 0)
	for _, s := range Ss.ScheduleList {
		d := &schedule.Schedule{}
		schedule.Copy(d, s)
		d.Job = nil
		d.Jobs = nil
		sl = append(sl, d)
	}
	r.JSON(200, sl)

} // }}}

//调度信息结构
type Schedule struct { // {{{
	Id           int64           `json: Id`           //调度ID
	Name         string          `json: Name`         //调度名称
	Count        int8            `json: Count`        //调度次数
	Cyc          string          `json: Cyc`          //调度周期
	StartSecond  []time.Duration `json: StartSecond`  //周期内启动时间
	StartMonth   []int           `json: StartMonth`   //周期内启动月份
	NextStart    time.Time       `json: NextStart`    //下次启动时间
	TimeOut      int64           `json: TimeOut`      //最大执行时间
	JobId        int64           `json: JobId`        //作业ID
	Job          []Job           `json: Job`          //作业
	Desc         string          `json: Desc`         //调度说明
	JobCnt       int64           `json: JobCnt`       //调度中作业数量
	TaskCnt      int64           `json: TaskCnt`      //调度中任务数量
	CreateUserId int64           `json: CreateUserId` //创建人
	CreateTime   time.Time       `json: CreateTime`   //创人
	ModifyUserId int64           `json: ModifyUserId` //修改人
	ModifyTime   time.Time       `json: ModifyTime`   //修改时间
}

//作业信息结构
type Job struct {
	Id           int64     `json: Id`           //作业ID
	ScheduleId   int64     `json: ScheduleId`   //调度ID
	ScheduleCyc  string    `json: ScheduleCyc`  //调度周期
	Name         string    `json: Name`         //作业名称
	Desc         string    `json: Desc`         //作业说明
	PreJobId     int64     `json: PreJobId`     //上级作业ID
	NextJobId    int64     `json: NextJobId`    //下级作业ID
	Tasks        []*Task   `json: Tasks`        //作业中的任务
	TaskCnt      int64     `json: TaskCnt`      //调度中任务数量
	CreateUserId int64     `json: CreateUserId` //创建人
	CreateTime   time.Time `json: CreateTime`   //创人
	ModifyUserId int64     `json: ModifyUserId` //修改人
	ModifyTime   time.Time `json: ModifyTime`   //修改时间
}

// 任务信息结构
type Task struct {
	Id           int64             `json: Id`           // 任务的ID
	Address      string            `json: Address`      // 任务的执行地址
	Name         string            `json: Name`         // 任务名称
	JobType      string            `json: JobType`      // 任务类型
	ScheduleCyc  string            `json: ScheduleCyc`  //调度周期
	TaskCyc      string            `json: TaskCyc`      //调度周期
	StartSecond  time.Duration     `json: StartSecond`  //周期内启动时间
	Cmd          string            `json: Cmd`          // 任务执行的命令或脚本、函数名等。
	Desc         string            `json: Desc`         //任务说明
	TimeOut      int64             `json: TimeOut`      // 设定超时时间，0表示不做超时限制。单位秒
	Param        map[string]string `json: Param`        // 任务的参数信息
	Attr         map[string]string `json: Attr`         // 任务的属性信息
	JobId        int64             `json: JobId`        //所属作业ID
	RelTasks     []Task            `json: RelTasks`     //依赖的任务
	RelTaskCnt   int64             `json: RelTaskCnt`   //依赖的任务数量
	CreateUserId int64             `json: CreateUserId` //创建人
	CreateTime   time.Time         `json: CreateTime`   //创人
	ModifyUserId int64             `json: ModifyUserId` //修改人
	ModifyTime   time.Time         `json: ModifyTime`   //修改时间
} // }}}

func getScheduleById(params martini.Params, r render.Render, res http.ResponseWriter, Ss *schedule.ScheduleManager) { // {{{

	if i, ok := params["id"]; ok {
		id, _ := strconv.Atoi(i)
		for _, s := range Ss.ScheduleList {
			if s.Id == int64(id) {
				r.JSON(200, s)
				return
			}
		}
	}
} // }}}

func getScheduleDetail(s *schedule.Schedule) Schedule { // {{{
	d := &Schedule{}
	d.Id = s.Id
	d.Name = s.Name
	d.Count = s.Count
	d.Cyc = s.Cyc
	d.StartSecond = s.StartSecond
	d.StartMonth = s.StartMonth
	d.NextStart = s.NextStart
	d.TimeOut = s.TimeOut
	d.JobId = s.JobId
	d.Job = make([]Job, 0)
	for j := s.Job; j != nil; {
		jb := &Job{}
		jb.Id = j.Id
		jb.ScheduleId = j.ScheduleId
		jb.ScheduleCyc = j.ScheduleCyc
		jb.Name = j.Name
		jb.Desc = j.Desc
		jb.PreJobId = j.PreJobId
		jb.NextJobId = j.NextJobId
		jb.TaskCnt = j.TaskCnt
		jb.CreateUserId = j.CreateUserId
		jb.CreateTime = j.CreateTime
		jb.ModifyUserId = j.ModifyUserId
		jb.ModifyTime = j.ModifyTime

		jb.Tasks = make([]*Task, 0)

		for _, t := range j.Tasks {
			tt := &Task{}
			tt.Id = t.Id
			tt.Address = t.Address
			tt.Name = t.Name
			tt.JobType = t.JobType
			tt.ScheduleCyc = t.ScheduleCyc
			tt.TaskCyc = t.TaskCyc
			tt.StartSecond = t.StartSecond
			tt.Cmd = t.Cmd
			tt.Desc = t.Desc
			tt.TimeOut = t.TimeOut
			tt.Param = t.Param
			tt.Attr = t.Attr
			tt.JobId = t.JobId

			tt.RelTasks = make([]Task, 0)
			for _, t1 := range t.RelTasks {
				t2 := &Task{}
				t2.Id = t1.Id
				t2.Address = t1.Address
				t2.Name = t1.Name
				tt.RelTasks = append(tt.RelTasks, *t2)
			}

			tt.RelTaskCnt = t.RelTaskCnt
			tt.CreateUserId = t.CreateUserId
			tt.CreateTime = t.CreateTime
			tt.ModifyUserId = t.ModifyUserId
			tt.ModifyTime = t.ModifyTime

			jb.Tasks = append(jb.Tasks, tt)

		}
		d.Job = append(d.Job, *jb)
		j = j.NextJob
	}

	d.Desc = s.Desc
	d.JobCnt = s.JobCnt
	d.TaskCnt = s.TaskCnt
	d.CreateUserId = s.CreateUserId
	d.CreateTime = s.CreateTime
	d.ModifyUserId = s.ModifyUserId
	d.ModifyTime = s.ModifyTime

	return *d

} // }}}

func addSchedule(ctx *web.Context, res http.ResponseWriter, Ss *schedule.ScheduleManager) { // {{{
	fmt.Println(ctx.Params)
	fmt.Println(ctx.Request)

} // }}}

func deleteJob(params martini.Params, ctx *web.Context, r render.Render, Ss *schedule.ScheduleManager) { // {{{

	sid, sidok := params["sid"]
	id, idok := params["id"]

	if !sidok || !idok {
		ctx.WriteHeader(500)
		return
	}

	ssid, _ := strconv.Atoi(sid)
	iid, _ := strconv.Atoi(id)

	if s := Ss.GetScheduleById(int64(ssid)); s != nil {
		if err := s.DeleteJob(int64(iid)); err != nil {
			ctx.WriteHeader(500)
			fmt.Println(err)
			fmt.Println("-----------------------------------------")
			ctx.WriteString("error:")
		} else {
			ctx.WriteHeader(204)
			ctx.WriteString("success")
			fmt.Println("-----------------ok------------------------")
		}

	}

} // }}}

//addJob获取客户端发送的Job信息，并调用Schedule的AddJob方法将其
//持久化并添加至Schedule中。
//成功返回添加好的Job信息
//错误返回err信息
func addJob(ctx *web.Context, r render.Render, Ss *schedule.ScheduleManager, job schedule.Job) { // {{{
	if job.Name == "" {
		ctx.WriteHeader(500)
		return
	}
	if s := Ss.GetScheduleById(int64(job.ScheduleId)); s != nil {
		job.ScheduleCyc = s.Cyc
		job.CreateUserId = 1
		job.ModifyUserId = 1
		if err := s.AddJob(&job); err != nil {
			ctx.WriteHeader(500)
			fmt.Println(err)
		} else {
			r.JSON(200, job)
		}
	} else {
		ctx.WriteHeader(500)
	}
} // }}}

//updateJob获取客户端发送的Job信息，并调用Schedule的UpdateJob方法将其
//持久化并更新至Schedule中。
//成功返回更新后的Job信息
func updateJob(ctx *web.Context, r render.Render, Ss *schedule.ScheduleManager, job schedule.Job) {
	if job.Name == "" {
		ctx.WriteHeader(500)
		return
	}
	if s := Ss.GetScheduleById(int64(job.ScheduleId)); s != nil {
		if err := s.UpdateJob(&job); err != nil {
			ctx.WriteHeader(500)
			fmt.Println(err)
		} else {
			r.JSON(200, job)
		}
	} else {
		ctx.WriteHeader(500)
	}

}

func getJobsForSchedule(ctx *web.Context, params martini.Params, r render.Render, res http.ResponseWriter, Ss *schedule.ScheduleManager) { // {{{

	sid, sidok := params["sid"]
	if !sidok {
		ctx.WriteHeader(500)
		return
	}

	ssid, _ := strconv.Atoi(sid)
	if s := Ss.GetScheduleById(int64(ssid)); s != nil {
		r.JSON(200, s.Jobs)
	}
	return
} // }}}

func deleteSchedule(ctx *web.Context, res http.ResponseWriter, Ss *schedule.ScheduleManager) { // {{{
	for _, s := range Ss.ScheduleList {
		res.Write([]byte(s.String()))
	}

} // }}}

func updateSchedule(params martini.Params, ctx *web.Context, res http.ResponseWriter) { // {{{
	fmt.Println(params)

} // }}}

func Logger() martini.Handler { // {{{
	return func(res http.ResponseWriter, req *http.Request, ctx martini.Context, log *log.Logger) {

		start := time.Now()
		log.Printf("Started %s %s", req.Method, req.URL.Path)

		rw := res.(martini.ResponseWriter)
		ctx.Next()

		content := fmt.Sprintf("Completed %v %s in %v", rw.Status(), http.StatusText(rw.Status()), time.Since(start))
		switch rw.Status() {
		case 200:
			content = fmt.Sprintf("\033[1;32m%s\033[0m", content)
		case 304:
			content = fmt.Sprintf("\033[1;33m%s\033[0m", content)
		case 404:
			content = fmt.Sprintf("\033[1;31m%s\033[0m", content)
		case 500:
			content = fmt.Sprintf("\033[1;36m%s\033[0m", content)
		}
		log.Println(content)
	}
} // }}}
