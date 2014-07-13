Spine = require('spine')
Raphael = require('raphael')
Eve = require('eve')
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

    sp = paper.rect(100, 20, 800, 200, 6)
    sp.attr({stroke: "blue", "fill-opacity": 0, "stroke-width": 1})

    sp1=paper.rect(100,250,800,200)
    sp1.attr({stroke: "blue", "fill-opacity": 0, "stroke-width": 1})

    rect= new Shape(paper)

    t1 = rect.draw(120,185,60,30,"task 1")
    t2 = rect.draw(200,185,60,30,"task 2")
    t3 = rect.draw(280,185,60,30,"task 3")
    t4 = rect.draw(360,185,60,30,"task 4")
    t5 = rect.draw(440,185,60,30,"task 5")
    t6 = rect.draw(520,185,60,30,"task aaaaaaaaaaaaaaaaa\ndddddd\ndddd\nddddf6")

    t6.insertAfter(t1)

    t1.rel=[t2,t3]
    t1.re=t1.connect()
    t2.re=t1.re
    t3.re=t1.re

class TaskSymbol
  constructor: (@paper, @x, @y, @name) ->

  click: ->
    alert(@.data('a'))

  dragger: ->
    @ox = @attr("x")
    @oy = @attr("y")
    @animate({"fill-opacity": .2}, 500) if @type isnt "text"

    @pair.ox = @attr("x")
    @pair.oy = @attr("y")
    @pair.animate({"fill-opacity": .2}, 500) if @pair.type isnt "text"

  move: (dx, dy) ->
    att =
        x:@ox + dx
        y:@oy + dy
    @attr(att)
    att =
        x:@pair.ox + (dx + @width/2)
        y:@pair.oy + (dy + @height/2)
    @pair.attr(att)
    @refresh()

  up: ->
    @animate({"fill-opacity": 0}, 500) if @type isnt "text"
    @pair.animate({"fill-opacity": 0}, 500) if @pair.type isnt "text"

  draw: (x, y, width, height, text) ->
      @paper.setStart()
      color = Raphael.getColor()
      @sp=@paper.rect(x, y, width, height, 10)
      @sp.width = width
      @sp.height = height
      @sp.attr({fill: color, stroke: color, "fill-opacity": 0, "stroke-width": 2, cursor: "move"})

      @sp.connect = ->
        if @rel
          for r,i in @rel
            @paper.connection(@,r,"#333")

      @sp.refresh = ->
        if @re
          for r,i in @re
            @paper.connection(r)

      @sp.drag(@move, @dragger, @up)

      @text = @paper.text(x+width/2, y+height/2, text)
      @text.width = width

      @text.height = height
      @text.toBack()
      @text.attr({fill: color, stroke: "none", "font-size": 15, "fill-opacity": 1, "stroke-width": 1, cursor: "move"})
      @sp.pair=@text

      an = Raphael.animation({"fill-opacity": .2}, 200)
      @sp.animate(an.repeat(10)) 

      st = @paper.setFinish()
      @sp

module.exports = ScheduleInfo
