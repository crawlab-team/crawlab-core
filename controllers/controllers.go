package controllers

import (
	"github.com/crawlab-team/crawlab-core/inject"
	"github.com/crawlab-team/crawlab-core/interfaces"
)

func InitControllers() (err error) {
	NodeController = NewListControllerDelegate(ControllerIdNode, inject.ModelServiceStore.MustGetModelService(interfaces.ModelIdNode))
	ProjectController = NewListControllerDelegate(ControllerIdProject, inject.ModelServiceStore.MustGetModelService(interfaces.ModelIdProject))
	SpiderController = NewListPostActionControllerDelegate(ControllerIdSpider, inject.ModelServiceStore.MustGetModelService(interfaces.ModelIdSpider), SpiderActions)
	TaskController = NewListPostActionControllerDelegate(ControllerIdTask, inject.ModelServiceStore.MustGetModelService(interfaces.ModelIdTask), TaskActions)
	UserController = NewListControllerDelegate(ControllerIdUser, inject.ModelServiceStore.MustGetModelService(interfaces.ModelIdUser))
	TagController = NewListControllerDelegate(ControllerIdTag, inject.ModelServiceStore.MustGetModelService(interfaces.ModelIdTag))
	LoginController = NewActionControllerDelegate(ControllerIdLogin, LoginActions)
	ColorController = NewActionControllerDelegate(ControllerIdColor, ColorActions)

	return nil
}
