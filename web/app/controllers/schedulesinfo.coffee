Spine = require('spine')
Raphael = require('raphael')
Schedule = require('models/schedule')
$       = Spine.$

class ScheduleInfo extends Spine.Controller
  className: 'scheduleinfo'

  elements:
    ".pant":          "pant"

  constructor: ->
    super
    @active @change

  change: (params) =>
    @item = Schedule.find(params.id)
    @render(@item)

  render: (item) =>
    @item = item if item
    @html(@template(@item))
    @draw()
    
  template: (items) ->
    require('views/schedule-show-info')(items)

  draw: ->
    paper = Raphael(@pant.get(0),'100%','100%')
    circle = paper.circle(200,40,50)
    cc = paper.text(200,30,"xxxx")
    circle.attr("fill","#f00")

module.exports = ScheduleInfo
