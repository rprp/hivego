Spine = require('spineify')
Events  = Spine.Events
Module  = Spine.Module
Style = require('controllers/style')
Schedule = require('models/schedule')
Raphael = require('raphaelify')
Eve = require('eve')
Job = require('models/job')
$       = Spine.$

class JobManager extends Spine.Controller
  elements:
    ".close":  "close"
    "#jobname":  "jobname"
    "#jobdesc":  "jobdesc"
    "#prejobid":  "prejobid"
    "#jobid":  "jobid"

  events:
    "click .close": "postJob"
    "keypress #jobname": "keypress"
    "keypress #jobdesc": "keypress"

  constructor: (@paper, @color, @item, @width, @sinfo) -># {{{
    super
    Job.fetch({url:"/schedules/#{@item.Id}/jobs"}) if @item.Jobs

    @isRefresh = true

    @height = 0

    @paper.setStart()

    [top,left] = [30, 10]
    @title = @paper.text(left, top, "作业列表").attr(Style.fontStyle)
    @titlerect = @paper.rect(left,top-20,190,35,3).attr(Style.titlerectStyle)
    @titlerect.hover(@hoveron,@hoverout)
    @titlerect.click(@showJob,@)

    @set = @paper.setFinish()
    
    [top,left]=@refreshJobList(top+40,left)

    @height = top# }}}

  refreshJobList:(top,left) =># {{{
    return [top,left] unless @item.Jobs
    return [top,left] unless @isRefresh
    if @list
      for s in @list
        s.pop().remove()
        s.pop().remove()
        s.pop().remove()
        s.pop()?.remove()

    @list = []
    for job,i in @item.Jobs when job.Id isnt 0
      s = @paper.set()
      s1 = @paper.set()
      jobname = @paper.text(left+80, top, job.Name).attr(Style.jobFontStyle)
      jobrect = @paper.rect(left,top-20,190,40,4).attr(Style.jobrectStyle)
      jobcir = @paper.circle(left+25,top,15).attr(Style.jobcirStyle)
      jobcir.click(@editJob,job)

      if job.TaskCnt is 0 and job.NextJobId is 0
        subButton = @paper.rect(left+150,top-4,30,7,4).attr(Style.jobrectStyle)
        subButton.attr(Style.jobrectStyle)
        subButton.attr("cursor","hand")
        subButton.attr("fill-opacity",0.01)
        subButton.data("Id",job.Id)
        subButton.data("Sid",@item.Id)
        subButton.data("item",@item)
        subButton.data("job",job)
        subButton.data("this",@)
        subButton.click(@delJob)
        s.push(subButton)
        s1.push(subButton)

      s.push(jobname, jobcir, jobrect)
      s.attr("stroke", @color[i])
      s.attr("fill", @color[i])
      s.data("Id",job.Id)
      s.data("sinfo",@sinfo)


      s1.push(jobcir)
      s.hover(@hoveron,@hoverout,s1,s1)
      s.hover(@hlightTasks,@nlightTasks)

      @set.push(s)
      @list.push(s)
      @lastJob = job

      @isRefresh = false
      [top,left]=[top+50,left]# }}}

  hlightTasks: -># {{{
    @data("sinfo")?.taskManager.hlight(@data("Id"))# }}}
      
  hoveron: -># {{{
    a = Raphael.animation({"fill-opacity": 0.8}, 200)
    @.animate(a)
    # }}}
      
  nlightTasks: -># {{{
    @data("sinfo")?.taskManager.nlight(@data("Id"))# }}}

  hoverout: -># {{{
    b = Raphael.animation({"fill-opacity": 0.1}, 200)
    @.animate(b)
    # }}}

  render: (x, y, job) =># {{{
    @html(require('views/schedule-add-job')(job))
    @el.css("display","block")
    @el.css("left",x-300)
    @el.css("top",y-120)
    @el.css("position","absolute")# }}}
    
  editJob: (e) -># {{{
    e = e||window.event
    @.opt = "edit"
    Spine.trigger("addJobRender", e.clientX, e.clientY, @)
    e# }}}

  showJob: (e) -># {{{
    e = e||window.event
    @.lastJob = new Job() unless @.lastJob
    @.lastJob.opt = "add"
    Spine.trigger("addJobRender",e.clientX,e.clientY,@.lastJob)
    @jobname.focus()
    e# }}}

  keypress: (e) -># {{{
    e = e||window.event
    if e.ctrlKey and e.keyCode in [13,10]
      @postJob(e)
    else if e.keyCode in [13,10]
      @jobdesc.focus()# }}}

  delJob: (e) ->
    jb = Job.find(@.data("Id"))
    ts = @.data("this")
    jb.bind("change",ts?.delJobAndRefresh)
    jb.destroy({url:"/schedules/#{@.data("Sid")}/jobs/#{@.data("Id")}"})

  postJob: (e) -># {{{
    @el.css("display","none")
    if @jobid.val()
      jb = Job.find(@jobid.val())
      jb.bind("ajaxSuccess",@addJobAndRefresh)
      jb.Name = @jobname.val()
      jb.Desc = @jobdesc.val()
      jb.save({url:"/schedules/#{@item.Id}/jobs/#{jb.Id}"})
    else
      jb = new Job()
      jb.bind("ajaxSuccess",@addJobAndRefresh)
      jb.ScheduleId = @item.Id
      jb.PreJobId = if @prejobid.val() then parseInt(@prejobid.val()) else 0
      jb.Id = -1
      jb.Name = @jobname.val()
      jb.Desc = @jobdesc.val()
      jb.create({url:"/schedules/#{@item.Id}/jobs"}) if jb.Name
  # }}}

  addJobAndRefresh: (data, status, xhr) =># {{{
    if xhr is "success"
      id = @item.Id
      Schedule.fetch({Id:id})
      @item = Schedule.find(id)
      @isRefresh = true
    # }}}
 
  delJobAndRefresh: (data, status, xhr) =># {{{
    id = @item.Id
    Schedule.fetch({Id:id})
    @item = Schedule.find(id)
    @isRefresh = true
    # }}}
 
module.exports = JobManager
