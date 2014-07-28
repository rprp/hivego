Spine = require('spineify')
Raphael = require('raphaelify')
Eve = require('eve')
Schedule = require('models/schedule')
$       = Spine.$

class ScheduleManager extends Spine.Controller
  elements:
    ".close":  "close"
    ".addScheduleHead":  "scheduleHead"
    ".cyclbl": "cycGroup"
    ".startList":"startList"
    ".start":"start"
    ".startInput":"startInput"
    "#scheduleName":"scheduleName"
    "#scheduleDesc":"scheduleDesc"

  events:
    "click .close": "postSchedule"
    "click .cyclbl": "setCyc"
    "click .addStart": "addStart"
    "click .delStart": "delStart"
    "click .start":"editStart"
    "keypress .startInput":"setStartSecond"
    "keypress .addSchedule":  "addSchedule"

    "mouseenter .list-group-item":  "showDelStart"
    "mouseleave .list-group-item":  "hideDelStart"

    "mousedown .addScheduleHead": "setMoveFlg"
    "mouseup .addScheduleHead": "clearMoveFlg"
    "mousemove .addScheduleHead": "movePanel"

  constructor: (@paper, @color, @item, @width) -># {{{
    super
    @isMove = false
    @font = "Helvetica, Tahoma, Arial, STXihei, '华文细黑', Heiti, '黑体', 'Microsoft YaHei', '微软雅黑', SimSun, '宋体', sans-serif"
    @titlerectStyle = {fill: "#98FB98", stroke: "#98FB98", "fill-opacity": 0.05, "stroke-width": 0, cursor: "hand"}
    @fontStyle = {fill: "#333", "font-family":@font, "text-anchor": "start", stroke: "none", "font-size": 14, "fill-opacity": 1, "stroke-width": 1}
    @height = 0

    @isRefresh = true

    @refreshSchedule(20, 10)# }}}

  refreshSchedule: (top, left) =># {{{
    return [top,left] unless @isRefresh

    @st.pop().remove() while @st?.length

    @paper.setStart()
    [top,left] = [top + (@item.Name.length//7) * 20, left]
    #标题，调度名称，每行超过7个字符后要换行
    @title = @paper.text(left, top, @item.SplitName(7)).attr(@fontStyle)
    @title.attr("font-size", 22)
    [top,left] = [top + 30 + (@item.Name.length//7) * 20, left]
    
    #调度周期
    @cyc = @paper.text(left, top, "调度周期：#{@item.GetCyc()}").attr(@fontStyle)

    #调度时间
    gs=@item.GetSecond()
    [top,left] = [top+30, left]
    @start = @paper.text(left, top, "启动时间：").attr(@fontStyle)

    @startSecondList = []
    for ss in gs
      [top,left] = [top+30, left]
      @startSecondList.push(@paper.text(left+20, top, "#{ss}").attr(@fontStyle))

    #任务数量
    [top,left] = [top+30, left]
    @taskCnt = @paper.text(left, top, "任务数量：#{@item.TaskCnt}").attr(@fontStyle)

    #下次执行时间
    [top,left] = [top+30, left]
    @nextStart = @paper.text(left, top, "下次执行：#{@item.GetNextStart()}").attr(@fontStyle)
    
    #当前状态
    #所有者

    [top,left] = [top+30, left]
    @betweenline = @paper.path("M #{left},#{top}L #{@width-30},#{top}").attr({stroke: "#A0522D", "stroke-width": 2, "stroke-opacity": 0.2})

    @titlerect = @paper.rect(left,0,190,top-10,3).attr(@titlerectStyle)
    @titlerect.hover(@hoveron,@hoverout)
    @titlerect.click(@showSchedule,@)

    @st = @paper.setFinish()

    @isRefresh = false
    @height = top# }}}

  showDelStart: (e) ->
    $(e.target).children(".delStart").css("display","")

  hideDelStart: (e) ->
    $(e.target).children(".delStart").css("display","none")

  delStart: (e) ->
    $(e.target).parent("li").remove()

  addStart: ->
    @startList.append(require('views/schedule-start')(@item.GetDefaultSecond()))
    $(".startInput").focus()

  setStartSecond: (e) ->
    e = e||window.event
    if e.keyCode in [13,10]
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

  editStart: (e) ->
    $(e.target).siblings().not(".delStart").css("display","")
    $(e.target).siblings().focus()
    $(e.target).css("display","none")

  addSchedule: (e) ->
    e = e||window.event
    if e.ctrlKey and e.keyCode in [13,10]
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
      @item.bind("ajaxSuccess",@scheduleRefresh)
      @item.save()

  scheduleRefresh:  (data, status, xhr) =>
    if xhr is "success"
      id = @item.Id
      Schedule.fetch({Id:id})
      @item = Schedule.find(id)
      @isRefresh = true


  setCyc: (e) -># {{{
    @cycGroup.removeClass("label-success")
    @cycGroup.addClass("label-default")

    $(e.target).removeClass("label-default")
    $(e.target).addClass("label-success")

    @item.SetCyc($(e.target).text())
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

  hoveron: -># {{{
    a = Raphael.animation({"fill-opacity": 0.4}, 200)
    @.animate(a)
      # }}}

  hoverout: -># {{{
    b = Raphael.animation({"fill-opacity": 0.1}, 200)
    @.animate(b)
# }}}

  render: (x, y, schedule) =># {{{
    @html(require('views/schedule')(schedule))

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
    
  postSchedule: -># {{{
    @el.css("display","none")
# }}}

module.exports = ScheduleManager
