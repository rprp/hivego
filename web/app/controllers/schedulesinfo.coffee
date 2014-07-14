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
    Schedule.bind("findRecord",  @draw)
    @active @change

  change: (params) =>
    Schedule.fetch({Id:params.id})
    @render()

  render: =>
    @html(require('views/schedule-show-info')())
    
  draw: (rs) =>
    @item = Schedule.find(rs.Id)

    paper = Raphael(@pant.get(0),'100%','100%')

    @width = parseFloat(@pant.css("width"))
    @height = parseFloat(@pant.css("height"))

    @ssl = new ScheduleSymbol(paper,@width,@height,@item)

    #Raphael.getColor()# {{{
    #Raphael.getColor()
    #color = Raphael.getColor()

    #ts= new TaskSymbol(paper,120,185,"ttttttttttttt",color)
    #Raphael.getColor()
    #Raphael.getColor()
    #Raphael.getColor()
    #color = Raphael.getColor()
    #ts1= new TaskSymbol(paper,120,285,"tt",color)
    #ts2= new TaskSymbol(paper,220,285,"tt",color)
    #ts3= new TaskSymbol(paper,320,285,"tt",color)
    #ts4= new TaskSymbol(paper,420,285,"tt",color)


    #Raphael.getColor()
    #Raphael.getColor()
    #Raphael.getColor()
    #color = Raphael.getColor()
    #ts5= new TaskSymbol(paper,420,385,"tt",color)


    #Raphael.getColor()
    #Raphael.getColor()
    #Raphael.getColor()
    #color = Raphael.getColor()
    #ts6= new TaskSymbol(paper,420,485,"tt",color)
    #ts7= new TaskSymbol(paper,520,485,"tt",color)

    #ts.addNext(ts1)
    #ts.addNext(ts2)
    #ts.addNext(ts3)
    #ts.addNext(ts4)

    #ts1.addNext(ts5)
    #ts2.addNext(ts5)
    #ts3.addNext(ts5)
    #ts4.addNext(ts6)
    #ts4.addNext(ts7)
    #ts5.addNext(ts6)
    #ts5.addNext(ts7)
    #ts.addNext(ts7)# }}}

class ScheduleSymbol
  constructor: (@paper, @width, @height, @item) ->

    st = Raphael.animation({"fill-opacity": .2}, 2000, ->
            @.animate(ed))
    ed = Raphael.animation({"fill-opacity": .5}, 2000, ->
            @.animate(st))

    title = @paper.text(@width/2, 30, @item.Name)
    title.attr({fill: "#333", stroke: "none", "font-size": 30, "fill-opacity": 1, "stroke-width": 2})

    jobflg=@paper.circle(@width-60, 100, 15)
    jobflg.attr({fill: "green", stroke: "green", "fill-opacity": 0.5, "stroke-width": 1, cursor: "hand"})
    jobflg.animate(st)
    jobflg.hover (-> @.attr({r: 17})), (-> @.attr({r: 15}))

    color=['#FFD700', '#008000', '#DA70D6', '#98FB98', '#6495ED', '#B8860B', '#2F4F4F', '#FF4500', '#AFEEEE', '#DB7093', 
        '#CD853F', '#FFC0CB', '#B0E0E6', '#BC8F8F', '#4169E1', '#8B4513', '#00FFFF', '#FF8C00', '#00BFFF', '#008B8B', 
        '#ADFF2F', '#4B0082', '#F0E68C', '#7CFC00', '#0000FF', '#8A2BE2', '#7FFF00', '#DEB887', '#5F9EA0', '#D2691E', '#A9A9A9',
        '#8B008B', '#556B2F', '#9932CC', '#8FBC8B', '#483D8B', '#00CED1', '#9400D3', '#FF69B4', '#228B22', '#1E90FF', '#FF00FF',
        '#FFB6C1', '#FFA07A', '#20B2AA', '#87CEFA', '#00FF00', '#B0C4DE', '#FF00FF', '#32CD32', '#0000CD', '#66CDAA', '#BA55D3',
        '#9370DB', '#3CB371', '#7B68EE', '#00FA9A', '#48D1CC', '#C71585', '#191970', '#000080', '#808000', '#6B8E23', '#FFA500',
        '#F4A460', '#2E8B57', '#A0522D', '#87CEEB', '#6A5ACD', '#708090', '#00FF7F', '#4682B4', '#D2B48C', '#008080', '#40E0D0',
         '#006400', '#BDB76B','#EE82EE', '#F5DEB3', '#FFFF00', '#9ACD32']

    #jobflg.click =>
        #for job,i in @item.Job
          #jobname = @paper.text(@width-100, 30 + (i*30), job.Name)
          #jobname.attr({fill: "#333", stroke: "none", "font-size": 18, "fill-opacity": 1, "stroke-width": 2})

    top = 0
    @ts = []
    for job,i in @item.Job
      top += 120

      spacing = (@width-100)/job.Tasks.length if job.Tasks.length > 0
      if spacing < 100 and spacing >50 
        r = 15
        spacing = 60
      else if spacing <= 50 
        r = 8
        spacing = 40
      else
        r = 25
        spacing = 100

      left = (@width-100)/2-(job.Tasks.length/2) * spacing if job.Tasks.length > 0
      for task,j in job.Tasks
        t= new TaskSymbol(paper,left,top,task.Name,color[i],r)
        t.Id = task.Id
        t.RelTaskId = (rt.Id for rt in task.RelTasks)
        @ts.push(t)

        rts.addNext(t) for rts in @getTaskSymbol(t.RelTaskId)
        left += spacing
    
  getTaskSymbol: (Ids) ->
     t for t in @ts when t.Id in Ids

class JobSymbol
  constructor: (@paper, @cx, @cy, @name, @color) ->


class TaskSymbol
  constructor: (@paper, @cx, @cy, @name, @color, @r=20) ->
    @pre=[]
    @preRel=[]

    @next=[]
    @nextRel=[]

    @paper.setStart()
    @sp=@paper.circle(@cx, @cy, @r)
    @sp.ts=@
    @sp.cx = @cx
    @sp.cy = @cy
    @sp.hover(@hoveron,@hoverout)
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
    a = Raphael.animation({"stroke-width": 6}, 600)

    @sp.animate(a)
    for r in @nextRel
      r.line.animate(a)
    for n in @next
      n.sp.animate(a)
    for rp in @preRel
      rp.line.animate(a)
    for p in @pre
      p.sp.animate(a)

  hoverout: =>
    b = Raphael.animation({"stroke-width": 1}, 600)

    @sp.animate(b)
    for r in @nextRel
      r.line.animate(b)
    for n in @next
      n.sp.animate(b)
    for rp in @preRel
      rp.line.animate(b)
    for p in @pre
      p.sp.animate(b)
      
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
