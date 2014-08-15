Spine = require('spineify')
Schedule = require('models/schedule')
$       = Spine.$

class Navbar extends Spine.Controller
  className: 'hnavbar'
  
  #elements:

  events:
   "mouseenter h1":   "mouseover"
   "mouseleave h1":   "mouseout"

  constructor: ->
    super

  render: ->
    @html(require('views/navbar')())

  mouseover: (e) ->
    $(e.target).stop().animate({backgroundColor:'#777'},"fast")


  mouseout: (e) ->
    $(e.target).stop().animate({backgroundColor:'#333'},"fast")


module.exports = Navbar
