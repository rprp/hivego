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

  events:
    "mouseenter h1":   "mouseover"
    "mouseleave h1":   "mouseout"

    "click #addSchedule": "showAddSchedule"
    "click #addTask": "showAddTask"

  constructor: ->
    super

  render: ->
    @html(require('views/navbar')())

  showAddTask: (e) ->
    e = e||window.event
    @trigger('addtask',e)

  showAddSchedule: (e) ->
    e = e||window.event
    s = new Schedule({Id:-1})
    @sf = new ScheduleManager.Form("c",s)
    @append(@sf.render(550,100,s))
    @sf.el.css("z-index",1000)

  mouseover: (e) ->
    $(e.target).stop().animate({backgroundColor:'#777'},400)


  mouseout: (e) ->
    $(e.target).stop().animate({backgroundColor:'#333'},200)

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
