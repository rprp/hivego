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
    Schedule.bind("refresh",  @getData)
    @active @change

  getData: ->
      console.log(Schedule.find(id:1))

  change: (params) =>

    a=Schedule.fetch({id:'0'})
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

    Raphael.getColor()
    Raphael.getColor()
    color = Raphael.getColor()

    ts= new TaskSymbol(paper,120,185,"ttttttttttttt",color)

    Raphael.getColor()
    Raphael.getColor()
    Raphael.getColor()
    color = Raphael.getColor()
    ts1= new TaskSymbol(paper,120,285,"tt",color)
    ts2= new TaskSymbol(paper,220,285,"tt",color)
    ts3= new TaskSymbol(paper,320,285,"tt",color)
    ts4= new TaskSymbol(paper,420,285,"tt",color)


    Raphael.getColor()
    Raphael.getColor()
    Raphael.getColor()
    color = Raphael.getColor()
    ts5= new TaskSymbol(paper,420,385,"tt",color)


    Raphael.getColor()
    Raphael.getColor()
    Raphael.getColor()
    color = Raphael.getColor()
    ts6= new TaskSymbol(paper,420,485,"tt",color)
    ts7= new TaskSymbol(paper,520,485,"tt",color)

    ts.addNext(ts1)
    ts.addNext(ts2)
    ts.addNext(ts3)
    ts.addNext(ts4)

    ts1.addNext(ts5)
    ts2.addNext(ts5)
    ts3.addNext(ts5)
    ts4.addNext(ts6)
    ts4.addNext(ts7)
    ts5.addNext(ts6)
    ts5.addNext(ts7)

#class ScheduleSymbol
  #constructor: (@paper, @cx, @cy, @name, @color) ->


#class JobSymbol
  #constructor: (@paper, @cx, @cy, @name, @color) ->


class TaskSymbol
  constructor: (@paper, @cx, @cy, @name, @color) ->
    @pre=[]
    @preRel=[]

    @next=[]
    @nextRel=[]

    @paper.setStart()
    @sp=@paper.circle(@cx, @cy, 30)
    @sp.ts=@
    @sp.cx = @cx
    @sp.cy = @cy
    @sp.attr({fill: color, stroke: color, "fill-opacity": 0.2, "stroke-width": 2, cursor: "move"})

    @sp.refresh = ->
      if @ts.nextRel
        for r,i in @ts.nextRel
          @paper.connection(r)
      if @ts.preRel
        for r,i in @ts.preRel
          @paper.connection(r)

    @sp.drag(@move, @dragger, @up)

    @text = @paper.text(@cx, @cy, @name)
    @text.toBack()
    @text.attr({fill: color, stroke: "none", "font-size": 15, "fill-opacity": 1, "stroke-width": 1, cursor: "move"})
    @sp.pair=@text

    an = Raphael.animation({"fill-opacity": .2}, 200)
    @sp.animate(an.repeat(10)) 

    st = @paper.setFinish()
    @sp

  addNext: (ts) ->
    @next.push(ts)
    r=@paper.connection(@sp,ts.sp,@sp.attr('fill'),"#{@sp.attr('fill')}|2")
    @nextRel.push(r)

    ts.pre.push(@)
    ts.preRel.push(r)

  click: ->
    alert(@.data('a'))

  dragger: ->
    @ox = @attr("cx")
    @oy = @attr("cy")
    @animate({"fill-opacity": .5}, 500) if @type isnt "text"

    @pair.ox = @attr("x")
    @pair.oy = @attr("y")
    @pair.animate({"fill-opacity": .2}, 500) if @pair.type isnt "text"

  move: (dx, dy) ->
    att =
        cx:@ox + dx
        cy:@oy + dy
    @attr(att)
    att =
        x:@ox + (dx)
        y:@oy + (dy)
    @pair.attr(att)
    @refresh()

  up: ->
    @animate({"fill-opacity": 0.2}, 500) if @type isnt "text"
    @pair.animate({"fill-opacity": 0.2}, 500) if @pair.type isnt "text"

module.exports = ScheduleInfo
