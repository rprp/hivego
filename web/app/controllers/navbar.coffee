Spine = require('spineify')
Schedule = require('models/schedule')
ScheduleManager = require('controllers/schedule.info')
$       = Spine.$

class Navbar extends Spine.Controller
  className: 'hnavbar'
  
  #elements:

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
    s = new Schedule({Id:-1})
    @sf = new ScheduleManager.SForm("c",s)
    @append(@sf.render(550,100,s))
    @sf.el.css("z-index",1000)

  showAddSchedule: (e) ->
    e = e||window.event
    s = new Schedule({Id:-1})
    @sf = new ScheduleManager.SForm("c",s)
    @append(@sf.render(550,100,s))
    @sf.el.css("z-index",1000)

  mouseover: (e) ->
    $(e.target).stop().animate({backgroundColor:'#777'},200)


  mouseout: (e) ->
    $(e.target).stop().animate({backgroundColor:'#333'},200)

module.exports = Navbar
