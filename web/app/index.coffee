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
    Spine.bind("showaddschedule", @showAddSchedule)
    Schedule.fetch()

    Schedule.bind "ajaxError", (record, xhr, settings, error) ->
        console.log(error)
    Schedule.bind "ajaxSuccess", (data,status,xhr) ->
        console.log("ajaxSuccess#{data}    #{xhr}   #{status}")


    nv = new Navbar
    @append nv.render()

    @main = new Main
    @append @main
    
    @routes
      '': (params)-> @main.scheduleList.active(params)
      '/schedules': (params)-> @main.scheduleList.active(params)
      '/schedules/:id': (params) -> @main.scheduleInfo.active(params)

    Spine.Route.setup()

  showAddSchedule: (t) =>
    @main.scheduleInfo.active(t)


module.exports = App
