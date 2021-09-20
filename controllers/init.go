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
	UserController = newUserController()
	TagController = NewListControllerDelegate(ControllerIdTag, modelSvc.GetBaseService(interfaces.ModelIdTag))
	LoginController = NewActionControllerDelegate(ControllerIdLogin, getLoginActions())
	ColorController = NewActionControllerDelegate(ControllerIdColor, getColorActions())
	PluginController = newPluginController()
	DataCollectionController = NewListControllerDelegate(ControllerIdDataCollection, modelSvc.GetBaseService(interfaces.ModelIdDataCollection))
	ResultController = NewActionControllerDelegate(ControllerIdResult, getResultActions())
	ScheduleController = newScheduleController()
	StatsController = NewActionControllerDelegate(ControllerIdStats, getStatsActions())
	TokenController = newTokenController()
	FilerController = NewActionControllerDelegate(ControllerIdFiler, getFilerActions())
	PluginProxyController = NewActionControllerDelegate(ControllerIdPluginDo, getPluginProxyActions())

	return nil
}
