Spine = require('spine')
Raphael = require('raphael')
Eve = require('eve')
Schedule = require('models/schedule')
ScheduleDetail = require('controllers/schedule.info.detail')
Job = require('controllers/schedule.info.job')
Task = require('controllers/schedule.info.task')
$       = Spine.$
wheel = require("jquery-mousewheel")($)

class ScheduleInfo extends Spine.Controller
  className: 'scheduleinfo'


  elements:
    ".pant":          "pant"

  events:
    "mousewheel .pant": "mousewheel"

  constructor: ->
    super
    Schedule.bind("findRecord",  @draw)
    @active @change

  change: (params) =>
    Schedule.fetch({Id:params.id})
    @render()

  render: =>
    @html(require('views/schedule-show-info')())


  mousewheel: (event, delta, deltaX, deltaY)->
    if delta > 0
      @ssl.job.set.transform("...s1.1")
    else
      @ssl.job.set.transform("...s0.9")

    event.stopPropagation()
    
  draw: (rs) =>
    @item = Schedule.find(rs.Id)

    paper = Raphael(@pant.get(0),'100%','100%')

    [@width, @height] = [parseFloat(@pant.css("width")), parseFloat(@pant.css("height"))]

    @ssl = new ScheduleSymbol(paper,@width,@height,@item)

class ScheduleSymbol
  constructor: (@paper, @width, @height, @item) ->

    @color=['#FF8C00', '#008000', '#2F4F4F', '#DA70D6', '#0000FF', '#8A2BE2', '#6495ED', '#B8860B', '#FF4500', '#AFEEEE', '#DB7093',
        '#CD853F', '#FFC0CB', '#B0E0E6', '#BC8F8F', '#4169E1', '#8B4513', '#00FFFF', '#00BFFF', '#008B8B',
        '#ADFF2F', '#4B0082', '#F0E68C', '#7CFC00', '#7FFF00', '#DEB887', '#98FB98', '#FFD700', '#5F9EA0', '#D2691E', '#A9A9A9',
        '#8B008B', '#556B2F', '#9932CC', '#8FBC8B', '#483D8B', '#00CED1', '#9400D3', '#FF69B4', '#228B22', '#1E90FF', '#FF00FF',
        '#FFB6C1', '#FFA07A', '#20B2AA', '#87CEFA', '#00FF00', '#B0C4DE', '#FF00FF', '#32CD32', '#0000CD', '#66CDAA', '#BA55D3',
        '#9370DB', '#3CB371', '#7B68EE', '#00FA9A', '#48D1CC', '#C71585', '#191970', '#000080', '#808000', '#6B8E23', '#FFA500',
        '#F4A460', '#2E8B57', '#A0522D', '#87CEEB', '#6A5ACD', '#708090', '#00FF7F', '#4682B4', '#D2B48C', '#008080', '#40E0D0',
         '#006400', '#BDB76B','#EE82EE', '#F5DEB3', '#FFFF00', '#9ACD32']

    #[@st, @ed] = [Raphael.animation({"fill-opacity": .2}, 2000, -> @.animate(ed)),
                #Raphael.animation({"fill-opacity": .5}, 2000, -> @.animate(st))]
    [@st, @ed] = [Raphael.animation({"fill-opacity": .2}, 1000),
                Raphael.animation({"fill-opacity": .5}, 1000)]

    top = 80
    @ts = []
    for job,i in @item.Job

      spacing = (@width-200)/job.Tasks.length if job.Tasks.length > 0
      r = 25
      spacing = 100

      left = (@width-200)/2-(job.Tasks.length/2) * spacing if job.Tasks.length > 0
      for task,j in job.Tasks
        t= new Task(paper,left,top,task.Name,@color[i],r)
        t.Id = task.Id
        t.JobId = job.Id
        t.RelTaskId = (rt.Id for rt in task.RelTasks)
        @ts.push(t)

        rts.addNext(t) for rts in @getTaskSymbol(t.RelTaskId)
        left += spacing

      top += 120

    slider = @paper.path("M #{@width-220},10L #{@width-220},#{@height}")
    slider.attr({fill: "#333", "fill-opacity": 0.3, "stroke-width": 2, "stroke-opacity": 0.1})
    
    @scheduleDetail = new ScheduleDetail(@paper,@color,@item,220)
    @job = new Job(@paper,@color,@item,220,@)

    @layout()

    #jobflg=@paper.circle(@width-60, 100, 15)
    #jobflg.attr({fill: "green", stroke: "green", "fill-opacity": 0.5, "stroke-width": 1, cursor: "hand"})
    #jobflg.animate(st)
    #jobflg.hover (-> @.attr({r: 17})), (-> @.attr({r: 15}))

    #jobflg.click =>

  getTaskSymbol: (Ids) ->
     t for t in @ts when t.Id in Ids

  hlight: (Id) ->
    a = Raphael.animation({"fill-opacity": 0.7}, 500)
    for t in @ts
      if t.JobId is Id
        t.sp.animate(a)
        t.sp.g = t.sp.glow({color: t.sp.attr("fill"), "fill-opacity": 0.2, width:10})

  nlight: (Id) ->
    a = Raphael.animation({"fill-opacity": 0.2}, 500)
    for t in @ts
      if t.JobId is Id
        t.sp.animate(a)
        t.sp.g.remove()

  layout: ->
    @scheduleDetail.set.transform("t#{@width-220},10")
    @job.set.transform("t#{@width-220},#{@scheduleDetail.height+10}")
    @job.addButton.transform("t1020,#{@scheduleDetail.height + @job.height}s2.5")

module.exports = ScheduleInfo
