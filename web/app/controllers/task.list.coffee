Spine = require('spineify')
Raphael = require('raphaelify')
Eve = require('eve')
$       = Spine.$

class TaskManager
  constructor: (@paper, @color, @item, @width, @height) ->
    @set = @paper.set()
    top = 80
    @ts = []
    for job,i in @item.Job

      spacing = (@width-200)/job.Tasks.length if job.Tasks.length > 0
      r = 25
      spacing = 100

      left = (@width-200)/2-(job.Tasks.length/2) * spacing if job.Tasks.length > 0
      for task,j in job.Tasks
        t= new TaskSymbol(paper,left,top,task.Name,@color[i],r)
        t.Id = task.Id
        t.JobId = job.Id
        t.RelTaskId = (rt.Id for rt in task.RelTasks)
        @ts.push(t)
        @set.push(t.sp)
        @set.push(t.text)

        rts.addNext(t) for rts in @getTaskSymbol(t.RelTaskId)
        left += spacing

      top += 120

  getTaskSymbol: (Ids) ->
     t for t in @ts when t.Id in Ids

  hlight: (Id) ->
    a = Raphael.animation({"fill-opacity": 0.7}, 500)
    for t in @ts
      if t.JobId is Id
        t.sp.animate(a)
        t.sp.g = t.sp.glow({color: t.sp.attr("fill"), "fill-opacity": 0.2, width:10})

  nlight: (Id) ->
    a = Raphael.animation({"fill-opacity": 0.2}, 500)
    for t in @ts
      if t.JobId is Id
        t.sp.animate(a)
        t.sp.g.remove()


class TaskSymbol
  constructor: (@paper, @cx, @cy, @name, @color, @r=20) ->
    @pre=[]
    @preRel=[]

    @next=[]
    @nextRel=[]

    @paper.setStart()
    @sp=@paper.circle(@cx, @cy, @r)
    @sp.ts=@
    @sp.hover(@hoveron,@hoverout)
    @sp.attr({fill: @color, stroke: @color, "fill-opacity": 0.2, "stroke-width": 1, cursor: "move"})

    @sp.refresh = ->
      if @ts.nextRel then for r,i in @ts.nextRel
        @paper.connection(r)

      if @ts.preRel then for r,i in @ts.preRel
        @paper.connection(r)

    @sp.drag(@move, @dragger, @up)

    @text = @paper.text(@cx, @cy, @name)
    @text.toBack()
    @text.attr({fill: "#333", stroke: "none", "font-size": 15, "fill-opacity": 1, "stroke-width": 1, cursor: "move"})
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

  hoveron: =>
    a = Raphael.animation({"stroke-width": 6, "fill-opacity": 0.5}, 300)

    @sp.animate(a)
    r.line.animate(a)  for r in @nextRel
    n.sp.animate(a)    for n in @next
    rp.line.animate(a) for rp in @preRel
    p.sp.animate(a)    for p in @pre
      
  hoverout: =>
    b = Raphael.animation({"stroke-width": 1,"fill-opacity": 0.2}, 300)

    @sp.animate(b)
    r.line.animate(b)  for r in @nextRel
    n.sp.animate(b)    for n in @next
    rp.line.animate(b) for rp in @preRel
    p.sp.animate(b)    for p in @pre
      
  dragger: ->
    [@ox, @oy]  = [@attr("cx"), @attr("cy")]
    @animate({"fill-opacity": .5}, 500) if @type isnt "text"

    [@pair.ox, @pair.oy] = [@attr("x"),@attr("y")]
    @pair.animate({"fill-opacity": .2}, 500) if @pair.type isnt "text"

  move: (dx, dy) ->
    @attr([ cx:@ox + dx, cy:@oy + dy])
    @pair.attr([x:@ox + dx, y:@oy + dy])
    @refresh()

  up: ->
    @animate({"fill-opacity": 0.2}, 500) if @type isnt "text"
    @pair.animate({"fill-opacity": 0.2}, 500) if @pair.type isnt "text"

module.exports = TaskManager
