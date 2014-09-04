require('lib/setup')

Spine    = require('spineify')
Schedule = require('models/schedule')
MainList = require('controllers/main.list')
MainInfo = require('controllers/main.info')
Navbar = require('controllers/navbar')
Msg = require('controllers/msg')

class Main extends Spine.Stack
  className: 'smain'

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

    Schedule.bind "ajaxError", (xhr, st, error) ->
        stxt = "#{st.status} #{st.statusText} #{st.responseText}"
        Spine.trigger("msg",st.status,stxt)

    @nv = new Navbar
    @msg = new Msg
    @main = new Main
    @append @nv.render(), @main, @msg.render()
    
    @nv.bind('addtask', @main.mainInfo.addTaskRender)
    @nv.bind('refreshAllTask', =>
        Schedule.fetch({Id: @main.mainInfo.item.Id})
        Schedule.bind("findRecord", (rs) =>
                @main.mainInfo.paper
                @main.mainInfo.draw(rs)
                @main.mainInfo.taskShape.refreshTaskList()
          )
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
