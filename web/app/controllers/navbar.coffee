Spine = require('spineify')
Schedule = require('models/schedule')
ScheduleManager = require('controllers/schedule.info')
$       = Spine.$

class Navbar extends Spine.Controller
  className: 'hnavbar'
  
  elements:
    "#addSchedule": "addSchedule"
    "#addTask": "addTask"
    "#refreshAll": "refreshAll"
    ".sinput": "sinput"

  events:
    "mouseenter h1": (e) -> $(e.target).stop().animate({backgroundColor:'#777'},400)
    "mouseleave h1": (e) -> $(e.target).stop().animate({backgroundColor:'#333'},200)

    "click #addSchedule": (e) ->
        if @sf
          @sf.el.remove()

        e = e||window.event
        s = new Schedule({Id:-1})
        @sf = new ScheduleManager.Form("c",s)
        @append(@sf.render(550,100,s))
        @sf.el.css("z-index",1000)

    "click #addTask": (e) ->
        e = e||window.event
        @trigger('addtask',e)

    "click #refreshAll": (e) ->
        e = e||window.event
        @trigger('refreshAllTask',e)

    "click #home": (e) -> @navigate('/schedules')

    "focus .sinput": (e) ->
        $(e.target).stop().animate({width:'300px'},400)
        window.setTimeout( =>
            @sinput.val("")
          ,1)
    "blur .sinput": (e) ->
        $(e.target).stop().animate({width:'50px'},400)
        @sinput.val("")

    "keypress": (e) ->
        console.log(e.keyCode)
        e = e||window.event
        if e.keyCode in [27, 13]
          @sinput.blur()

  constructor: ->
    super

  render: ->
    @html(require('views/navbar')())

  show: (param) ->
    if param is 'list'
      @addSchedule.css('display','block')
      @addTask.css('display','none')
      @refreshAll.css('display','none')
    else
      @addTask.css('display','block')
      @refreshAll.css('display','block')
      @addSchedule.css('display','none')

module.exports = Navbar
