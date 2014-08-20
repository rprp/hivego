require('lib/setup')

Spine    = require('spineify')
Manager = require('spineify/lib/manager')
Schedule = require('models/schedule')
ScheduleList = require('controllers/main.list')
ScheduleInfo = require('controllers/main.info')
Navbar = require('controllers/navbar')

class Main extends Spine.Stack
  className: 'smain'

  controllers:
    scheduleList: ScheduleList
    scheduleInfo: ScheduleInfo


class App extends Spine.Controller
  events:
    "keypress": "keypress"

  keypress: (e) ->
    e = e||window.event
    console.log(e.keyCode)

  constructor: ->
    super
    Spine.bind("showaddschedule", @showAddSchedule = (t) =>
        @main.scheduleInfo.active(t)
      )

    Schedule.fetch()

    Schedule.bind "ajaxError", (record, xhr, settings, error) ->
        console.log(error)
    Schedule.bind "ajaxSuccess", (data,status,xhr) ->
        console.log("ajaxSuccess #{xhr} #{status}")

    @nv = new Navbar
    @append @nv.render()

    @main = new Main
    @append @main
    
    @nv.bind('addtask',@main.scheduleInfo.renderTask)
    @routes
      '': (params)->
          @main.scheduleList.active(params)
          @nv.show('list')
      '/schedules': (params) ->
          @main.scheduleList.active(params)
          @nv.show('list')
      '/schedules/:id': (params) =>
          @main.scheduleInfo.active(params)
          @nv.show('info')

    Spine.Route.setup()



module.exports = App
