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

func StartManager(sl *schedule.ScheduleManager) { // {{{
	m := martini.Classic()
	m.Use(Logger)
	m.Use(martini.Static("web/public"))
	m.Use(web.ContextWithCookieSecret(""))
	m.Use(render.Renderer(render.Options{
		Directory:       "web/templates",            // Specify what path to load the templates from.
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
} // }}}

func controller(m *martini.ClassicMartini) { // {{{
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Group("/schedules", func(r martini.Router) {
		r.Get("", getSchedules)
		r.Post("", addSchedule)
		r.Get("/:id", getScheduleById)
		r.Put("/:id", binding.Bind(schedule.Schedule{}), updateSchedule)
		r.Delete("/:id", deleteSchedule)

		r.Get("/:sid/jobs", getJobsForSchedule)
		r.Post("/:sid/jobs", binding.Bind(schedule.Job{}), addJob)
		r.Put("/:sid/jobs/:id", binding.Bind(schedule.Job{}), updateJob)
		r.Delete("/:sid/jobs/:id", deleteJob)

		r.Get("/:sid/jobs/:jid/tasks", getJobsForSchedule)
		r.Post("/:sid/jobs/:jid/tasks", binding.Bind(schedule.Task{}), addTask)
		r.Put("/:sid/jobs/:jid/tasks/:id", binding.Bind(schedule.Task{}), updateTask)
		r.Delete("/:sid/jobs/:jid/tasks/:id", deleteTask)

		r.Post("/:sid/jobs/:jid/tasks/:id/reltask/:relid", addRelTask)
		r.Delete("/:sid/jobs/:jid/tasks/:id/reltask/:relid", deleteRelTask)
	})

} // }}}

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
		job.CreateTime = time.Now()
		job.ModifyTime = time.Now()
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
func updateJob(ctx *web.Context, r render.Render, Ss *schedule.ScheduleManager, job schedule.Job) { // {{{
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

} // }}}

//addTask获取客户端发送的Task信息，调用Task的AddTask方法持久化。
//成功后根据其中的JobId找到对应Job将其添加
//成功返回添加好的Job信息
//错误返回err信息
func addTask(params martini.Params, ctx *web.Context, r render.Render, Ss *schedule.ScheduleManager, task schedule.Task) { // {{{
	sid, sidok := params["sid"]
	ssid, _ := strconv.Atoi(sid)

	if !sidok || task.Name == "" || task.JobId == 0 {
		ctx.WriteHeader(500)
		return
	}

	task.TaskType = 1
	task.CreateUserId = 1
	task.ModifyUserId = 1
	task.CreateTime = time.Now()
	task.ModifyTime = time.Now()

	if err := task.AddTask(); err != nil {
		ctx.WriteHeader(500)
		fmt.Println(err)
	} else {
		if s := Ss.GetScheduleById(int64(ssid)); s != nil {
			if j := s.GetJobById(task.JobId); j != nil {
				j.Tasks[string(task.Id)] = &task
				j.TaskCnt++
			}
		}
		r.JSON(200, task)
	}
} // }}}

//deleteTask从调度结构中删除指定的Task，并持久化。
func deleteTask(params martini.Params, ctx *web.Context, r render.Render, Ss *schedule.ScheduleManager) { // {{{
	sid, _ := strconv.Atoi(params["sid"])
	jid, _ := strconv.Atoi(params["jid"])
	id, _ := strconv.Atoi(params["id"])

	if sid == 0 || jid == 0 || id == 0 {
		ctx.WriteHeader(500)
		return
	}

	if s := Ss.GetScheduleById(int64(sid)); s != nil {
		if err := s.DeleteTask(int64(id)); err != nil {
			ctx.WriteHeader(500)
			fmt.Println(err)
			return
		} else {
			r.JSON(200, nil)
		}
	}

} // }}}

//updateTask获取客户端发送的Task信息，并调用Job的UpdateTask方法将其
//持久化并更新至Job中。
//成功返回更新后的Task信息
func updateTask(params martini.Params, ctx *web.Context, r render.Render, Ss *schedule.ScheduleManager, task schedule.Task) { // {{{
	var err error
	sid, sidok := params["sid"]
	ssid, _ := strconv.Atoi(sid)

	if !sidok || task.Name == "" || task.JobId == 0 {
		ctx.WriteHeader(500)
		return
	}

	if s := Ss.GetScheduleById(int64(ssid)); s != nil {
		if j := s.GetJobById(task.JobId); j != nil {
			err = j.UpdateTask(&task)
		}
	}

	if err == nil {
		r.JSON(200, task)
	} else {
		ctx.WriteHeader(500)
		return
	}

} // }}}

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

//updateSchedule获取客户端发送的Schedule信息，并调用Schedule的Update方法将其
//持久化并更新至Schedule中。
//成功返回更新后的Schedule信息
func updateSchedule(params martini.Params, ctx *web.Context, r render.Render, Ss *schedule.ScheduleManager, scd schedule.Schedule) { // {{{
	if scd.Name == "" {
		ctx.WriteHeader(500)
		return
	}
	if s := Ss.GetScheduleById(int64(scd.Id)); s != nil {
		if err := s.UpdateSchedule(&scd); err != nil {
			ctx.WriteHeader(500)
			fmt.Println(err)
		} else {
			r.JSON(200, s)
		}
	} else {
		ctx.WriteHeader(500)
	}
} // }}}

//addRelTask根据Url参数获取到要添加的Task关系
func addRelTask(params martini.Params, ctx *web.Context, r render.Render, Ss *schedule.ScheduleManager) { // {{{
	sid, _ := strconv.Atoi(params["sid"])
	jid, _ := strconv.Atoi(params["jid"])
	id, _ := strconv.Atoi(params["id"])
	relid, _ := strconv.Atoi(params["relid"])

	if sid == 0 || jid == 0 || id == 0 || relid == 0 {
		ctx.WriteHeader(500)
		return
	}

	if s := Ss.GetScheduleById(int64(sid)); s != nil {
		t := s.GetTaskById(int64(id))
		rt := s.GetTaskById(int64(relid))

		if t == nil || rt == nil {
			ctx.WriteHeader(500)
			return
		}

		err := t.AddRelTask(rt)
		if err != nil {
			ctx.WriteHeader(500)
			return
		}
		r.JSON(200, t)
	}

} // }}}

func deleteRelTask(params martini.Params, ctx *web.Context, r render.Render, Ss *schedule.ScheduleManager) { // {{{
	sid, _ := strconv.Atoi(params["sid"])
	jid, _ := strconv.Atoi(params["jid"])
	id, _ := strconv.Atoi(params["id"])
	relid, _ := strconv.Atoi(params["relid"])

	if sid == 0 || jid == 0 || id == 0 || relid == 0 {
		ctx.WriteHeader(500)
		return
	}

	if s := Ss.GetScheduleById(int64(sid)); s != nil {
		t := s.GetTaskById(int64(id))

		if t == nil {
			ctx.WriteHeader(500)
			return
		}

		err := t.DeleteRelTask(int64(relid))
		if err != nil {
			ctx.WriteHeader(500)
			return
		}
		r.JSON(200, t)
	}

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
