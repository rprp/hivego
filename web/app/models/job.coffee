Spine = require('spineify')

class Job extends Spine.Model
  @configure 'Schedule', 'Id', 'Name', 'TaskCnt', 'Job', 'Count', 'Cyc', 'StartMonth', 'StartSecond', 'NextStart', 'TimeOut', 'Desc', 'CreateTime', 'CreateUserId', 'ModifyTime', 'ModifyUserId'

  @extend Spine.Model.Ajax
  
  constructor: ->
    Moment.lang('zh-cn')
    super
 
module.exports = Job
