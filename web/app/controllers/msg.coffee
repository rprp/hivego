Spine = require('spineify')
Schedule = require('models/schedule')
ScheduleManager = require('controllers/schedule.info')
$       = Spine.$

class Msg extends Spine.Controller
  className: 'msg'
  
  elements:
    "#msg": "msg"
    ".alert": "alert"

  events:
    "click": (e) -> @alert.css('display','none')

  constructor: ->
    super
    Spine.bind("msg", @show)

  render: ->
    @html(require('views/msg')())

  show: (status,msg) =>
    if status >= 400
      @alert.removeClass("alert-success")
      @alert.addClass("alert-danger")
      @alert.stop().animate({opacity: 0},0)
      @alert.css('display','block')
      @alert.stop().animate({opacity: 1},800)
      @msg.text("错误：#{msg}")
      @timout=window.setTimeout( =>
          @alert.stop().animate({opacity: 0},1000)
        ,20000)

module.exports = Msg
