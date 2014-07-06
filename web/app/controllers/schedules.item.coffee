Spine = require('spine')
Schedule = require('models/schedule')
$       = Spine.$

class ScheduleItem extends Spine.Controller
  className: 'scheduleitem'
  
  events:
   "click ": "show"
   "click .cyc": "showcyc"
   "mouseover .panel": "mouseover"
   "mouseout .panel": "mouseout"
   
  elements:
    ".panel":     "panel"
    "#jobcnt":     "jobcnt"
    
  constructor: ->
    super
    @el.addClass('col-sm-3')
    throw "@item required" unless @item
    @item.bind("update", @render)
    @item.bind("destroy", @remove)
  
  render: (item) =>
    @item = item if item
    @html(@template(@item))
    @
    
  template: (items) ->
    require('views/schedule.show')(items)

  remove: ->
    @el.remove()

  show: ->

  showcyc: ->
    alert('！')

  mouseover: ->
      @panel.css("box-shadow","0 0 22px #777")
      @jobcnt.css("background-color","#d9edf7")

  mouseout: ->
      @panel.css("box-shadow","")
      @jobcnt.css("background-color","")

module.exports = ScheduleItem
