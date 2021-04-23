package controllers

import "github.com/crawlab-team/crawlab-core/models"

func InitControllers() (err error) {
	NodeController = NewListControllerDelegate(ControllerIdNode, models.NodeService)
	ProjectController = NewListControllerDelegate(ControllerIdProject, models.ProjectService)
	SpiderController = NewListPostActionControllerDelegate(ControllerIdSpider, models.SpiderService, SpiderActions)
	TaskController = NewListPostActionControllerDelegate(ControllerIdTask, models.TaskService, TaskActions)
	UserController = NewListControllerDelegate(ControllerIdUser, models.UserService)
	TagController = NewListControllerDelegate(ControllerIdTag, models.TagService)
	LoginController = NewActionControllerDelegate(ControllerIdLogin, LoginActions)
	ColorController = NewActionControllerDelegate(ControllerIdColor, ColorActions)

	return nil
}
