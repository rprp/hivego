Spine = require('spineify')
Raphael = require('raphaelify')
Eve = require('eve')
Schedule = require('models/schedule')
ScheduleManager = require('controllers/schedule.info')
JobManager = require('controllers/job.list')
TaskManager = require('controllers/task.list')
$       = Spine.$
wheel = require("jquery-mousewheel")($)

class ScheduleInfo extends Spine.Controller
  className: 'scheduleinfo'

  elements:
    ".pant":          "pant"

  #events:
    #"mousewheel .pant": "mousewheel"

  constructor: ->
    super
    Schedule.bind("findRecord",  @draw)
    Spine.bind("addJobRender", @renderJob)
    Spine.bind("editScheduleRender", @renderSchedule)
    @active @change

  change: (params) =>
    Schedule.fetch({Id:params.id})
    @render()

  render: =>
    @html(require('views/schedule-show-info')())

  mousewheel: (event, delta, deltaX, deltaY)->
    if delta > 0
      @ssl.taskManager.set.transform("...s1.1")
      tt.sp.refresh() for tt in @ssl.taskManager.ts
    else
      @ssl.taskManager.set.transform("...s0.9")
      tt.sp.refresh() for tt in @ssl.taskManager.ts
    event.stopPropagation()
    
  draw: (rs) =>
    @item = Schedule.find(rs.Id)

    h = @item?.Jobs?.length*140
    h = 800 unless h
    h = 800 if h < 800
    @pant.css("height", h)

    [@width, @height] = [parseFloat(@pant.css("width")), parseFloat(@pant.css("height"))]

    if @ssl
      @ssl.scheduleManager.refreshSchedule(20,10)
      @ssl.jobManager.refreshJobList(70,10)
      @ssl.layout()
    else
      paper = Raphael(@pant.get(0),'100%','100%')
      @ssl = new ScheduleSymbol(paper,@width,@height,@item) 
    @ssl

  renderSchedule: (x, y, schedule) =>
    @append(@ssl.scheduleManager.render(x, y, schedule))

  renderJob: (x, y, job) =>
    @append(@ssl.jobManager.render(x, y, job))

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

    @taskManager = new TaskManager(@paper,@color,@item,@width,@height)

    slider = @paper.path("M #{@width-220},10L #{@width-220},#{@height}")
    slider.attr({fill: "#333", "fill-opacity": 0.3, "stroke-width": 2, "stroke-opacity": 0.1})
    
    @scheduleManager = new ScheduleManager(@paper,@color,@item,220)
    @newJobManager()

    @layout()

  newJobManager: =>
    @jobManager = new JobManager(@paper,@color,@item,220,@)
    @jobManager.bind("rfJobList",@layout)
    @layout()

  layout: =>
    @scheduleManager.set.transform("t#{@width-220},10")
    @jobManager.set.transform("t#{@width-220},#{@scheduleManager.height+10}")

module.exports = ScheduleInfo
