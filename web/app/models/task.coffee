Spine = require('spineify')

class Task extends Spine.Model
  @configure 'Task', "Id", "Address", "Name", "JobType", "ScheduleCyc", "TaskCyc", "StartSecond", "Cmd", "Desc", "TimeOut", "Param", "Attr", "JobId", "RelTasks", "RelTaskCnt", "CreateUserId", "CreateTime", "ModifyUserId", "ModifyTime"

  @extend Spine.Model.Ajax

module.exports = Task
