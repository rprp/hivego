require('lib/setup')

Spine    = require('spine')
Manager = require('spine/lib/manager')
Schedule = require('models/schedule')
ScheduleList = require('controllers/schedules')
ScheduleInfo = require('controllers/schedulesinfo')

class Main extends Spine.Stack
  className: 'smain'

  controllers:
    scheduleList: ScheduleList
    scheduleInfo: ScheduleInfo

class App extends Spine.Controller

  constructor: ->
    super
    Schedule.fetch()

    Schedule.bind "ajaxError", (record, xhr, settings, error) -> 
        console.log(error)
    Schedule.bind "ajaxSuccess", (status,xhr) -> 
        console.log(xhr)

    main = new Main
    @append main
    
    @routes
      '': (params)-> main.scheduleList.active(params)
      '/schedules': (params)-> main.scheduleList.active(params)
      '/schedules/:id': (params) -> main.scheduleInfo.active(params)
      #'/contacts/:id':    (params) -> @show.active(params)
      #'/contacts':        (params) -> @list.active(params)

    Spine.Route.setup()

module.exports = App