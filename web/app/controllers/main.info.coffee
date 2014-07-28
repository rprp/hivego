Spine = require('spineify')
Raphael = require('raphaelify')
Eve = require('eve')
Schedule = require('models/schedule')
Style = require('controllers/style')
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
    @color = Style.color
    [@st, @ed] = [Style.sopt,Style.eopt]
    @taskManager = new TaskManager(@paper,@color,@item,@width,@height)
    slider = @paper.path("M #{@width-220},10L #{@width-220},#{@height}")
    slider.attr(Style.slider)
    
    @scheduleManager = new ScheduleManager(@paper,@color,@item,220)
    @newJobManager()
    @layout()

  newJobManager: =>
    @jobManager = new JobManager(@paper,@color,@item,220,@)
    @jobManager.bind("rfJobList",@layout)
    @layout()

  layout: =>
    @scheduleManager.st.transform("t#{@width-220},10")
    @jobManager.set.transform("t#{@width-220},#{@scheduleManager.height+10}")

module.exports = ScheduleInfo
