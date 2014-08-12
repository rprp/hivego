Spine = require('spineify')
Ajax  = Spine.Ajax.Base
Raphael = require('raphaelify')
Eve = require('eve')
Style = require('controllers/style')
Task = require('models/task')
Job = require('models/job')
$       = Spine.$

class TaskManager extends Spine.Controller
  elements:
    "#taskName": "taskName"
    "#taskAddr": "taskAddr"
    "#taskid": "taskId"
    "#taskCmd":  "taskCmd"
    ".taskParam":"taskParam"
    "#taskDesc":"taskDesc"
    ".tcyclbl": "cycGroup"
    ".startList":"startList"
    ".taskParamList":"taskParamList"
    "#jobid": "JobId"
    "#confirmdeltaskrel": "confirmdeltaskrel"

  events:
    "click .tclose"       :  "hideTask"
    "click #delrelclose"  :  "hideDelTaskRel"
    "click #deltaskrel"   :  "postDelTaskRel"
    "click .tparam"       :  "editParam"
    "click .addParam"     :  "appendParam"
    "click .delParam"     :  "delParam"
    "click #submitTask": "postTask"
    "click .jobli": "setJob"

    "keypress .addTask":  "addTaskKeyPress"
    "keypress .taskParam":  "paramKeyPress"
    "blur .taskParam":  "setTaskParamVal"

    "mouseenter .list-group-item":  "showDelParam"
    "mouseleave .list-group-item":  "hideDelParam"

    "mousedown .addTaskHead": "setMoveFlg"
    "mouseup .addTaskHead": "clearMoveFlg"
    "mousemove .addTaskHead": "movePanel"

  constructor: (@paper, @color, @item, @w, @h) -># {{{
    super
    Spine.bind("connectTaskStart", @connectStart)
    Spine.bind("connectTaskFinish", @connectFinish)
    Spine.bind("deleteTaskRelStart", @delTaskRelStart)
    Spine.bind("deleteTaskRel", @addRemoveTaskRel)
    Spine.bind("deleteTask", @deleteTask)
    @setpp = @paper.set()
    @isRefresh = true
    @isMove = false
    @jobList = []
    @taskList = []
    @delTaskRels = []
    @currentTask
    @relTask
    top = 80
    
    if @item.Jobs
      @refreshTaskList(top)
  # }}}

  refreshTaskList: (top) =># {{{
    for jb,i in @item.Jobs
      job = new Job()
      for key, value of jb
        job[key] = value
      @jobList.push(job)

      tasks = (v for k,v of job.Tasks)
      spacing = 100
  
      left = (@w-200)/2-(tasks.length/2) * spacing if tasks.length > 0
      for task,j in tasks
        tk = new Task()
        for key, value of task
          tk[key] = value
        tk.JobNo = i
        t= new TaskShape(@paper,left,top,tk,@color[i],25)
        t.conn.drag(t.connMove, t.connDragger, t.connUp,@)
        @taskList.push(t)
        
        @setpp.push(t.sp)
        @setpp.push(t.text)
  
        rts.addNext(t) for rts in @getTaskShape(t.RelTaskId)
        left += spacing
      top += 100
  # }}}

  showDelParam: (e) -># {{{
    $(e.target).children(".delParam").css("display","")
  # }}}

  hideDelParam: (e) -># {{{
    $(e.target).children(".delParam").css("display","none")
  # }}}

  delParam: (e) -># {{{
    $(e.target).parent("li").remove()
  # }}}

  appendParam: -># {{{
    @taskParamList.append(require('views/task-param')())
    $(".taskParam").focus()
  # }}}

  getTaskShape: (Ids) -># {{{
     t for t in @taskList when t.task.Id in Ids
  # }}}

  addTaskKeyPress: (e) -># {{{
    e = e||window.event
    if e.ctrlKey and e.keyCode in [13,10]
      @postTask(e)
  # }}}

  postTask: (e) -># {{{
    @el.css("display","none")
    if @taskId.val()
      tk = t.task for t in @taskList when t.task.Id is parseInt(@taskId.val())
      tk.bind("ajaxSuccess",@updateTaskAndRefresh)
      tk.Name = @taskName.val()
      tk.JobId = parseInt(@JobId.val())
      tk.Address = @taskAddr.val()
      tk.Cmd = @taskCmd.val()
      tk.Desc = @taskDesc.val()

      tk.Param = []
      for li,i in @taskParamList.children("li")
        tp = $(li).children(".taskParam").val()
        if tp isnt ""
          tk.Param.push(tp)
      tk.save({method: "PUT", url:"/schedules/#{@item.Id}/jobs/#{tk.JobId}/tasks/#{tk.Id}"})
    else
      tk = new Task()
      tk.bind("ajaxSuccess",@addTaskAndRefresh)
      tk.Name = @taskName.val()
      tk.Address = @taskAddr.val()
      tk.Cmd = @taskCmd.val()
      tk.Desc = @taskDesc.val()
      tk.Id = -1
      tk.JobId = parseInt(@JobId.val())

      tk.Param = []
      for li,i in @taskParamList.children("li")
        tp = $(li).children(".taskParam").val()
        if tp isnt ""
          tk.Param.push(tp)

      tk.create({url:"/schedules/#{@item.Id}/jobs/0/tasks"}) if tk.Name
  # }}}

  updateTaskAndRefresh: (task, status, xhr) =># {{{
    s = Raphael.animation({"fill-opacity": 1, "stroke-width": 6}, 1200, -> @.animate(e))
    e = Raphael.animation({"fill-opacity": .2, "stroke-width": 1}, 300)
    if xhr is "success"
      tp = t for t in @taskList when t.task.Id is parseInt(task.Id)
      tp.task = task
      tp.sp.animate(s)

      @isRefresh = true
  # }}}

  addTaskAndRefresh: (task, status, xhr) =># {{{
    if xhr is "success"
      Spine.Module.extend.call(task, Task)

      ci = i for j,i in @item.Jobs when j.Id is task.JobId

      task.JobNo = ci
      t= new TaskShape(@paper,150,0,task,Style.getRgbColor()[ci],25)
      t.conn.drag(t.connMove, t.connDragger, t.connUp,@)
      t.sp.animate({"cx": 150, "cy": ci*100+80}, 2000, "elastic")
      t.text.animate({"x": 150, "y": ci*100+80}, 2000, "elastic")

      for el in t.toolset
        el.attr([cy:ci*100+80]) if el.attr("cx")
        el.attr([y:ci*100+80]) if el.attr("x")

      @taskList.push(t)
      
      @setpp.push(t.sp)
      @setpp.push(t.text)

      @isRefresh = true
  # }}}

  hlight: (Id) -># {{{
    a = Raphael.animation({"fill-opacity": 0.5}, 500)
    for t in @taskList
      if t.task.JobId is Id
        t.sp.animate(a)
        t.sp.g = t.sp.glow({color: t.sp.attr("fill"), "fill-opacity": 0.2, width:10})
  # }}}

  nlight: (Id) -># {{{
    a = Raphael.animation({"fill-opacity": 0.2}, 500)
    for t in @taskList
      if t.task.JobId is Id
        t.sp.animate(a)
        t.sp.g.remove()
  # }}}

  render: (task) -># {{{
    unless task.Name
      task = new Task()
      task.Param=[]
    
    task.JobList = @item.Jobs
    [task.JobName,task.JobNo] = [n.Name,i] for n,i in @item.Jobs when n.Id is task.JobId
    task.RgbColor = Style.getRgbColor()
    @html(require("views/task")(task))

    cs = c for c in @cycGroup when c.textContent is @item.GetCyc()
    $(cs).removeClass("label-default")
    $(cs).addClass("label-success")
    $(cs).css("display","none")
    $(cs).prevAll().css("display","none")

    window.setTimeout( =>
            @taskName.focus()
      ,500)
    @el.css("display","block")
    @el.css("position","absolute")
    @el.css("left",200)
    @el.css("top", 60)

