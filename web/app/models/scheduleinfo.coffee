Spine = require('spine')

class MScheduleInfo extends Spine.Model
  @configure 'Schedule', 'Id', 'Name', 'TaskCnt', 'Count', 'Cyc', 'StartMonth', 'StartSecond', 'NextStart', 'TimeOut', 'Desc', 'CreateTime', 'CreateUserId', 'ModifyTime', 'ModifyUserId'

  @extend Spine.Model.Ajax

  @url: "/schedules"
  
module.exports = MScheduleInfo
