package manager

import (
	"fmt"
	"github.com/go-martini/martini"
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

	m.Group("", func(r martini.Router) {
		r.Get("/schedules", getSchedules)
		r.Post("/schedules", addSchedule)
		r.Get("/schedules/:id", getScheduleById)
		r.Put("/schedules/:id", updateSchedule)
		r.Delete("/schedules/:id", deleteSchedule)
	})

	m.Post("/hello", func(ctx *web.Context) {
		fmt.Println(ctx.Params)
		ctx.WriteString("Hello World!")
	})

}

func getSchedules(ctx *web.Context, r render.Render, res http.ResponseWriter, Ss *schedule.ScheduleManager) {
	sl := make([]*schedule.Schedule, 0)
	for _, s := range Ss.ScheduleList {
		d := &schedule.Schedule{}
		schedule.Copy(d, s)
		d.Job = nil
		sl = append(sl, d)
	}
	r.JSON(200, sl)

}

//调度信息结构
type Schedule struct {
	Id           int64           //调度ID
	Name         string          //调度名称
	Count        int8            //调度次数
	Cyc          string          //调度周期
	StartSecond  []time.Duration //周期内启动时间
	StartMonth   []int           //周期内启动月份
	NextStart    time.Time       //下次启动时间
	TimeOut      int64           //最大执行时间
	JobId        int64           //作业ID
	Job          []Job           //作业
	Desc         string          //调度说明
	JobCnt       int64           //调度中作业数量
	TaskCnt      int64           //调度中任务数量
	CreateUserId int64           //创建人
	CreateTime   time.Time       //创人
	ModifyUserId int64           //修改人
	ModifyTime   time.Time       //修改时间
}

//作业信息结构
type Job struct {
	Id           int64     //作业ID
	ScheduleId   int64     //调度ID
	ScheduleCyc  string    //调度周期
	Name         string    //作业名称
	Desc         string    //作业说明
	PreJobId     int64     //上级作业ID
	NextJobId    int64     //下级作业ID
	Tasks        []Task    //作业中的任务
	TaskCnt      int64     //调度中任务数量
	CreateUserId int64     //创建人
	CreateTime   time.Time //创人
	ModifyUserId int64     //修改人
	ModifyTime   time.Time //修改时间
}

// 任务信息结构
type Task struct {
	Id           int64             // 任务的ID
	Address      string            // 任务的执行地址
	Name         string            // 任务名称
	JobType      string            // 任务类型
	ScheduleCyc  string            //调度周期
	TaskCyc      string            //调度周期
	StartSecond  time.Duration     //周期内启动时间
	Cmd          string            // 任务执行的命令或脚本、函数名等。
	Desc         string            //任务说明
	TimeOut      int64             // 设定超时时间，0表示不做超时限制。单位秒
	Param        map[string]string // 任务的参数信息
	Attr         map[string]string // 任务的属性信息
	JobId        int64             //所属作业ID
	RelTasks     []Task            //依赖的任务
	RelTaskCnt   int64             //依赖的任务数量
	CreateUserId int64             //创建人
	CreateTime   time.Time         //创人
	ModifyUserId int64             //修改人
	ModifyTime   time.Time         //修改时间
}

func getScheduleById(params martini.Params, r render.Render, res http.ResponseWriter, Ss *schedule.ScheduleManager) {

	if i, ok := params["id"]; ok {
		id, _ := strconv.Atoi(i)
		for _, s := range Ss.ScheduleList {
			if s.Id == int64(id) {
				r.JSON(200, getScheduleDetail(s))
				return
			}
		}
	}
}

func getScheduleDetail(s *schedule.Schedule) Schedule {
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

		jb.Tasks = make([]Task, 0)

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

			jb.Tasks = append(jb.Tasks, *tt)

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

}

func addSchedule(ctx *web.Context, res http.ResponseWriter, Ss *schedule.ScheduleManager) {
	fmt.Println(ctx.Params)
	fmt.Println(ctx.Request)

}

func deleteSchedule(ctx *web.Context, res http.ResponseWriter, Ss *schedule.ScheduleManager) {
	for _, s := range Ss.ScheduleList {
		res.Write([]byte(s.String()))
	}

}

func updateSchedule(params martini.Params, ctx *web.Context, res http.ResponseWriter) {
	fmt.Println(params)

}

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