# }}}

  setCyc: (e) -># {{{
    @cycGroup.removeClass("label-success")
    @cycGroup.addClass("label-default")

    $(e.target).removeClass("label-default")
    $(e.target).addClass("label-success")

    @item.SetCyc($(e.target).text())
  # }}}

  hideTask: -># {{{
    @el.css("display","none")
# }}}

  setMoveFlg: (e) -># {{{
    @isMove = true
    @preLeft = e.clientX
    @preTop = e.clientY
  # }}}

  clearMoveFlg: (e) -># {{{
    @isMove = false
    @el.css("opacity", 1)
  # }}}

  movePanel: (e) -># {{{
    return unless @isMove
    e = e||window.event

    dx = (e.clientX - @preLeft) + parseInt(@el.css("left"))
    dy = (e.clientY - @preTop) + parseInt(@el.css("top"))
    @el.css("left", dx)
    @el.css("top", dy)
    @el.css("opacity", 0.4)

    @preLeft = e.clientX
    @preTop = e.clientY
  # }}}

  editParam: (e) -># {{{
    $(e.target).siblings().not(".delParam").css("display","")
    $(e.target).siblings().focus()
    $(e.target).css("display","none")
  # }}}

  setTaskParamVal: (e) -># {{{
    e = e||window.event
    $(e.target).css("display","none")
    $(e.target).siblings().not(".delParam").css("display","")
    $(e.target).siblings().not(".delParam").text(" #{$(e.target).val()}           ")
  # }}}

  paramKeyPress: (e) -># {{{
    e = e||window.event
    if e.keyCode in [13,10]
      @setTaskParamVal(e)
  # }}}

  setJob: (e) -># {{{
    $('#jobid').val(@$(e.target).attr("data"))
    $('.jobbtn').text(@$(e.target).text())
    $('.jobbtn').css("background-color",@$(e.target).css("background-color"))
    $('.jobbl').css("background-color",@$(e.target).css("background-color"))
  # }}}
  
  delTaskRelStart: (ts, e) =>
    s1 = Raphael.animation({"fill-opacity": .05, "stroke-width": 0}, 800)
    @currentTask = ts
    @delTaskRelFlg = true

    so = ["stroke-opacity",0]

    for t,i in @taskList
      t.sp.unhover(t.hoveron,t.hoverout)
      t.sp.unmousedown(t.sp.md,t)
      t.sp.unmouseup(t.showTool,t)
      if t isnt ts and t not in ts.pre
        t.sp.animate(s1)
        t.text.animate(s1)
        [r.bg.attr(so...), r.line.attr(so...)] for r,i in t.preRel
        [r.bg.attr(so...), r.line.attr(so...)] for r,i in t.nextRel
      else if t is ts
        t.sp.click(t.sp.ck = ->
            if @delTaskRelFlg
              @delTaskRelEnd()
              @delTaskRelFlg = false
              @confirmdeltaskrel.css("display","none")
          ,@)
        t.sp.attr("cursor","pointer")

        for r,i in t.preRel
          r.head = t.pre[i]
          r.tail = t
          r.line.animate({"stroke-opacity": 0.05, "stroke-width": 12}, 500)
          r.bg.animate({"stroke-opacity": 1, "stroke-width": 2}, 500)
          r.line.hover(r.mouseover = ->
              @animate({"stroke-opacity": 0.8}, 100)
            ,r.mouseout = ->
              @animate({"stroke-opacity": 0.05, "stroke-width": 12}, 200))

          r.line.click(r.click = ->
              @bg.animate({"stroke-opacity": 0.05}, 500)
              @line.animate({"stroke-opacity": 0.05}, 500)
              @line.unhover(@.mouseover,@.mouseout)
              Spine.trigger("deleteTaskRel",@)
            ,r)

          #cc = @paper.circle(0, 0, 3)
          #cc.attr({fill: '#00FF00', stroke:  '#00FF00', "fill-opacity": 1, "stroke-width": 0})
          #$(cc.node).html("<animateMotion fill='freeze' begin='2s' dur='2s' repeatCount='8' path='#{r.line.getSubpath().end}' rotate='auto' />")

    ts.toolset.show()
    ts.toolset.attr({"fill-opacity": 0.1, "stroke-width": 0.5})
    ts.showTool()
    $("svg").css("cursor","url('img/scissors.cur'),auto")

  postDelTaskRel: =>
    ajax = new Ajax()
    for r,i in @delTaskRels
      param = "tasks/#{r.tail.task.Id}/reltask/#{r.head.task.Id}"
      ajax.ajaxQueue(
        {}, {
        type: 'DELETE'
        contentType: 'application/json'
        data: ""
        url: "/schedules/#{@item.Id}/jobs/#{r.tail.task.JobId}/#{param}"
        parallel:{}
        }
      )
      r.head.removeNext(r)

    @delTaskRelEnd()
    @delTaskRelFlg = false
    @confirmdeltaskrel.css("display","none")

  delTaskRelEnd: =>
    s1 = Raphael.animation({"fill-opacity": .2, "stroke-width": 1}, 1200)
    txt = Raphael.animation({"fill-opacity": 1, "stroke-width": 1}, 1200)

    @delTaskRelFlg = false
    ts = @currentTask
    ts.sp.unclick(ts.sp.ck)
    ts.sp.attr("cursor","move")

    so = [{"stroke-opacity":1},1200]

    for t,i in @taskList
      t.sp.hover(t.hoveron,t.hoverout)
      t.sp.mousedown(t.sp.md,t)
      t.sp.mouseup(t.showTool,t)
      if t isnt ts and t not in ts.pre
        t.sp.animate(s1)
        t.text.animate(txt)
        [r.bg.animate(so...), r.line.animate(so...)] for r,i in t.preRel
        [r.bg.animate(so...), r.line.animate(so...)] for r,i in t.nextRel
      else if t is ts
        for r,i in t.preRel
          r.head = t.pre[i]
          r.tail = t
          r.line.animate({"stroke-opacity": 1, "stroke-width": 1}, 800)
          r.bg.animate({"stroke-width": 1}, 800)
          r.line.unhover(r.mouseover, r.mouseout)

          r.line.unclick(r.click)

    ts.toolset.attr({"fill-opacity": 0.1, "stroke-width": 0.5})
    @delTaskRels = []
    $("svg").css("cursor","auto")

  addRemoveTaskRel: (r) =>
    @delTaskRels.push(r)
    $("#delcnt").text(@delTaskRels.length)
    
    if @confirmdeltaskrel.length < 1
      @html(require('views/taskrel')())
      @el.css("position","absolute")
      @el.css("left",r.tail.sp.ox)
      @el.css("top",r.tail.sp.oy+92)
    else
      @confirmdeltaskrel.css("display","block")

  hideDelTaskRel: (e) =>
    @confirmdeltaskrel.css("display","none")
    if @delTaskRelFlg
      @delTaskRelEnd()

  deleteTask: (ts,e) =>
    @taskList = (t for t,i in @taskList when t isnt ts)
    tk = new Task()
    tk.destroy({url:"/schedules/#{@item.Id}/jobs/#{ts.task.JobId}/tasks/#{ts.task.Id}"})
    ts.remove()

  connectStart: (ts, e) =>
    s1 = Raphael.animation({"fill-opacity": .05, "stroke-width": 0}, 200)

    @currentTask = ts
    c = ts.conn

    so = ["stroke-opacity",0]
    so2 = ["stroke-opacity",0.2]
    cnt = [c,ts.sp,ts.sp.attr("fill"),"#{ts.sp.attr('fill')}|4"]

    c.rel = ts.paper.connection(cnt...)
    [c.rel.bg.attr(so2...), c.rel.line.attr(so2...)]
    c.toFront()

    for t,i in @taskList
      t.sp.unhover(t.hoveron,t.hoverout)
      if t.task.Id isnt ts.task.Id and t.task.JobNo >= ts.task.JobNo
        t.sp.animate(s1)
        t.text.animate(s1)
        [r.bg.attr(so...), r.line.attr(so...)] for r,i in t.preRel
        [r.bg.attr(so...), r.line.attr(so...)] for r,i in t.nextRel

    tpre = ts.pre
    while tpre.length > 0
      tmp = []
      for rts,i in tpre
        [rts.sp.animate(s1), rts.text.animate(s1)]
        [r.bg.attr(so...),r.line.attr(so...)] for r,i in rts.nextRel
        if rts.pre.length > 0
          tmp.push(r) for r,j in rts.pre
            
      tpre = tmp

  connectFinish: (ts, e) =>
    s1 = Raphael.animation({"fill-opacity": .2, "stroke-width": 1}, 300)
    txt = Raphael.animation({"fill-opacity": 1, "stroke-width": 1}, 300)
    
    cr = ts.conn.rel
    cr.line.remove()
    cr.bg.remove()
    cr = null
    if @relTask
      ajax = new Ajax()
      param = "tasks/#{ts.task.Id}/reltask/#{@relTask.task.Id}"
      ajax.ajaxQueue(
        {}, {
        type: 'POST'
        contentType: 'application/json'
        data: ""
        url: "/schedules/#{@item.Id}/jobs/#{ts.task.JobId}/#{param}"
        parallel:{}
        }
      )
      @relTask.addNext(ts)

    so = ["stroke-opacity",1]
    for t,i in @taskList
      t.sp.hover(t.hoveron,t.hoverout)
      if t.task.Id isnt ts.task.Id and t.task.JobNo >= ts.task.JobNo
        [t.sp.animate(s1), t.text.animate(txt)]
        [r.bg.attr(so...),r.line.attr(so...)] for r,i in t.preRel
        [r.bg.attr(so...),r.line.attr(so...)] for r,i in t.nextRel

    tpre = ts.pre
    while tpre.length > 0
      tmp = []
      for rts,i in tpre
        [rts.sp.animate(s1), rts.text.animate(txt)]
        [r.bg.attr(so...),r.line.attr(so...)] for r,i in rts.nextRel
        if rts.pre.length > 0
          for r,j in rts.pre
            tmp.push(r)
      tpre = tmp

