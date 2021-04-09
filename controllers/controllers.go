package controllers

import "github.com/crawlab-team/crawlab-core/models"

func InitControllers() (err error) {
	ProjectController = NewListControllerDelegate(ControllerIdProject, models.ProjectService)
	UserController = NewListControllerDelegate(ControllerIdUser, models.UserService)
	SpiderController = NewListControllerDelegate(ControllerIdSpider, models.SpiderService)

	return nil
}
