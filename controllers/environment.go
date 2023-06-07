package controllers

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/user"
	"go.uber.org/dig"
)

var EnvironmentController *environmentController

var EnvironmentActions []Action

type environmentController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *environmentContext
}

type environmentContext struct {
	modelSvc service.ModelService
	userSvc  interfaces.UserService
}

func newEnvironmentContext() *environmentContext {
	// context
	ctx := &environmentContext{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		panic(err)
	}
	if err := c.Provide(user.ProvideGetUserService()); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		userSvc interfaces.UserService,
	) {
		ctx.modelSvc = modelSvc
		ctx.userSvc = userSvc
	}); err != nil {
		panic(err)
	}

	return ctx
}

func newEnvironmentController() *environmentController {
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdEnvironment, modelSvc.GetBaseService(interfaces.ModelIdEnvironment), EnvironmentActions)
	d := NewListPostActionControllerDelegate(ControllerIdEnvironment, modelSvc.GetBaseService(interfaces.ModelIdEnvironment), EnvironmentActions)
	ctx := newEnvironmentContext()

	return &environmentController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