class TaskShape
  constructor: (@paper, @cx, @cy, @task, @color="#FF8C00", @r=20) ->
    
    @RelTaskId = (v.Id for k,v of @task.RelTasks)
    @pre=[]
    @preRel=[]
    @next=[]
    @nextRel=[]

    @draw()

  draw: ->
    @toolset = @paper.set()

    imgStyle = [@cx, @cy, 15, 15]
    @editImg=@paper.image("img/edit.png", imgStyle...)
    @deleteRelImg=@paper.image("img/delrel.png", imgStyle...)
    @deleteImg=@paper.image("img/delete.png", imgStyle...)
    @connImg=@paper.image("img/conn.png", imgStyle...)
    @connImg.toBack()

    @edit=@paper.circle(@cx, @cy , 14)
    @edit.click(mm = (e) ->
        @sp.flg = true
        @showTool()
        @.task.opt = "edit"
        Spine.trigger("addTaskRender", @.task)
      ,@)

    @deleteRel=@paper.circle(@cx, @cy, 14)
    @deleteRel.click(mm = (e) ->
        Spine.trigger("deleteTaskRelStart", @, e||window.event)
      ,@)

    @delete=@paper.circle(@cx, @cy, 14)
    @delete.click(mm = (e) ->
        e = e||window.event
        @.task.opt = "delete"
        Spine.trigger("deleteTask", @, e||window.event)
      ,@)

    @conn=@paper.circle(@cx, @cy, 14)
    @conn.refresh = =>
      if @conn.rel
        @paper.connection(@conn.rel)

    @conn.mousedown(mm = (e) ->
        Spine.trigger("connectTaskStart", @, e||window.event)
      ,@)

    @conn.mouseup(mm = (e) ->
        Spine.trigger("connectTaskFinish", @, e||window.event)
      ,@)

    @toolset.push(@editImg,@deleteRelImg,@deleteImg,@connImg,@edit,@deleteRel,@delete,@conn)
    @toolset.attr({fill: @color, stroke: @color, "fill-opacity": 0.1, "stroke-width": .5, cursor: "hand"})
    @toolset.hover(hh = ->
        @animate({"fill-opacity": 0.5}, 200)
      ,nn = ->
        @animate({"fill-opacity": 0.1}, 200))
    
    @toolset.hide()

    @sp=@paper.circle(@cx, @cy, @r)
    @sp.ts=@
    @sp.mousedown(@sp.md = ->
            @sp.flg = true
        ,@)

    @sp.mouseup(@showTool,@)
    @sp.hover(@hoveron,@hoverout)
    @sp.attr({fill: @color, stroke: @color, "fill-opacity": 0.2, "stroke-width": 1, cursor: "move"})
    @sp.refresh = ->
      if @ts.nextRel then for r,i in @ts.nextRel
        @paper.connection(r)

      if @ts.preRel then for r,i in @ts.preRel
        @paper.connection(r)
    @sp.drag(@move, @dragger, @up)

    @text = @paper.text(@cx, @cy, @task.Name)
    @text.toBack()
    @text.attr({fill: "#333", stroke: "none", "font-size": 10, "fill-opacity": 1, "stroke-width": 1, cursor: "move"})


  addNext: (taskShape) ->
    @next.push(taskShape)
    r=@paper.connection(@sp,taskShape.sp,@sp.attr('fill'),"#{@sp.attr('fill')}|1")
    @nextRel.push(r)

    taskShape.pre.push(@)
    taskShape.preRel.push(r)

  removeNext: (rel) ->
    @next = (t for t,i in @next when t isnt rel.tail)
    @nextRel = (r for r,i in @nextRel when r isnt rel)

    rel.tail.pre = (t for t,i in rel.tail.pre when t isnt rel.head)
    rel.tail.preRel = (r for r,i in rel.tail.preRel when r isnt rel)

    rel.bg.remove()
    rel.line.remove()
    rel = null

  remove: ->
    for r,j in @nextRel
      r.bg.remove()
      r.line.remove()
      r = null

    for t,j in @next
      t.pre = (p for p,i in t.pre when p isnt @)
      t.preRel = (r for r,i in t.preRel when r.bg.id isnt null)

    for r,j in @preRel
      r.bg.remove()
      r.line.remove()
      r = null

    for t,j in @pre
      t.next = (n for n,i in t.next when n isnt @)
      t.nextRel = (r for r,i in t.nextRel when r.bg.id isnt null)

    for key, value of @
      unless value is @paper
        value.remove?()
    

  showTool: ->
    return unless @sp.flg

    s = @sp.ts
    mc = Math.cos(45*Math.PI/180)
    ms = Math.sin(45*Math.PI/180)
    mc1 = Math.cos(90*Math.PI/180)
    ms1 = Math.sin(90*Math.PI/180)
    if @sp.isShowTool
      [x,y] = [@sp.attr("cx"), @sp.attr("cy")]
      s.toolset.animate({"x": x, "y": y, "cx": x, "cy": y}, 80, "backin",-> @.hide())
      @sp.isShowTool = false
    else
      s.editImg.animate({"x": s.editImg.ox + 50, "y": s.editImg.oy - 7.5}, 600, "elastic")
      s.deleteImg.animate({"x": s.deleteImg.ox + 50 * mc, "y": s.deleteImg.oy + 50 * ms - 7.5}, 600, "elastic")
      s.connImg.animate({"x": s.connImg.ox + 50 * mc, "y": s.connImg.oy - 50 * ms - 7.5}, 600, "elastic")
      s.deleteRelImg.animate({"x": s.deleteRelImg.ox + 50 * mc1, "y": s.deleteRelImg.oy - 50 * ms1 - 7.5}, 600, "elastic")
      s.edit.animate({"cx": s.edit.ox + 57, "cy": @sp.ts.edit.oy}, 600, "elastic")
      s.delete.animate({"cx": s.delete.ox + 60 * mc, "cy": s.delete.oy + 50 * ms}, 600, "elastic")
      s.deleteRel.animate({"cx": s.deleteRel.ox + 60 * mc1 + 7, "cy": s.deleteRel.oy - 50 * ms1}, 600, "elastic")
      s.conn.animate({"cx": s.conn.ox + 60 * mc, "cy": s.conn.oy - 50 * ms}, 600, "elastic")
      s.toolset.show()
      @sp.isShowTool = true

  hoveron: =>
    a = Raphael.animation({"stroke-width": 3, "fill-opacity": 0.7}, 200)

    @sp.animate(a)
    r.line.animate(a)  for r in @nextRel
    n.sp.animate(a)    for n in @next
    rp.line.animate(a) for rp in @preRel
    p.sp.animate(a)    for p in @pre
      
  hoverout: =>
    b = Raphael.animation({"stroke-width": 1,"fill-opacity": 0.2}, 200)
    @sp.animate(b)
    r.line.animate(b)  for r in @nextRel
    n.sp.animate(b)    for n in @next
    rp.line.animate(b) for rp in @preRel
    p.sp.animate(b)    for p in @pre
      
  dragger: ->
    [@ox, @oy]  = [@attr("cx"), @attr("cy")]
    @animate({"fill-opacity": .5}, 500) if @type isnt "text"

    for el in@ts.toolset
      if el.type is "image"
        [el.ox,el.oy] = [el.attr("x"), el.attr("y")]
      else
        [el.ox,el.oy]  = [el.attr("cx"), el.attr("cy")]
    [@ts.text.ox, @ts.text.oy] = [@attr("x"),@attr("y")]
    @ts.text.animate({"fill-opacity": .2}, 500) if @ts.text.type isnt "text"


  move: (dx, dy) ->
    @flg = false
    @attr([ cx:@ox + dx, cy:@oy + dy])

    for el in@ts.toolset
      el.attr([ cx:el.ox + dx, cy:el.oy + dy]) if el.attr("cx")
      el.attr([ x:el.ox + dx, y:el.oy + dy]) if el.attr("x")
    @ts.text.attr([x:@ox + dx, y:@oy + dy])
    @refresh()

  up: ->
    a = [{"fill-opacity": 0.2}, 500]
    @animate(a...) if @type isnt "text"
    @ts.text.animate(a...) if @ts.text.type isnt "text"

  connDragger: ->
    c = @currentTask
    [c.conn.ox, c.conn.oy]  = [c.conn.attr("cx"), c.conn.attr("cy")]
    c.conn.animate({"fill-opacity": .5}, 500)

    c.connImg.ox = c.connImg.attr("x")
    c.connImg.oy = c.connImg.attr("y")
    c.connImg.hide()
    c.editImg.hide()
    c.deleteRelImg.hide()
    c.deleteImg.hide()
    c.toolset.attr({"fill-opacity": 0, "stroke-width": 0})

  connMove: (dx, dy) ->
    flg = false
    c = @currentTask
    c.sp.flg = true
    c.conn.attr([ cx:c.conn.ox + dx, cy:c.conn.oy + dy])

    c.connImg.attr([ x:c.connImg.ox + dx, y:c.connImg.oy + dy])
    for t, i in @taskList when t.sp.attr("fill-opacity") isnt .05
      if t.sp.isPointInside(c.conn.attr("cx"),c.conn.attr("cy"))
        @relTask = t
        c.conn.animate({fill: t.sp.attr("fill"), "fill-opacity": 1, stroke: t.sp.attr("fill"), "stroke-width": 4},100)
        c.conn.rel.line.animate({"stroke": t.sp.attr("fill"), "stroke-width": 6},100)
        c.conn.rel.bg.animate({"stroke": t.sp.attr("fill"), "stroke-width": 6},100)
        flg = true

    unless flg
      @relTask = null
      c.conn.animate({fill: c.sp.attr("fill"), stroke: c.sp.attr("fill"), "stroke-width": 1}, 50)
      c.conn.rel.line.animate({"stroke": c.sp.attr("fill"), "stroke-width": 2}, 50)
      c.conn.rel.bg.animate({"stroke": c.sp.attr("fill"), "stroke-width": 2}, 50)

    c.conn.refresh()

  connUp: ->
    c = @currentTask
    c.conn.animate({fill: c.sp.attr("fill"), stroke: c.sp.attr("fill"), "fill-opacity": 0.2}, 500)
    
    c.toolset.show()
    c.toolset.attr({"fill-opacity": 0.1, "stroke-width": 0.5})

    c.showTool()
    @currentTask = null
    @relTask = null

module.exports = TaskManager
