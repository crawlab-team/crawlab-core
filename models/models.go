package models

var initialized = false

func InitModelServices() (err error) {
	// skip if already initialized
	if initialized {
		return nil
	}

	// system model services
	ArtifactService = NewArtifactService()
	TagService = NewTagService()
	ColorService = NewColorService()

	// operation model services
	NodeService = NewNodeService()
	ProjectService = NewProjectService()
	ScheduleService = NewScheduleService()
	SettingService = NewSettingService()
	SpiderService = NewSpiderService()
	TaskService = NewTaskService()
	JobService = NewJobService()
	TokenService = NewTokenService()
	UserService = NewUserService()
	VariableService = NewVariableService()

	// mark as initialized
	initialized = true

	return nil
}
