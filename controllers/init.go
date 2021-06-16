package controllers

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
)

func InitControllers() (err error) {
	modelSvc, err := service.GetService()
	if err != nil {
		return err
	}

	NodeController = newNodeController()
	ProjectController = newProjectController()
	SpiderController = newSpiderController()
	TaskController = newTaskController()
	UserController = NewListControllerDelegate(ControllerIdUser, modelSvc.GetBaseService(interfaces.ModelIdUser))
	TagController = NewListControllerDelegate(ControllerIdTag, modelSvc.GetBaseService(interfaces.ModelIdTag))
	LoginController = NewActionControllerDelegate(ControllerIdLogin, LoginActions)
	ColorController = NewActionControllerDelegate(ControllerIdColor, ColorActions)
	PluginController = NewListControllerDelegate(ControllerIdPlugin, modelSvc.GetBaseService(interfaces.ModelIdPlugin))
	DataCollectionController = NewListControllerDelegate(ControllerIdDataCollection, modelSvc.GetBaseService(interfaces.ModelIdDataCollection))
	ResultController = NewActionControllerDelegate(ControllerIdResult, ResultActions)
	ScheduleController = newScheduleController()

	return nil
}
