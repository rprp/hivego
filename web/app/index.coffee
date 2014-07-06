require('lib/setup')

Spine    = require('spine')
Manager = require('spine/lib/manager')
ScheduleList = require('controllers/schedules')

class Main extends Spine.Stack
  controllers:
    schedules: ScheduleList

class App extends Spine.Controller

  constructor: ->
    super

    main = new Main
    @append main.schedules.active()

    Spine.Route.setup()

module.exports = App
