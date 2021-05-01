package controllers

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/store"
)

func InitControllers() (err error) {
	NodeController = NewListControllerDelegate(ControllerIdNode, store.ModelServiceStore.MustGetModelService(interfaces.ModelIdNode))
	ProjectController = NewListControllerDelegate(ControllerIdProject, store.ModelServiceStore.MustGetModelService(interfaces.ModelIdProject))
	SpiderController = NewListPostActionControllerDelegate(ControllerIdSpider, store.ModelServiceStore.MustGetModelService(interfaces.ModelIdSpider), SpiderActions)
	TaskController = NewListPostActionControllerDelegate(ControllerIdTask, store.ModelServiceStore.MustGetModelService(interfaces.ModelIdTask), TaskActions)
	UserController = NewListControllerDelegate(ControllerIdUser, store.ModelServiceStore.MustGetModelService(interfaces.ModelIdUser))
	TagController = NewListControllerDelegate(ControllerIdTag, store.ModelServiceStore.MustGetModelService(interfaces.ModelIdTag))
	LoginController = NewActionControllerDelegate(ControllerIdLogin, LoginActions)
	ColorController = NewActionControllerDelegate(ControllerIdColor, ColorActions)

	return nil
}
