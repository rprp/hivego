Spine = require('spineify')
Schedule = require('models/schedule')
$       = Spine.$

class ScheduleItem extends Spine.Controller
  className: 'scheduleitem'
  
  elements:
    ".panel":          "panel"
    ".panel-heading":  "header"
    ".panel-body":     "body"
    ".pbmask":         "pbmask"
    ".panel-footer":   "footer"

    ".sname":          "sname"
    ".cyc":            "cyc"
    ".sstart":          "sstart"
    ".startlist":      "startlist"
    ".jobcnt":         "jobcnt"
    ".nextstart":      "nextstart"
    ".scopy":          "scopy"
    ".sdelete":        "sdelete"
    ".slog":          "slog"
    ".srun":          "srun"
    ".addstart":      "addstart"
    
  events:
   "click .cyc":    "showcyc"
   "click .sname":  "showschedule"
   "click .sstart": "ck"
   "click .sdelete": "deleteSchedule"

   "mouseenter": @mouseover = (e)->
        @panel.stop().animate({boxShadow:'0 0 20px #777'},"fast")
        @cyc.stop().animate({opacity: 1},200)
        @timout=window.setTimeout( =>
            @pbmask.fadeOut(400)
            @body.stop().animate({opacity: 1},800)
            @footer.stop().animate({opacity: 1},800)
            @sname.stop().animate({color:"#333", opacity: 1},800)
            @header.stop().animate({backgroundColor:"#E0E0E0", opactiy: 1},"fast")

            @srun.stop().animate({backgroundColor:"#999", opactiy: 1},800)
            @sdelete.stop().animate({backgroundColor:"#999", opactiy: 1},800)
            @scopy.stop().animate({backgroundColor:"#999", opactiy: 1},800)
          ,800)

   "mouseleave": @mouseout = (e)->
        window.clearTimeout(@timout)
        @pbmask.fadeIn(200)
        @cyc.stop().animate({opacity: 0},200)
        @body.stop().animate({opacity: 0},800)
        @footer.stop().animate({opacity: 0},800)
        @sname.stop().animate({color:"transparent"},"fast")
        @header.stop().animate({backgroundColor:"transparent"},"fast")
        @panel.stop().animate({boxShadow:''},"fast")

        @srun.stop().animate({backgroundColor:"transparent"},"fast")
        @scopy.stop().animate({backgroundColor:"transparent"},"fast")
        @sdelete.stop().animate({backgroundColor:"transparent"},"fast")

   "mouseenter .sstart":    @sstartmouseover = (e)->
          @sstart.css("background-color","#E0E0E0")
          @addstart.stop().animate({backgroundColor:"#999"},1)
   "mouseleave .sstart":     @sstartmouseout = (e)->
          @sstart.css("background-color","transparent")
          @addstart.stop().animate({backgroundColor:"transparent"},10)

   "mouseenter .jobcnt": @jobcntmouseover = (e)->
          @jobcnt.css("background-color","#E0E0E0")
   "mouseleave .jobcnt": @jobcntmouseout = (e)->
          @jobcnt.css("background-color","transparent")

   "mouseenter .nextstart": @nextstartmouseover = (e)->
          @nextstart.css("background-color","#E0E0E0")
   "mouseleave .nextstart": @nextstartmouseout = (e)->
          @nextstart.css("background-color","transparent")

  constructor: ->
    super
    @el.addClass('col-sm-3')
    throw "@item required" unless @item
    @item.bind("update", @render)
    @item.bind("destroy", @remove)

  render: (item) =>
    @item = item if item
    @html(@template(@item))
    @
    
  template: (items) ->
    require('views/main-list')(items)

  remove: ->

  showcyc: ->
    alert('！')

  showschedule: (e)->
    @navigate('/schedules', @item.Id)

  ck: (e) ->
    if e.target.className.indexOf("glyphicon-plus")>=0
      alert(e.target.className)
      e.stopPropagation()

  deleteSchedule: (e) ->
    s = Schedule.find(@item.Id)
    s.bind("refresh",ScheduleList.addAll)
    @el.remove()
    s.destroy()

class ScheduleList extends Spine.Controller
  className: 'schedulelist'

  constructor: ->
    super
    Schedule.bind("create",  @addOne)
    Schedule.bind("refresh", @addAll)
    @html(require('views/main')())

    @active @change

  change: (params) =>
    @addAll()

  addOne: (it) =>
    view = new ScheduleItem(item: it)
    $('#smain').append(view.render().el)

    view.pbmask.css("position","absolute")
    view.pbmask.css("z-index","100")
    view.pbmask.css("width",view.sstart.css("width"))
    view.pbmask.css("height",view.body.css("height"))

  addAll: =>
    $('.scheduleitem').remove()
    Schedule.comparator = (a, b) ->
      return a.Id - b.Id
    Schedule.sort()
    Schedule.each(@addOne)

module.exports = ScheduleItem
module.exports = ScheduleList
