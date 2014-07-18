Spine = require('spineify')
Schedule = require('models/schedule')
ScheduleItem = require('controllers/schedule.item')
$       = Spine.$

class ScheduleList extends Spine.Controller
  className: 'schedulelist'

  constructor: ->
    super
    Schedule.bind("create",  @addOne)
    Schedule.bind("refresh", @addAll)

  addOne: (it) =>
    view = new ScheduleItem(item: it)
    @append(view.render())

    view.pbmask.css("position","absolute")
    view.pbmask.css("z-index","1000")
    view.pbmask.css("width",view.sstart.css("width"))
    view.pbmask.css("height",view.body.css("height"))
    view.pbmask.css("background-color","#f5f5f5")
    view.sname.css("color","#f5f5f5")

  addAll: =>
    Schedule.each(@addOne)
    
module.exports = ScheduleList
