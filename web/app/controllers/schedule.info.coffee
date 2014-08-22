Spine = require('spineify')
Events  = Spine.Events
Module  = Spine.Module
Raphael = require('raphaelify')
Style = require('controllers/style')
Eve = require('eve')
Schedule = require('models/schedule')
$       = Spine.$

class Form extends Spine.Controller
  elements:
    ".cyclbl": "cycGroup"
    ".startList":"startList"
    ".start":"start"
    ".startInput":"startInput"
    "#scheduleName":"scheduleName"
    "#scheduleDesc":"scheduleDesc"

  events:
    "click #submitSchedule": "addSchedule"
    "blur .startInput":  "setStartVal"

    "click .close": -> @el.css("display","none")

    "click .cyclbl": @setCyc = (e) ->
            @cycGroup.removeClass("label-success")
            @cycGroup.addClass("label-default")
            $(e.target).removeClass("label-default")
            $(e.target).addClass("label-success")
            @item.SetCyc($(e.target).text())

    "click .addStart": @appendStart = ->
            @startList.append(require('views/schedule-start')(@item.GetDefaultSecond()))
            $(".startInput").focus()
    "click .delStart": (e) -> $(e.target).parent("li").remove()

    "click .start": @editStart = (e) ->
            $(e.target).siblings().not(".delStart").css("display","")
            $(e.target).siblings().focus()
            $(e.target).css("display","none")

    "keypress .startInput":  @startKeypress = (e) ->
            e = e||window.event
            if e.keyCode in [13,10]
              @setStartVal()
  
    "keypress .addSchedule": @keypress = (e) ->
            e = e||window.event
            if e.ctrlKey and e.keyCode in [13,10]
              @addSchedule(e)

    "mouseenter .list-group-item": @showDelStart = (e) ->
            $(e.target).children(".delStart").css("display","")
    "mouseleave .list-group-item": @hideDelStart = (e) ->
            $(e.target).children(".delStart").css("display","none")
  

    "mousedown .addScheduleHead": @setMoveFlg = (e) ->
            @isMove = true
            @preLeft = e.clientX
            @preTop = e.clientY
  
    "mouseup .addScheduleHead": @clearMoveFlg = (e) ->
            @isMove = false
            @el.css("opacity", 1)
  
    "mousemove .addScheduleHead": @movePanel = (e) ->
            return unless @isMove
            e = e||window.event

            dx = (e.clientX - @preLeft) + parseInt(@el.css("left"))
            dy = (e.clientY - @preTop) + parseInt(@el.css("top"))
            @el.css("left", dx)
            @el.css("top", dy)
            @el.css("opacity", 0.4)

            @preLeft = e.clientX
            @preTop = e.clientY


  constructor: (@c, @item) -># {{{
    super
  # }}}

  setStartVal: (e) -># {{{
    e = e||window.event
    $(e.target).css("display","none")
    $(e.target).siblings().not(".delStart").css("display","")
    $(e.target).siblings().not(".delStart").text(" #{$(e.target).val()}")
    [m,t] = @item.ParseSecond($(e.target).val())

    if t is -1
      $(e.target).siblings().not(".delStart").addClass("alert-danger")
    else
      $(e.target).siblings().not(".delStart").removeClass("alert-danger")
      $(e.target).siblings().filter(".startSecond").val(t)
      $(e.target).siblings().filter(".startMonth").val(m)
  # }}}

  addSchedule: (e) -># {{{
    @el.css("display","none")
    @item.Name = @scheduleName.val()
    @item.Desc = @scheduleDesc.val()
    @item.StartMonth = []
    @item.StartSecond = []
    for li,i in @startList.children("li")
      ss = $(li).children(".startSecond").val()
      sm = $(li).children(".startMonth").val()
      if ss isnt -1 and ss isnt ""
        @item.StartMonth.push(parseInt(sm))
        @item.StartSecond.push(parseInt(ss))

    if @item.Id is -1
      @item.create()
    else
      @item.bind("ajaxSuccess",@scheduleRefresh)
      @item.save()
  # }}}

  scheduleRefresh:  (data, status, xhr) =># {{{
    if xhr is "success"
      id = @item.Id
      Schedule.fetch({Id:id})
      @item = Schedule.find(id)
  # }}}


  render: (x, y, schedule) =># {{{
    @html(require('views/schedule')(schedule))
    window.setTimeout( =>
            @scheduleName.focus()
      ,500)

    cs = c for c in @cycGroup when c.textContent is @item.GetCyc()
    $(cs).removeClass("label-default")
    $(cs).addClass("label-success")

    @el.css("display","block")
    @el.css("position","absolute")
    @el.css("left", x-400)
    @el.css("top", y-50)
  # }}}

  showSchedule: (e) -># {{{
    e = e||window.event
    Spine.trigger("editScheduleRender",e.clientX,e.clientY,@.item)
    e
  # }}}
    
  

class Shape extends Spine.Controller

  constructor: (@paper, @color, @item, @width) -># {{{
    super
    @isMove = false
    @height = 0
    @refreshSchedule(20, 10)
  # }}}

  refreshSchedule: (top, left, isRefresh = true) =># {{{
    return [top,left] unless isRefresh
    @st.pop().remove() while @st?.length

    @paper.setStart()
    [top,left] = [top + (@item.Name.length//7) * 20, left]
    #标题，调度名称，每行超过7个字符后要换行
    @title = @paper.text(left, top, @item.SplitName(7)).attr(Style.fontStyle)
    @title.attr("font-size", 22)
    [top,left] = [top + 30 + (@item.Name.length//7) * 20, left]
    
    #调度周期
    @cyc = @paper.text(left, top, "调度周期：#{@item.GetCyc()}").attr(Style.fontStyle)

    #调度时间
    gs=@item.GetSecond()
    [top,left] = [top+30, left]
    @start = @paper.text(left, top, "启动时间：").attr(Style.fontStyle)

    @startSecondList = []
    for ss in gs
      [top,left] = [top+30, left]
      @startSecondList.push(@paper.text(left+20, top, "#{ss}").attr(Style.fontStyle))

    #任务数量
    [top,left] = [top+30, left]
    @taskCnt = @paper.text(left, top, "任务数量：#{@item.TaskCnt}").attr(Style.fontStyle)

    #下次执行时间
    [top,left] = [top+30, left]
    @nextStart = @paper.text(left, top, "下次执行：#{@item.GetNextStart()}").attr(Style.fontStyle)
    
    [top,left] = [top+30, left]
    @betweenline = @paper.path("M #{left},#{top}L #{@width-30},#{top}").attr({stroke: "#A0522D", "stroke-width": 2, "stroke-opacity": 0.2})

    @titlerect = @paper.rect(left,0,190,top-10,18).attr(Style.titlerectStyle)
    @titlerect.hover(@hoveron,@hoverout)

    @st = @paper.setFinish()

    @height = top
  # }}}

  hoveron: -># {{{
    a = Raphael.animation({"fill-opacity": 0.1}, 200)
    @.animate(a)
    # }}}

  hoverout: -># {{{
    b = Raphael.animation({"fill-opacity": 0.01}, 200)
    @.animate(b)
  # }}}

ScheduleManager = {}
ScheduleManager.Form = Form
ScheduleManager.Shape = Shape
module.exports = ScheduleManager
