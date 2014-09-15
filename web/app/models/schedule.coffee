Moment = require('momentifycn')
Spine = require('spineify')
Job = require('models/job')

class Schedule extends Spine.Model
  @configure 'Schedule', 'Id', 'Name', 'TaskCnt', "Job", "JobCnt", "Jobs", "JobId", 'Count', 'Cyc', 'StartMonth', 'StartSecond', 'NextStart', 'TimeOut', 'Desc', 'CreateTime', 'CreateUserId', 'ModifyTime', 'ModifyUserId'

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

  SetCyc: (cyc) ->
    switch cyc
      when "Second" then @Cyc = "ss"
      when  "Minute" then @Cyc = "mi"
      when "Hour" then @Cyc = "h"
      when "Day" then @Cyc = "d"
      when "Week" then @Cyc = "w"
      when "Month" then @Cyc = "m"
      when "Year" then @Cyc = "y"
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

  ParseSecond: (sd) ->
    ss = 0
    sc = ""
    second = 1000 * 1000 * 1000
    mi = 60 * second
    h = 60 * mi
    d = 24 * h
    k = 0
    month = 0

    clst = ["年","月","周","日","点","分","秒"]
    switch @Cyc
      when "mi" then sc="分"
      when "h" then sc="小时"
      when "d" then sc="日"
      when "w" then sc="周"
      when "m" then sc="月"
      when "y" then sc="年"
      else ""

    console.log(sd)
    console.log(sd[1..2])
    if ((sc is sd[1]) or (@Cyc is "h" and sd[1..2] is "小时"))
      tp = ""
      j = []
      for v,i in sd[1..]
        if v in clst
          if tp isnt ""
            switch v
              when "秒"
                if parseInt(tp) > 60
                  return [0,-1]
                else
                  k += parseInt(tp) * second
                  j.push(parseInt(tp) * second)
              when "分"
                if parseInt(tp) > 60
                  return [0,-1]
                else
                  k += parseInt(tp) * mi
                  j.push(parseInt(tp) * mi)
              when "点"
                if parseInt(tp) > 24
                  return [0,-1]
                else
                  k += parseInt(tp) * h
                  j.push(parseInt(tp) * h)
              when "日"
                if parseInt(tp) > 30
                  return [0,-1]
                else
                  k += parseInt(tp) * d
                  j.push(parseInt(tp) * d)
              when "月"
                if parseInt(tp) > 12
                  return [0,-1]
                else
                  month = parseInt(tp)
              else ""
            tp = ""
        else if not isNaN(v)
          tp = "#{tp}#{v}"
          
    else
      return [0,-1]

    return [month, k]


  GetDefaultSecond: ->
    switch @Cyc
      when "mi" then sc="每分"
      when "h" then sc="每小时0分"
      when "d" then sc="每日0点0分"
      when "w" then sc="每周1 0点0分"
      when "m" then sc="每月1日0点0分"
      when "y" then sc="每年1月1日0点0分"
      else ""
    sc="#{sc}1秒"

  GetSecond: ->
    return [] unless @StartSecond
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
