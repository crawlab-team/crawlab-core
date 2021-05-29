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

	NodeController = NewListControllerDelegate(ControllerIdNode, modelSvc.NewBaseService(interfaces.ModelIdNode))
	ProjectController = NewListControllerDelegate(ControllerIdProject, modelSvc.NewBaseService(interfaces.ModelIdProject))
	SpiderController = NewListPostActionControllerDelegate(ControllerIdSpider, modelSvc.NewBaseService(interfaces.ModelIdSpider), SpiderActions)
	TaskController = NewListPostActionControllerDelegate(ControllerIdTask, modelSvc.NewBaseService(interfaces.ModelIdTask), TaskActions)
	UserController = NewListControllerDelegate(ControllerIdUser, modelSvc.NewBaseService(interfaces.ModelIdUser))
	TagController = NewListControllerDelegate(ControllerIdTag, modelSvc.NewBaseService(interfaces.ModelIdTag))
	LoginController = NewActionControllerDelegate(ControllerIdLogin, LoginActions)
	ColorController = NewActionControllerDelegate(ControllerIdColor, ColorActions)

	return nil
}
