Spine = require('spineify')
Events  = Spine.Events
Module  = Spine.Module
Schedule = require('models/schedule')
Raphael = require('raphaelify')
Eve = require('eve')
Job = require('models/job')
$       = Spine.$

class JobManager extends Spine.Controller
  elements:
    ".close":  "close"
    ".jobpanel":  "jobpanel"
    "#jobname":  "jobname"
    "#jobdesc":  "jobdesc"
    "#prejobid":  "prejobid"

  events:
    "click .close": "postAddJob"
    "keypress #jobname": "keypress"
    "keypress #jobdesc": "keypress"

  constructor: (@paper, @color, @item, @width, @sinfo) -># {{{
    super
    Job.fetch({url:"/schedules/#{@item.Id}/jobs"}) if @item.Jobs

    @font = "Heiti, '黑体', 'Microsoft YaHei', '微软雅黑', SimSun, '宋体', '华文细黑', Helvetica, Tahoma, Arial, STXihei, sans-serif"
    @fontStyle = {fill: "#333", "font-family":@font, "text-anchor": "start", stroke: "none", "font-size": 18, "fill-opacity": 1, "stroke-width": 1}
    @jobFontStyle = {"font-family":@font , "font-size": 18, "stroke-opacity":1, "fill-opacity": 1, "stroke-width": 0}
    @jobcirStyle = {"fill-opacity": 0.2, "stroke-width": 1, cursor: "hand"}
    @jobrectStyle = {"fill-opacity": 0.1, "stroke-width": 0}
    @titlerectStyle = {fill: "#31708f", stroke: "#31708f", "fill-opacity": 0.05, "stroke-width": 0, cursor: "hand"}
    @addButtonStyle = {fill: "#31708f", stroke: "#31708f", "fill-opacity": 0.1, "stroke-opacity":0.2, "stroke-dasharray" : "-","stroke-width": 1, cursor: "hand"}

    @icoplus = "M25.979,12.896 19.312,12.896 19.312,6.229 12.647,6.229 12.647,12.896 5.979,12.896 5.979,19.562 12.647,19.562 12.647,26.229 19.312,26.229 19.312,19.562 25.979,19.562z"

    @height = 0

    @paper.setStart()

    [top,left] = [30, 10]
    @title = @paper.text(left, top, "作业列表").attr(@fontStyle)

    @titlerect = @paper.rect(left,top-20,190,35,3).attr(@titlerectStyle)
    @titlerect.hover(@hoveron,@hoverout)
    @titlerect.click(@showAddJob)

    @addButton = @paper.path(@icoplus)
    @addButton.attr(@addButtonStyle)
    @addButton.toBack()

    @set = @paper.setFinish()
    
    [top,left]=@refreshJobList(top+40,left)

    @height = top# }}}

  refreshJobList:(top,left) =># {{{
    @list = []
    return [top,left] unless @item.Jobs 
    for job,i in @item.Jobs when job.Id isnt 0
      s = @paper.set()
      s1 = @paper.set()
      jobname = @paper.text(left+80, top, job.Name).attr(@jobFontStyle)
      jobrect = @paper.rect(left,top-20,190,40,4).attr(@jobrectStyle)
      jobcir = @paper.circle(left+25,top,15).attr(@jobcirStyle)


      if job.TaskCnt is 0 and job.NextJobId is 0
        subButton = @paper.rect(left+150,top-5,25,8,4).attr(@jobrectStyle)
        subButton.attr(@jobrectStyle)
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

      [top,left]=[top+50,left]# }}}

  hlightTasks: -># {{{
    @data("sinfo")?.taskManager.hlight(@data("Id"))# }}}
      
  hoveron: -># {{{
    a = Raphael.animation({"fill-opacity": 0.8}, 200)
    @.animate(a)
      
  nlightTasks: -># {{{
    @data("sinfo")?.taskManager.nlight(@data("Id"))# }}}

  hoverout: -># {{{
    b = Raphael.animation({"fill-opacity": 0.1}, 200)
    @.animate(b)

  render: (x, y) =># {{{
    @html(require('views/schedule-add-job')(@lastJob))
    @el.css("display","block")
    @el.css("left",x)
    @el.css("top",y)
    @el.css("position","absolute")# }}}
    
  showAddJob: (e) =># {{{
    Spine.trigger("addJobRender",e.screenX,e.screenY)
    @jobname.focus()
    e# }}}

  keypress: (e) -># {{{
    if e.keyCode is 13 and e.ctrlKey
      @postAddJob(e)
    else if e.keyCode is 13
      @jobdesc.focus()# }}}

  delJob: (e) ->
    jb = Job.find(@.data("Id"))
    ts = @.data("this")
    jb.bind("change",ts?.delJobAndRefresh)
    jb.destroy({url:"/schedules/#{@.data("Sid")}/jobs/#{@.data("Id")}"})

  postAddJob: (e) -># {{{
    @el.css("display","none")
    jb = new Job()
    jb.bind("ajaxSuccess",@addJobAndRefresh)
    jb.Name = @jobname.val()
    jb.Desc = @jobdesc.val()
    jb.ScheduleId = @item.Id
    jb.PreJobId = if @prejobid.val() then parseInt(@prejobid.val()) else 0
    jb.Id = -1
    jb.create({url:"/schedules/#{@item.Id}/jobs"}) if jb.Name# }}}

  addJobAndRefresh: (data, status, xhr) =># {{{
    if xhr is "success"
      id = @item.Id
      Schedule.fetch({Id:id})
      @item = Schedule.find(id)
      for s in @list
        s.pop().remove()
        s.pop().remove()
        s.pop().remove()
      @refreshJobList(70, 10)
      @trigger("refreshJobList")# }}}
 
  delJobAndRefresh: (data, status, xhr) =># {{{
    id = @item.Id
    Schedule.fetch({Id:id})
    @item = Schedule.find(id)
    for s in @list
      s.pop().remove()
      s.pop().remove()
      s.pop().remove()
    @refreshJobList(70, 10)
    @trigger("refreshJobList")# }}}
 
module.exports = JobManager
