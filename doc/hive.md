
管理模块的API

/v0.0.1/schedules
	GET  列出所有调度
	POST 新建一个调度

/v0.0.1/schedules/ID
	GET  列出指定调度的信息
	PUT 更新指定调度
	DELETE 删除指定调度


/v0.0.1/jobs
	GET  列出所有作业
	POST 新建一个作业

/v0.0.1/jobs/ID
	GET  列出指定作业
	PUT 更新指定作业
	DELETE 删除指定作业


/v0.0.1/tasks
	GET  列出所有任务
	POST 新建一个任务

/v0.0.1/tasks/ID
	GET  列出指定任务
	PUT 更新指定任务
	DELETE 删除指定任务

/v0.0.1/tasks/ID/tasks
	GET  列出指定任务的依赖任务
	POST 为指定任务新建一个依赖任务

/v0.0.1/tasks/ID/tasks/ID
	DELETE 删除指定任务的依赖任务

/v0.0.1/schedules/ID/jobs
	GET  列出指定调度下的所有作业
	POST 指定调度下新建一个作业
/v0.0.1/schedules/ID/jobs/ID
	DELETE 删除指定调度下指定作业


/v0.0.1/schedules/ID/tasks
	GET  列出指定调度下的所有任务

/v0.0.1/schedules/ID/jobs/ID/tasks
	GET  列出指定作业下的所有任务
	POST 指定作业下新建一个作业
/v0.0.1/schedules/ID/jobs/ID/tasks/ID
	DELETE 删除指定作业下指定任务


