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

class MainInfo extends Spine.Controller
  className: 'maininfo'

  elements:
    ".pant": "pant"

  constructor: ->
    super
    Schedule.bind("findRecord",  @draw)

    Spine.bind("addJobRender", (x, y, job) => @append(@jobManager.render(x, y, job)))
    Spine.bind("addTaskRender", @addTaskRender = (task) =>
        unless @item.Jobs
          return
        @append (@taskForm.render(task))
      )

    @active @change

  change: (params) =>
    @paper = null
    Schedule.fetch({Id:params.id})
    @render()

  render: =>
    @html(require('views/main-info')())

  draw: (rs) =>
    if rs
      @item = Schedule.find(rs.Id)

    h = @item?.Jobs?.length*140
    @pant.css("height", 800 if h? or h < 800)
    [@width, @height] = [parseFloat(@pant.css("width")), parseFloat(@pant.css("height"))]

    unless @paper
      @paper = Raphael(@pant.get(0),'100%','100%')
      @color = Style.color
      @taskShape = new TaskManager.Shape(@paper,@color,@item,@width,@height)
      @taskShape.bind('refresh',@draw)
      @taskForm = new TaskManager.Form("c",@item)
      @taskForm.bind('updateTaskAndRefresh',@taskShape.updateTaskAndRefresh)
      @taskForm.bind('addTaskAndRefresh',@taskShape.addTaskAndRefresh)
      @taskForm.bind('refresh',@draw)
  
      slider = @paper.path("M #{@width-220},10L #{@width-220},#{@height}")
      slider.attr(Style.slider)
      
      @scheduleShape = new ScheduleManager.Shape(@paper,@color,@item,220)
      @scheduleForm = new ScheduleManager.Form("c",@item)
      @scheduleForm.bind("editScheduleRender", (x, y, schedule) => @append(@scheduleForm.render(x, y, schedule)))
  
      @jobManager = new JobManager(@paper,@color,@item,220,@)
      @jobManager.bind("rfJobList",@layout)
    else
      @scheduleShape.refreshSchedule(20,10)
      @jobManager.refreshJobList(70,10)

    @scheduleShape.titlerect.click(@scheduleForm.showSchedule,@scheduleForm)
    @layout()
    @append (@taskShape.el)
    @

  layout: =>
    @scheduleShape.st.transform("t#{@width-220},10")
    @jobManager.set.transform("t#{@width-220},#{@scheduleShape.height+10}")

module.exports = MainInfo
