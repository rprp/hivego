Moment = require('momentify')
Spine = require('spine')

class Schedule extends Spine.Model
  @configure 'Schedule', 'Id', 'Name', 'TaskCnt', 'Job', 'Count', 'Cyc', 'StartMonth', 'StartSecond', 'NextStart', 'TimeOut', 'Desc', 'CreateTime', 'CreateUserId', 'ModifyTime', 'ModifyUserId'

  @extend Spine.Model.Ajax
  
  constructor: ->
    Moment.lang('zh-cn')
    super

  GetNextStart: ->
    Moment(@NextStart).calendar()

  GetCycStyle: ->
    switch @Cyc
      when "ss" then "label-default"
      when "mi" then "label-warning"
      when "h" then "label-info"
      when "d" then "label-success"
      when "w" then "label-default"
      when "m" then "label-primary"
      when "y" then "label-danger"
      else "无"

  GetCyc: ->
    switch @Cyc
      when "ss" then "Second"
      when "mi" then "Minute"
      when "h" then "Hour"
      when "d" then "Day"
      when "w" then "Week"
      when "m" then "Month"
      when "y" then "Year"
      else "无"

  SplitName: (len) =>
    s = @Name.split("")
    sname = ""
    for si,i in s
      sname += si
      if i>0 and i%%len is 0
        sname += "\n"
    sname

  GetSecond: ->
    for t,i in @StartSecond

      startMonth = if @StartMonth[i] is 0 then 1 else @StartMonth[i]
      
      #转换成秒
      t = t/1000/1000/1000
      s =
          "y": Moment.duration(t,'seconds').years()
          "m": Moment.duration(t,'seconds').months()
          "d": Moment.duration(t,'seconds').days()
          "h": Moment.duration(t,'seconds').hours()
          "mi":Moment.duration(t,'seconds').minutes()
          "ss":Moment.duration(t,'seconds').seconds()
      
      cyc=""
      if s['d']>0
          cyc='d'
      else if s['h']>0
          cyc='h'
      else if s['mi']>0
          cyc='mi'
      else if s['ss']>0
          cyc='ss'

      sc=""
      switch cyc
        when "ss"
          switch @Cyc
            when "mi" then sc="每分"
            when "h" then sc="每小时0分"
            when "d" then sc="每日0点0分"
            when "w" then sc="每周1 0点0分"
            when "m" then sc="每月1日0点0分"
            when "y" 
                if startMonth >1
                    sc="每年#{startMonth}月1日0点0分"
                else
                    sc="每年1月1日0点0分"
            else ""
          sc="#{sc}#{s['ss']}秒"
        when "mi"
          switch @Cyc
            when "h" then sc="每小时"
            when "d" then sc="每日0点"
            when "w" then sc="每周1 0点"
            when "m" then sc="每月1日0点"
            when "y"
                if startMonth >1
                    sc="每年#{startMonth}月1日0点"
                else
                    sc="每年1月1日0点"
            else ""
          if @Cyc in ['ss','mi']
            sc="！！！#{sc}#{s['mi']}分#{s['ss']}秒"
          else
            sc="#{sc}#{s['mi']}分#{s['ss']}秒"
        when "h"
          switch @Cyc
            when "d" then sc="每日"
            when "m" then sc="每月1日"
            when "y"
                if startMonth >1
                    sc="每年#{startMonth}月1日"
                else
                    sc="每年1月1日"
            else ""
          if @Cyc in ['ss','mi','h']
            sc="！！！#{sc}#{s['h']}点#{s['mi']}分#{s['ss']}秒"
          else
            sc="#{sc}#{s['h']}点#{s['mi']}分#{s['ss']}秒"
        when "d"
          switch @Cyc
            when "w" then sc="每周"
            when "m" then sc="每月"
            when "y"
                if startMonth >1
                    sc="每年#{startMonth}月"
                else
                    sc="每年1月"
            else ""
          if @Cyc in ['ss','mi','h','d']
            sc="！！！#{sc}#{s['d']}日#{s['h']}点#{s['mi']}分#{s['ss']}秒"
          else
            sc="#{sc}#{s['d']}日#{s['h']}点#{s['mi']}分#{s['ss']}秒"
        else "未设置"
  
module.exports = Schedule
