Spine = require('spineify')

class Job extends Spine.Model
  @configure 'Job', 'Id', 'ScheduleId', 'ScheduleCyc', 'Name', 'Desc', 'PreJobId', 'NextJobId', 'Tasks', 'TaskCnt', 'CreateUserId', 'CreateTime', 'ModifyUserId', 'ModifyTime'

  @extend Spine.Model.Ajax
  
module.exports = Job
