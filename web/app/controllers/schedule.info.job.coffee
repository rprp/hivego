Spine = require('spineify')
Events  = Spine.Events
Module  = Spine.Module
Raphael = require('raphaelify')
Eve = require('eve')
ScheduleInfo = require('controllers/schedule.info')
Schedule = require('models/schedule')
$       = Spine.$

class Job extends Spine.Controller
  elements:
    ".close":  "close"

  events:
    "click .close": "cl"

  constructor: (@paper, @color, @item, @width, @sinfo) ->
    super
    @font = "Helvetica, Tahoma, Arial, STXihei, '华文细黑', Heiti, '黑体', 'Microsoft YaHei', '微软雅黑', SimSun, '宋体', sans-serif"
    @fontStyle = {fill: "#333", "font-family":@font, "text-anchor": "start", stroke: "none", "font-size": 18, "fill-opacity": 1, "stroke-width": 1}
    @jobFontStyle = {stroke: @color[i], "fill": @color[i], "font-family":@font , "font-size": 16, "stroke-opacity":1, "fill-opacity": 1, "stroke-width": 1}
    icoplus = "M25.979,12.896 19.312,12.896 19.312,6.229 12.647,6.229 12.647,12.896 5.979,12.896 5.979,19.562 12.647,19.562 12.647,26.229 19.312,26.229 19.312,19.562 25.979,19.562z"
    @height = 0

    @paper.setStart()

    [top,left] = [0, 10]
    [top,left]=[top+30,left]
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
    @data("sinfo")?.task.hlight(@data("Id"))
      
  hoverout: ->
    b = Raphael.animation({"fill-opacity": 0.1}, 300)
    @.animate(b)
    @data("sinfo")?.task.nlight(@data("Id"))

  render: (x, y) =>
    @html(require('views/schedule-add-job')())
    @el.css("left",x)
    @el.css("top",y)
    @el.css("position","absolute")
    @el.css("display","block")
    
  addjob: (e) =>
    Spine.trigger("addjob",e.screenX,e.screenY)

  cl: (e) ->
    @el.css("display","none")

 
module.exports = Job
