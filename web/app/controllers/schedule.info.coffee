Spine = require('spineify')
Raphael = require('raphaelify')
Eve = require('eve')
Schedule = require('models/schedule')
$       = Spine.$

class ScheduleManager
  constructor: (@paper, @color, @item, @width) ->
    @font = "Helvetica, Tahoma, Arial, STXihei, '华文细黑', Heiti, '黑体', 'Microsoft YaHei', '微软雅黑', SimSun, '宋体', sans-serif"
    @titlerectStyle = {fill: "#31708f", stroke: "#31708f", "fill-opacity": 0.05, "stroke-width": 0, cursor: "hand"}
    @fontStyle = {fill: "#333", "font-family":@font, "text-anchor": "start", stroke: "none", "font-size": 14, "fill-opacity": 1, "stroke-width": 1}
    @height = 0

    @paper.setStart()

    [top,left] = [20, 10]
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
    @startSecond = @paper.text(left, top, "启动时间：").attr(@fontStyle)

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
    #@titlerect.click(@showJob,@)

    @set = @paper.setFinish()
    @height = top 

  hoveron: ->
    a = Raphael.animation({"fill-opacity": 0.4}, 200)
    @.animate(a)
      
  hoverout: ->
    b = Raphael.animation({"fill-opacity": 0.1}, 200)
    @.animate(b)

module.exports = ScheduleManager
