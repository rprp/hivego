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

    [@width, @height] = [parseFloat(@pant.css("width")), parseFloat(@pant.css("height"))]

    @ssl = new ScheduleSymbol(paper,@width,@height,@item)

class ScheduleSymbol
  constructor: (@paper, @width, @height, @item) ->

    @color=['#FF8C00', '#008000', '#2F4F4F', '#DA70D6', '#0000FF', '#8A2BE2', '#6495ED', '#B8860B', '#FF4500', '#AFEEEE', '#DB7093', 
        '#CD853F', '#FFC0CB', '#B0E0E6', '#BC8F8F', '#4169E1', '#8B4513', '#00FFFF', '#00BFFF', '#008B8B', 
        '#ADFF2F', '#4B0082', '#F0E68C', '#7CFC00', '#7FFF00', '#DEB887', '#98FB98', '#FFD700', '#5F9EA0', '#D2691E', '#A9A9A9',
        '#8B008B', '#556B2F', '#9932CC', '#8FBC8B', '#483D8B', '#00CED1', '#9400D3', '#FF69B4', '#228B22', '#1E90FF', '#FF00FF',
        '#FFB6C1', '#FFA07A', '#20B2AA', '#87CEFA', '#00FF00', '#B0C4DE', '#FF00FF', '#32CD32', '#0000CD', '#66CDAA', '#BA55D3',
        '#9370DB', '#3CB371', '#7B68EE', '#00FA9A', '#48D1CC', '#C71585', '#191970', '#000080', '#808000', '#6B8E23', '#FFA500',
        '#F4A460', '#2E8B57', '#A0522D', '#87CEEB', '#6A5ACD', '#708090', '#00FF7F', '#4682B4', '#D2B48C', '#008080', '#40E0D0',
         '#006400', '#BDB76B','#EE82EE', '#F5DEB3', '#FFFF00', '#9ACD32']

    #[@st, @ed] = [Raphael.animation({"fill-opacity": .2}, 2000, -> @.animate(ed)),
                #Raphael.animation({"fill-opacity": .5}, 2000, -> @.animate(st))]
    [@st, @ed] = [Raphael.animation({"fill-opacity": .2}, 1000),
                Raphael.animation({"fill-opacity": .5}, 1000)]

    slider = @paper.path("M #{@width-220},10L #{@width-220},#{@height}")
    slider.attr({fill: "#333", "fill-opacity": 0.3, "stroke-width": 3, "stroke-opacity": 0.1})
    
    top = @printSchedule(25, @width-205)
    @printJob(top, @width-205)

    #jobflg=@paper.circle(@width-60, 100, 15)
    #jobflg.attr({fill: "green", stroke: "green", "fill-opacity": 0.5, "stroke-width": 1, cursor: "hand"})
    #jobflg.animate(st)
    #jobflg.hover (-> @.attr({r: 17})), (-> @.attr({r: 15}))

    #jobflg.click =>

    top = 40
    @ts = []
    for job,i in @item.Job

      spacing = (@width-200)/job.Tasks.length if job.Tasks.length > 0
      if spacing < 100 and spacing >50 
        r = 15
        spacing = 60
      else if spacing <= 50 
        r = 8
        spacing = 40
      else
        r = 25
        spacing = 100

      left = (@width-200)/2-(job.Tasks.length/2) * spacing if job.Tasks.length > 0
      for task,j in job.Tasks
        t= new TaskSymbol(paper,left,top,task.Name,@color[i],r)
        t.Id = task.Id
        t.RelTaskId = (rt.Id for rt in task.RelTasks)
        @ts.push(t)

        rts.addNext(t) for rts in @getTaskSymbol(t.RelTaskId)
        left += spacing

      top += 120
    
  getTaskSymbol: (Ids) ->
     t for t in @ts when t.Id in Ids

  printSchedule: (top, left)->
    ff = "Helvetica, Tahoma, Arial, STXihei, '华文细黑', Heiti, '黑体', 'Microsoft YaHei', '微软雅黑', SimSun, '宋体', sans-serif"

    #标题，调度名称
    [top,left] = [top + (@item.Name.length//7) * 20, left]
    title = @paper.text(left, top, @item.SplitName(7))
    title.attr({fill: "#333", "text-anchor": "start", stroke: "none", "font-size": 22, "fill-opacity": 1, "stroke-width": 2})

    att = {fill: "#333", "font-family":ff, "text-anchor": "start", stroke: "none", "font-size": 14, "fill-opacity": 1, "stroke-width": 1}
    
    [top,left] = [top + 30 + (@item.Name.length//7) * 20, left]
    #调度周期
    cyc = @paper.text(left, top, "调度周期：#{@item.GetCyc()}")
    cyc.attr(att)

    #调度时间
    gs=@item.GetSecond()
    [top,left] = [top+30, left]
    @paper.text(left, top, "启动时间：").attr(att)

    for ss in gs
      [top,left] = [top+30, left]
      c = @paper.text(left+20, top, "#{ss}")
      c.attr(att)

    #任务数量
    [top,left] = [top+30, left]
    @paper.text(left, top, "任务数量：#{@item.TaskCnt}").attr(att)

    #下次执行时间
    [top,left] = [top+30, left]
    @paper.text(left, top, "下次执行：#{@item.GetNextStart()}").attr(att)

    #当前状态

    #所有者

    [top,left] = [top+30, left]
    betweenline = @paper.path("M #{left},#{top}L #{@width-30},#{top}")
    betweenline.attr({stroke: "#A0522D", "stroke-width": 2, "stroke-opacity": 0.2})

    top


  printJob: (top,left)->
    ff = "Helvetica, Tahoma, Arial, STXihei, '华文细黑', Heiti, '黑体', 'Microsoft YaHei', '微软雅黑', SimSun, '宋体', sans-serif"

    [top,left]=[top+30,left]
    jtitle = @paper.text(left, top, "任务列表")
    jtitle.attr({fill: "#333", "text-anchor": "start", stroke: "none", "font-size": 18, "fill-opacity": 1, "stroke-width": 2})
    [top,left]=[top+40,left]
    for job,i in @item.Job when job.Id isnt 0
      jobname = @paper.text(left+80, top, job.Name)
      jobname.attr({stroke: @color[i], "fill": @color[i], "font-family":ff , "font-size": 16, "stroke-opacity":1, "fill-opacity": 1, "stroke-width": 1})
      jobcir = @paper.circle(left+25,top,15)
      jobcir.attr({fill: @color[i], stroke: @color[i], "fill-opacity": 0.2, "stroke-width": 1})
      [top,left]=[top+50,left]

    addbtn = @paper.path("M25.979,12.896 19.312,12.896 19.312,6.229 12.647,6.229 12.647,12.896 5.979,12.896 5.979,19.562 12.647,19.562 12.647,26.229 19.312,26.229 19.312,19.562 25.979,19.562z")
    addbtn.transform("t#{left+60},#{top-10}s2")
    addbtn.attr({fill: "#00FF00", stroke: "#00FF00", "fill-opacity": 0.1, "stroke-opacity":0.2,  "stroke-width": 1, cursor: "hand"})


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

module.exports = ScheduleInfo
