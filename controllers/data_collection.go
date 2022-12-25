package controllers

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	"go.uber.org/dig"
)

var DataCollectionController *dataCollectionController

func getDataCollectionActions() []Action {
	//ctx := newDataCollectionContext()
	return []Action{}
}

type dataCollectionController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *dataCollectionContext
}

type dataCollectionContext struct {
	modelSvc  service.ModelService
	resultSvc interfaces.ResultService
}

var _dataCollectionCtx *dataCollectionContext

func newDataCollectionContext() *dataCollectionContext {
	if _dataCollectionCtx != nil {
		return _dataCollectionCtx
	}

	// context
	ctx := &dataCollectionContext{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
	) {
		ctx.modelSvc = modelSvc
	}); err != nil {
		panic(err)
	}

	_dataCollectionCtx = ctx

	return ctx
}

func newDataCollectionController() *dataCollectionController {
	actions := getDataCollectionActions()
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdDataCollection, modelSvc.GetBaseService(interfaces.ModelIdDataCollection), actions)
	d := NewListPostActionControllerDelegate(ControllerIdDataCollection, modelSvc.GetBaseService(interfaces.ModelIdDataCollection), actions)
	ctx := newDataCollectionContext()

	return &dataCollectionController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
