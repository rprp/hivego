require('lib/setup')

Spine    = require('spineify')
Schedule = require('models/schedule')
MainList = require('controllers/main.list')
MainInfo = require('controllers/main.info')
Navbar = require('controllers/navbar')

class Main extends Spine.Stack
  className: 'smain container'

  controllers:
    mainList: MainList
    mainInfo: MainInfo

class App extends Spine.Controller
  events:
    "keypress": (e) ->
       e = e||window.event
       if e.keyCode is 47
         @nv.sinput.focus()

  constructor: ->
    super
    Schedule.fetch()

    Schedule.bind "ajaxError", (record, xhr, settings, error) ->
        console.log(error)
    Schedule.bind "ajaxSuccess", (data,status,xhr) ->
        console.log("ajaxSuccess #{xhr} #{status}")

    @nv = new Navbar
    @main = new Main
    @append @nv.render(), @main
    
    @nv.bind('addtask', @main.mainInfo.addTaskRender)
    @nv.bind('refreshAllTask', =>
        @main.mainInfo.draw()
        @main.mainInfo.taskShape.refreshTaskList()
        @main
      )

    @routes
      '': (params)->
          @main.mainList.active(params)
          @nv.show('list')
      '/schedules': (params) ->
          @main.mainList.active(params)
          @nv.show('list')
      '/schedules/:id': (params) =>
          @main.mainInfo.active(params)
          @nv.show('info')

    Spine.Route.setup()

module.exports = App
