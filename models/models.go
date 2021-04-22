package models

func InitModelServices() (err error) {
	// system model services
	ArtifactService = NewArtifactService()
	TagService = NewTagService()
	ColorService = NewColorService()

	// operation model services
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
