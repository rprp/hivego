Spine = require('spineify')
Events  = Spine.Events
Module  = Spine.Module
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
    "click .close": "hideAddJob"
    "keypress #jobname": "keypress"
    "keypress #jobdesc": "keypress"

  constructor: (@paper, @color, @item, @width, @sinfo) ->
    super
    @font = "Helvetica, Tahoma, Arial, STXihei, '华文细黑', Heiti, '黑体', 'Microsoft YaHei', '微软雅黑', SimSun, '宋体', sans-serif"
    @fontStyle = {fill: "#333", "font-family":@font, "text-anchor": "start", stroke: "none", "font-size": 18, "fill-opacity": 1, "stroke-width": 1}
    @jobFontStyle = {stroke: @color[i], "fill": @color[i], "font-family":@font , "font-size": 16, "stroke-opacity":1, "fill-opacity": 1, "stroke-width": 1}
    icoplus = "M25.979,12.896 19.312,12.896 19.312,6.229 12.647,6.229 12.647,12.896 5.979,12.896 5.979,19.562 12.647,19.562 12.647,26.229 19.312,26.229 19.312,19.562 25.979,19.562z"
    @height = 0

    @paper.setStart()

    [top,left] = [30, 10]
    @title = @paper.text(left, top, "作业列表").attr(@fontStyle)

    [top,left]=[top+40,left]
    @list = []
    for job,i in @item.Job when job.Id isnt 0
      s = @paper.set()
      jobname = @paper.text(left+80, top, job.Name).attr(@jobFontStyle)
      jobname.attr("stroke", @color[i])
      jobname.attr("fill", @color[i])

      jobcir = @paper.circle(left+25,top,15)
      jobcir.attr({fill: @color[i], stroke: @color[i], "fill-opacity": 0.2, "stroke-width": 1})

      jobrect = @paper.rect(left,top-20,190,40,5)
      jobrect.attr({fill: @color[i], stroke: @color[i], "fill-opacity": 0.1, "stroke-width": 0})

      s.push(jobname, jobcir, jobrect)
      s.data("Id",job.Id)
      s.data("sinfo",@sinfo)
      s.hover(@hoveron,@hoverout)
      @list.push(s)
      @lastJob = job
      [top,left]=[top+50,left]

    addbtn = @paper.rect(left,top-20,190,40,5).attr({fill: "#31708f", stroke: "#31708f", "fill-opacity": 0.1, "stroke-width": 0, cursor: "hand"})
    addbtn.hover(@hoveron,@hoverout)
    addbtn.click(@addjob)
    @addButton = @paper.path(icoplus)
    @addButton.attr({fill: "#31708f", stroke: "#31708f", "fill-opacity": 0.4, "stroke-opacity":0.8, "stroke-dasharray" : "-","stroke-width": 1, cursor: "hand"})
    @addButton.toBack()

    @set = @paper.setFinish()
    @height = top

  hoveron: ->
    a = Raphael.animation({"fill-opacity": 0.6}, 300)
    @.animate(a)
    @data("sinfo")?.taskManager.hlight(@data("Id"))
      
  hoverout: ->
    b = Raphael.animation({"fill-opacity": 0.1}, 300)
    @.animate(b)
    @data("sinfo")?.taskManager.nlight(@data("Id"))

  render: (x, y) =>
    @html(require('views/schedule-add-job')(@lastJob))
    @el.css("display","block")
    @el.css("left",x)
    @el.css("top",y)
    @el.css("position","absolute")
    
  addjob: (e) =>
    Spine.trigger("addjob",e.screenX,e.screenY)
    @jobname.focus()
    e

  keypress: (e) ->
    if e.keyCode is 13 and e.ctrlKey
      @hideAddJob(e)
    else if e.keyCode is 13
      @jobdesc.focus()

  hideAddJob: (e) ->
    @el.css("display","none")
    jb = new Job()
    jb.Name = @jobname.val()
    jb.Desc = @jobdesc.val()
    jb.ScheduleId = @item.Id
    jb.PreJobId = if @prejobid.val() then parseInt(@prejobid.val()) else 0
    jb.Id = -1
    jb.save() if jb.Name

 
module.exports = JobManager
