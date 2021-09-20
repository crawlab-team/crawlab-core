package controllers

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	delegate2 "github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/plugin"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

var PluginController *pluginController

func getPluginActions() []Action {
	return []Action{}
}

type pluginController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *pluginContext
}

func (ctr *pluginController) Put(c *gin.Context) {
	s, err := ctr.ctx._put(c)
	if err != nil {
		return
	}
	HandleSuccessWithData(c, s)
}

type pluginContext struct {
	modelSvc  service.ModelService
	pluginSvc interfaces.PluginService
}

var _pluginCtx *pluginContext

func (ctx *pluginContext) _put(c *gin.Context) (s *models.Plugin, err error) {
	// bind
	s = &models.Plugin{}
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return nil, err
	}

	// add
	if err := delegate2.NewModelDelegate(s, GetUserFromContext(c)).Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// install
	if err := ctx.pluginSvc.InstallPlugin(s.GetId()); err != nil {
		_ = delegate2.NewModelDelegate(s).Delete()
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	return s, nil
}

func newPluginContext() *pluginContext {
	if _pluginCtx != nil {
		return _pluginCtx
	}

	// context
	ctx := &pluginContext{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		panic(err)
	}
	if err := c.Provide(plugin.NewPluginService); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		pluginSvc interfaces.PluginService,
	) {
		ctx.modelSvc = modelSvc
		ctx.pluginSvc = pluginSvc
	}); err != nil {
		panic(err)
	}

	_pluginCtx = ctx

	return ctx
}

func newPluginController() *pluginController {
	actions := getPluginActions()
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdPlugin, modelSvc.GetBaseService(interfaces.ModelIdPlugin), actions)
	d := NewListPostActionControllerDelegate(ControllerIdPlugin, modelSvc.GetBaseService(interfaces.ModelIdPlugin), actions)
	ctx := newPluginContext()

	return &pluginController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
