Spine = require('spine')
Schedule = require('models/schedule')
ScheduleItem = require('controllers/schedules.item')
$       = Spine.$

class ScheduleList extends Spine.Controller
  className: 'schedulelist'

  constructor: ->
    super
    Schedule.bind("create",  @addOne)
    Schedule.bind("refresh", @addAll)
    Schedule.fetch()

    Schedule.bind "ajaxError", (record, xhr, settings, error) -> 
        console.log(error)
    Schedule.bind "ajaxSuccess", (status,xhr) -> 
        console.log(xhr)

  addOne: (it) =>
    view = new ScheduleItem(item: it)
    @append(view.render())

  addAll: =>
    Schedule.each(@addOne)
    
module.exports = ScheduleList
