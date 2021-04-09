package models

func InitModelServices() (err error) {
	JobService = NewJobService()
	NodeService = NewNodeService()
	ProjectService = NewProjectService()
	ScheduleService = NewScheduleService()
	SettingService = NewSettingService()
	SpiderService = NewSpiderService()
	TaskService = NewTaskService()
	TokenService = NewTokenService()
	UserService = NewUserService()
	VariableService = NewVariableService()

	return nil
}
