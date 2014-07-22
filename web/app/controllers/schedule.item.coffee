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
   "mouseenter":            "mouseover"
   "mouseleave":             "mouseout"

   "mouseenter .sstart":    "sstartmouseover"
   "mouseleave .sstart":     "sstartmouseout"

   "mouseenter .jobcnt":    "jobcntmouseover"
   "mouseleave .jobcnt":     "jobcntmouseout"

   "mouseenter .nextstart": "nextstartmouseover"
   "mouseleave .nextstart":  "nextstartmouseout"

   "click .cyc":       "showcyc"
   "click .sname":       "showschedule"
   "click .sstart":       "ck"

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
    require('views/schedule-show')(items)

  remove: ->

  showcyc: ->
    alert('！')

  showschedule: (e)->
    @navigate('/schedules', @item.Id)

  ck: (e) ->
    if e.target.className.indexOf("glyphicon-plus")>=0
      alert(e.target.className)
      e.stopPropagation()

  mouseover: (e)->
    @panel.stop().animate({boxShadow:'0 0 20px #777'},"fast")
    @timout=window.setTimeout( =>
        @pbmask.fadeOut(400)
        @sname.stop().animate({color:"#333"},800)
        @header.stop().animate({backgroundColor:"#E0E0E0"},"fast")

        @srun.stop().animate({backgroundColor:"#999"},800)
        @sdelete.stop().animate({backgroundColor:"#999"},800)
        @scopy.stop().animate({backgroundColor:"#999"},800)
      ,800)

  mouseout: (e)->
    window.clearTimeout(@timout)
    @pbmask.fadeIn(200)
    @sname.stop().animate({color:"#f5f5f5"},"fast")
    @header.stop().animate({backgroundColor:"#f5f5f5"},"fast")
    @panel.stop().animate({boxShadow:''},"fast")

    @srun.stop().animate({backgroundColor:"#f5f5f5"},"fast")
    @scopy.stop().animate({backgroundColor:"#f5f5f5"},"fast")
    @sdelete.stop().animate({backgroundColor:"#f5f5f5"},"fast")

  sstartmouseover: (e)->
      @sstart.css("background-color","#E0E0E0")
      @addstart.stop().animate({backgroundColor:"#999"},1)
  sstartmouseout: (e)->
      @sstart.css("background-color","#f5f5f5")
      @addstart.stop().animate({backgroundColor:"#f5f5f5"},10)

  jobcntmouseover: (e)->
      @jobcnt.css("background-color","#E0E0E0")
  jobcntmouseout: (e)->
      @jobcnt.css("background-color","#f5f5f5")

  nextstartmouseover: (e)->
      @nextstart.css("background-color","#E0E0E0")
  nextstartmouseout: (e)->
      @nextstart.css("background-color","#f5f5f5")

module.exports = ScheduleItem
