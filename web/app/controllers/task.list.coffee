Spine = require('spineify')
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

  events:
    "click .tclose"        :  "hideTask"
    "click .tparam"        :  "editParam"
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
    @setpp = @paper.set()
    @isRefresh = true
    @isMove = false
    @jobList = []
    @taskList = []
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
        t= new TaskShape(@paper,left,top,tk,@color[i],25)
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
      tk.bind("ajaxSuccess",@addTaskAndRefresh)
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
      tk.save({url:"/schedules/#{@item.Id}/jobs/#{tk.jobid}/tasks/#{tk.Id}"})
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

  addTaskAndRefresh: (task, status, xhr) =># {{{
    s = Raphael.animation({"fill-opacity": .3, "stroke-opacity": .3, "stroke-width": 6}, 1500, -> @.animate(e))
    e = Raphael.animation({"fill-opacity": .01, "stroke-opacity": .01, "stroke-width": 1}, 1500, -> @.animate(s))
    if xhr is "success"
      Spine.Module.extend.call(task, Task)

      ci = i for j,i in @item.Jobs when j.Id is task.JobId

      t= new TaskShape(@paper,150,0,task,Style.getRgbColor()[ci],25)

      t.sp.animate({"cx": 150, "cy": ci*100+80}, 2000, "elastic")
      t.text.animate({"x": 150, "y": ci*100+80}, 2000, "elastic")

      #t.sp.animate(s)
      #t.text.animate(s)
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
    @el.css("opacity", 0.4)
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
    @editImg=@paper.image("img/edit.png", @cx, @cy, 15, 15)
    @deleteImg=@paper.image("img/delete.png", @cx , @cy, 15, 15)
    @connImg=@paper.image("img/conn.png", @cx, @cy, 15, 15)

    @edit=@paper.circle(@cx, @cy , 14)
    @edit.click(@showEdit,@)
    @delete=@paper.circle(@cx, @cy, 14)
    @conn=@paper.circle(@cx, @cy, 14)
    @toolset.push(@editImg,@deleteImg,@connImg,@edit,@delete,@conn)
    @toolset.attr({fill: @color, stroke: @color, "fill-opacity": 0.1, "stroke-width": .5, cursor: "hand"})
    @toolset.hover(@hlight,@nlight)
    @toolset.hide()

    @sp=@paper.circle(@cx, @cy, @r)
    @sp.ts=@
    @sp.dblclick(@showTool,@)
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

  hlight: -># {{{
    a = Raphael.animation({"fill-opacity": 0.5}, 200)
    @animate(a)
  # }}}

  nlight: ->
    a = Raphael.animation({"fill-opacity": 0.1}, 200)
    @animate(a)

  addNext: (taskShape) ->
    @next.push(taskShape)
    r=@paper.connection(@sp,taskShape.sp,@sp.attr('fill'),"#{@sp.attr('fill')}|1")
    @nextRel.push(r)

    taskShape.pre.push(@)
    taskShape.preRel.push(r)

  showEdit: (e) ->
    e = e||window.event
    @showTool()
    @.task.opt = "edit"
    Spine.trigger("addTaskRender", @.task)
    e


  showTool: ->
    if @sp.isShowTool
      @sp.ts.toolset.animate({"x": @sp.ox, "y": @sp.oy, "cx": @sp.ox, "cy": @sp.oy}, 200, "backin",-> @.hide())
      @sp.isShowTool = false
    else
      @sp.ts.editImg.animate({"x": @sp.ts.editImg.ox + 50, "y": @sp.ts.editImg.oy - 7.5}, 600, "elastic")
      @sp.ts.deleteImg.animate({"x": @sp.ts.deleteImg.ox + 50 * Math.cos(45*Math.PI/180), "y": @sp.ts.deleteImg.oy + 50 * Math.sin(45*Math.PI/180) - 7.5}, 600, "elastic")
      @sp.ts.connImg.animate({"x": @sp.ts.connImg.ox + 50 * Math.cos(45*Math.PI/180), "y": @sp.ts.connImg.oy - 50 * Math.sin(45*Math.PI/180) - 7.5}, 600, "elastic")
      @sp.ts.edit.animate({"cx": @sp.ts.edit.ox + 57, "cy": @sp.ts.edit.oy}, 600, "elastic")
      @sp.ts.delete.animate({"cx": @sp.ts.delete.ox + 60 * Math.cos(45*Math.PI/180), "cy": @sp.ts.delete.oy + 50 * Math.sin(45*Math.PI/180)}, 600, "elastic")
      @sp.ts.conn.animate({"cx": @sp.ts.conn.ox + 60 * Math.cos(45*Math.PI/180), "cy": @sp.ts.conn.oy - 50 * Math.sin(45*Math.PI/180)}, 600, "elastic")
      @sp.ts.toolset.show()
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
        el.ox = el.attr("x")
        el.oy = el.attr("y")
      else
        el.ox = el.attr("cx")
        el.oy = el.attr("cy")

    [@ts.text.ox, @ts.text.oy] = [@attr("x"),@attr("y")]
    @ts.text.animate({"fill-opacity": .2}, 500) if @ts.text.type isnt "text"

  move: (dx, dy) ->
    @attr([ cx:@ox + dx, cy:@oy + dy])

    for el in@ts.toolset
      el.attr([ cx:el.ox + dx, cy:el.oy + dy]) if el.attr("cx")
      el.attr([ x:el.ox + dx, y:el.oy + dy]) if el.attr("x")

    @ts.text.attr([x:@ox + dx, y:@oy + dy])
    @refresh()

  up: ->
    @animate({"fill-opacity": 0.2}, 500) if @type isnt "text"
    @ts.text.animate({"fill-opacity": 0.2}, 500) if @ts.text.type isnt "text"

module.exports = TaskManager
