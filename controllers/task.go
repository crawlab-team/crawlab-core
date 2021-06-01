package controllers

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/spider/admin"
	clog "github.com/crawlab-team/crawlab-log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"net/http"
	"sync"
)

var TaskController ListActionController

var TaskActions = []Action{
	{
		Method:      http.MethodPut,
		Path:        "/run",
		HandlerFunc: taskCtx.run,
	},
	{
		Method:      http.MethodPost,
		Path:        "/:id/restart",
		HandlerFunc: taskCtx.restart,
	},
	{
		Method:      http.MethodGet,
		Path:        "/:id/logs",
		HandlerFunc: taskCtx.getLogs,
	},
}

type taskContext struct {
	modelSvc service.ModelService
	adminSvc interfaces.SpiderAdminService
	l        clog.Driver

	// internals
	drivers sync.Map
}

func (ctx *taskContext) run(c *gin.Context) {
	// task
	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// validate spider id
	if t.GetSpiderId().IsZero() {
		HandleErrorBadRequest(c, errors.ErrorTaskEmptySpiderId)
		return
	}

	// spider
	s, err := ctx.modelSvc.GetSpiderById(t.GetSpiderId())
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// options
	opts := &interfaces.SpiderRunOptions{
		Mode:    t.Mode,
		NodeIds: t.NodeIds,
		Param:   t.Param,
	}

	// run
	if err := ctx.adminSvc.Schedule(s.GetId(), opts); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *taskContext) restart(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// task
	t, err := ctx.modelSvc.GetTaskById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// options
	opts := &interfaces.SpiderRunOptions{
		Mode:    t.Mode,
		NodeIds: t.NodeIds,
		Param:   t.Param,
	}

	// run
	if err := ctx.adminSvc.Schedule(t.SpiderId, opts); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *taskContext) getLogs(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// pagination
	p, err := GetPagination(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// log driver
	l, err := ctx._getLogDriver(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// logs
	logs, err := l.Find("", (p.Page-1)*p.Size, p.Size)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	total, err := l.Count("")
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithListData(c, logs, total)
}

func (ctx *taskContext) _getLogDriver(id primitive.ObjectID) (l clog.Driver, err error) {
	// attempt to get from cache
	res, ok := ctx.drivers.Load(id)
	if ok {
		l, ok = res.(clog.Driver)
		if ok {
			return l, nil
		}
	}

	// TODO: other types of log drivers
	l, err = clog.NewSeaweedFsLogDriver(&clog.SeaweedFsLogDriverOptions{Prefix: id.Hex()})
	if err != nil {
		return nil, err
	}
	ctx.drivers.Store(id, l)

	return l, nil
}

var taskCtx = newTaskContext()

func newTaskContext() *taskContext {
	// context
	ctx := &taskContext{
		drivers: sync.Map{},
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		panic(err)
	}
	if err := c.Provide(admin.NewSpiderAdminService); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(modelSvc service.ModelService, adminSvc interfaces.SpiderAdminService) {
		ctx.modelSvc = modelSvc
		ctx.adminSvc = adminSvc
	}); err != nil {
		panic(err)
	}

	return ctx
}
