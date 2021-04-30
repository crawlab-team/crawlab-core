package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/utils"
)

func initModels() (err error) {
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

	return nil
}

func InitModels() (err error) {
	return utils.InitModule(interfaces.ModuleIdModels, initModels)
}
