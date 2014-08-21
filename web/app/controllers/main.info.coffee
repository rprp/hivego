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

  constructor: ->
    super
    Schedule.bind("findRecord",  @draw)

    Spine.bind("addJobRender", @renderJob = (x, y, job) =>
        @append(@ssl.jobManager.render(x, y, job))
      )
      
    Spine.bind("addTaskRender", @renderTask = (task) =>
        @append (@ssl.taskForm.render(task))
      )

    Spine.bind("editScheduleRender", @renderSchedule = (x, y, schedule) =>
        @append(@ssl.scheduleForm.render(x, y, schedule))
      )

    @active @change

  change: (params) =>
    Schedule.fetch({Id:params.id})
    @render()

  render: =>
    @html(require('views/main-info')())

  draw: (rs) =>
    if rs
      @item = Schedule.find(rs.Id)

    h = @item?.Jobs?.length*140
    h = 800 unless h
    h = 800 if h < 800
    @pant.css("height", h)

    [@width, @height] = [parseFloat(@pant.css("width")), parseFloat(@pant.css("height"))]

    if @ssl
      @ssl.scheduleShape.refreshSchedule(20,10)
      @ssl.jobManager.refreshJobList(70,10)
      @ssl.layout()
    else
      paper = Raphael(@pant.get(0),'100%','100%')
      @ssl = new ScheduleSymbol(paper,@width,@height,@item) 

    @append (@ssl.taskShape.el)
    @ssl

class ScheduleSymbol
  constructor: (@paper, @width, @height, @item) ->
    @color = Style.color
    [@st, @ed] = [Style.sopt,Style.eopt]
    @taskShape = new TaskManager.Shape(@paper,@color,@item,@width,@height)
    @taskForm = new TaskManager.Form("c",@item)
    @taskForm.bind('updateTaskAndRefresh',@taskShape.updateTaskAndRefresh)
    @taskForm.bind('addTaskAndRefresh',@taskShape.addTaskAndRefresh)

    slider = @paper.path("M #{@width-220},10L #{@width-220},#{@height}")
    slider.attr(Style.slider)
    
    @scheduleShape = new ScheduleManager.Shape(@paper,@color,@item,220)
    @scheduleForm = new ScheduleManager.Form("c",@item)
    @scheduleShape.titlerect.click(@scheduleForm.showSchedule,@scheduleForm)

    @jobManager = new JobManager(@paper,@color,@item,220,@)
    @jobManager.bind("rfJobList",@layout)
    @layout()

  layout: =>
    @scheduleShape.st.transform("t#{@width-220},10")
    @jobManager.set.transform("t#{@width-220},#{@scheduleShape.height+10}")

module.exports = ScheduleInfo
