package controllers

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/interfaces"
	delegate2 "github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/plugin"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"net/http"
)

var PluginController *pluginController

func getPluginActions() []Action {
	pluginCtx := newPluginContext()
	return []Action{
		{
			Method:      http.MethodPost,
			Path:        "/:id/run",
			HandlerFunc: pluginCtx.run,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/stop",
			HandlerFunc: pluginCtx.stop,
		},
	}
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

func (ctr *pluginController) Delete(c *gin.Context) {
	_, err := ctr.ctx._delete(c)
	if err != nil {
		return
	}
	HandleSuccess(c)
}

type pluginContext struct {
	modelSvc  service.ModelService
	pluginSvc interfaces.PluginService
}

var _pluginCtx *pluginContext

func (ctx *pluginContext) run(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := ctx.pluginSvc.RunPlugin(id); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *pluginContext) stop(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := ctx.pluginSvc.StopPlugin(id); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *pluginContext) _put(c *gin.Context) (p *models.Plugin, err error) {
	// bind
	p = &models.Plugin{}
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleErrorBadRequest(c, err)
		return nil, err
	}

	// TODO: check if exists

	// add
	if err := delegate2.NewModelDelegate(p, GetUserFromContext(c)).Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// install
	if err := ctx.pluginSvc.InstallPlugin(p.GetId()); err != nil {
		_ = delegate2.NewModelDelegate(p).Delete()
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// run
	if err := ctx.pluginSvc.RunPlugin(p.GetId()); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	return p, nil
}

func (ctx *pluginContext) _delete(c *gin.Context) (p *models.Plugin, err error) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// plugin
	p, err = ctx.modelSvc.GetPluginById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// delete
	if err := delegate2.NewModelDelegate(p, GetUserFromContext(c)).Delete(); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// stop
	if p.Status == constants.PluginStatusRunning {
		if err := ctx.pluginSvc.StopPlugin(p.GetId()); err != nil {
			HandleErrorInternalServerError(c, err)
			return nil, err
		}
	}

	// TODO: uninstall
	//if err := ctx.pluginSvc.UninstallPlugin(p.GetId()); err != nil {
	//	HandleErrorInternalServerError(c, err)
	//	return nil, err
	//}

	return p, nil
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
