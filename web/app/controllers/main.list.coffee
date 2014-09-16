Spine = require('spineify')
Schedule = require('models/schedule')
$       = Spine.$

class ScheduleItem extends Spine.Controller
  className: 'scheduleitem col-sm-3'
  
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
   "click .sdelete": "deleteSchedule"
   "click .sname": (e)-> @navigate('/schedules', @item.Id)
   "click .pbmask": (e)-> @navigate('/schedules', @item.Id)

   "mouseenter": @mouseover = (e)->
        @panel.stop().animate({boxShadow:'0 0 20px #777'},"fast")
        @cyc.stop().animate({opacity: 1},200)
        @timout=window.setTimeout( =>
            @pbmask.fadeOut(400)
            @body.stop().animate({opacity: 1},800)
            @footer.stop().animate({opacity: 1},800)
            @sname.stop().animate({color:"#333", opacity: 1},800)
            @header.stop().animate({backgroundColor:"rgba(196, 187, 142, 1)", opactiy: 1},"fast")

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

   "mouseenter .sstart": (e)->
          @sstart.css("background-color","rgba(196, 187, 142, 1)")
          @addstart.stop().animate({backgroundColor:"#999"},1)
   "mouseleave .sstart": (e)->
          @sstart.css("background-color","transparent")
          @addstart.stop().animate({backgroundColor:"transparent"},10)

   "mouseenter .jobcnt": (e)-> @jobcnt.css("background-color","rgba(196, 187, 142, 1)")
   "mouseleave .jobcnt": (e)-> @jobcnt.css("background-color","transparent")

   "mouseenter .nextstart": (e)-> @nextstart.css("background-color","rgba(196, 187, 142, 1)")
   "mouseleave .nextstart": (e)-> @nextstart.css("background-color","transparent")

  constructor: ->
    super
    throw "@item required" unless @item
    @item.bind("update", @render)

  render: (item) =>
    @item = item if item
    @html(require('views/main-list')(@item))
    @
    
  deleteSchedule: (e) ->
    s = Schedule.find(@item.Id)
    s.bind("refresh", MainList.addAll)
    @el.remove()
    s.destroy()

class MainList extends Spine.Controller
  className: 'mainlist container'

  elements:
    "#row": "row"

  constructor: ->
    super
    Schedule.bind("create",  @addOne)
    Schedule.bind("refresh", @addAll)
    @html(require('views/main')())

    @active (params) => @addAll()

  addOne: (it) =>
    view = new ScheduleItem(item: it)
    @row.append(view.render().el)

  addAll: =>
    $('.scheduleitem').remove()
    Schedule.comparator = (a, b) ->
      return a.Id - b.Id
    Schedule.sort()
    Schedule.each(@addOne)

module.exports = MainList
