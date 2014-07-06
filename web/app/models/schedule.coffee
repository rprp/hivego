Moment = require('momentify')
Spine = require('spine')

class Schedule extends Spine.Model
  @configure 'Schedule', 'Id', 'Name', 'TaskCnt', 'Count', 'Cyc', 'StartMonth', 'StartSecond', 'NextStart', 'TimeOut', 'Desc', 'CreateTime', 'CreateUserId', 'ModifyTime', 'ModifyUserId'

  @extend Spine.Model.Ajax
  
  constructor: ->
    super
    Moment.lang('zh-cn')

  GetNextStart: ->
    Moment(@NextStart).calendar()

  GetCycStyle: ->
    switch @Cyc
      when "ss" then "label-default"
      when "mi" then "label-warning"
      when "h" then "label-info"
      when "d" then "label-primary"
      when "w" then "label-default"
      when "m" then "label-success"
      when "q" then "label-default"
      when "y" then "label-danger"
      else "无"

  GetCyc: ->
    switch @Cyc
      when "ss" then "秒"
      when "mi" then "分"
      when "h" then "时"
      when "d" then "日"
      when "w" then "周"
      when "m" then "月"
      when "q" then "季"
      when "y" then "年"
      else "无"

  GetSecond: ->
    for t,i in @StartSecond
      #转换成秒
      t = t/1000/1000/1000

      #输出格式化时间
      switch
		#秒
        when t<60
          "#{@addStartPrefix('ss')}#{t}秒"
		#分钟
        when 60<=t<60*60
          "#{@addStartPrefix('mi')}#{@getMinuets(t)}"
		#小时
        when 60*60<=t<60*60*24
          "#{@addStartPrefix('h')}#{@getHour(t)}"
        #日
        when 60*60*24<=t<60*60*24*31
          if @StartMonth[i] > 0
            "#{@addStartPrefix('m')}#{@StartMonth[i]+1}月#{@getDay(t)}"
          else
            "#{@addStartPrefix('d')}#{@getDay(t)}"

        else "未知"
  
  #增加时间前缀
  #懒的弄了瞎写写吧，哪天心情好再改。
  addStartPrefix: (cyc) ->
    switch cyc
      when "ss"
        switch @Cyc
          when "h" then "每小时0分"
          when "d" then "每日0点0分"
          when "w" then "每周1 0点0分"
          when "m" then "每月1日0点0分"
          when "q" then "每季度1日0点0分"
          when "y" then "每年1月1日0点0分"
          else ""
      when "mi"
        switch @Cyc
          when "d" then "每日0点"
          when "w" then "每周1 0点"
          when "m" then "每月1日0点"
          when "q" then "每季度1日0点"
          when "y" then "每年1月1日0点"
          else ""
      when "h"
        switch @Cyc
          when "m" then "每月1日"
          when "q" then "每季度1日"
          when "y" then "每年1月1日"
          else ""
      when "d"
        switch @Cyc
          when "w" then "每周"
          when "m" then "每月"
          when "q" then "每季度"
          when "y" then "每年1月"
          else ""
      when "m"
        switch @Cyc
          when "y" then "每年"
          else ""
      else ""


  getMinuets: (tm) ->
    if (tm%60) > 0
      "#{tm//60}分#{tm%60}秒"
    else
      "#{tm/60}分"
     
  getHour: (tm) ->
    if tm%(60*60) > 0
      "#{tm//(60*60)}点#{@getMinuets(tm%(60*60))}"
    else
      "#{tm/60/60}点"

  getDay: (tm) ->
    if tm%(60*60*24) > 0
      "#{tm//(60*60*24)}日#{@getHour(tm%(60*60*24))}"
    else
      "#{tm/60/60/24}日"

module.exports = Schedule


