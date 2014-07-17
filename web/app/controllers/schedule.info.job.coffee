Spine = require('spineify')
Events  = Spine.Events
Module  = Spine.Module
Raphael = require('raphaelify')
Eve = require('eve')
ScheduleInfo = require('controllers/schedule.info')
Schedule = require('models/schedule')
$       = Spine.$

class Job extends Spine.Controller
  constructor: (@paper, @color, @item, @width, @sinfo) ->
    super
    @font = "Helvetica, Tahoma, Arial, STXihei, '华文细黑', Heiti, '黑体', 'Microsoft YaHei', '微软雅黑', SimSun, '宋体', sans-serif"
    @fontStyle = {fill: "#333", "font-family":@font, "text-anchor": "start", stroke: "none", "font-size": 18, "fill-opacity": 1, "stroke-width": 1}
    @jobFontStyle = {stroke: @color[i], "fill": @color[i], "font-family":@font , "font-size": 16, "stroke-opacity":1, "fill-opacity": 1, "stroke-width": 1}
    @height = 0

    @paper.setStart()

    [top,left] = [0, 10]
    [top,left]=[top+30,left]
    @title = @paper.text(left, top, "任务列表").attr(@fontStyle)

    [top,left]=[top+40,left]
    @list = []
    for job,i in @item.Job when job.Id isnt 0
      s = @paper.set()
      jobname = @paper.text(left+80, top, job.Name).attr(@jobFontStyle)
      jobname.attr("stroke", @color[i])
      jobname.attr("fill", @color[i])

      jobcir = @paper.circle(left+25,top,15)
      jobcir.attr({fill: @color[i], stroke: @color[i], "fill-opacity": 0.2, "stroke-width": 1})

      jobrect = @paper.rect(left,top-20,180,40,5)
      jobrect.attr({fill: @color[i], stroke: @color[i], "fill-opacity": 0.1, "stroke-width": 0})

      s.push(jobname, jobcir, jobrect)
      s.data("Id",job.Id)
      s.data("sinfo",@sinfo)
      s.hover(@hoveron,@hoverout)
      @list.push(s)
      [top,left]=[top+50,left]

    @addButton = @paper.path("M25.979,12.896 19.312,12.896 19.312,6.229 12.647,6.229 12.647,12.896 5.979,12.896 5.979,19.562 12.647,19.562 12.647,26.229 19.312,26.229 19.312,19.562 25.979,19.562z")
    @addButton.transform("t#{left+60},#{top-10}s2.5")
    @addButton.attr({fill: "#00FF00", stroke: "#00FF00", "fill-opacity": 0.2, "stroke-opacity":0.8, "stroke-dasharray" : "-","stroke-width": 1, cursor: "hand"})

    @set = @paper.setFinish()
    @height = top 

  hoveron: ->
    a = Raphael.animation({"fill-opacity": 0.6}, 300)
    @.animate(a)
    @data("sinfo").hlight(@data("Id"))
      
  hoverout: ->
    b = Raphael.animation({"fill-opacity": 0.1}, 300)
    @.animate(b)
    @data("sinfo").nlight(@data("Id"))
 
module.exports = Job
