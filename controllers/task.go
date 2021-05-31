package controllers

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/spider/admin"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"net/http"
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
}

type taskContext struct {
	modelSvc service.ModelService
	adminSvc interfaces.SpiderAdminService
}

func (ctx taskContext) run(c *gin.Context) {
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

func (ctx taskContext) restart(c *gin.Context) {
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

var taskCtx = newTaskContext()

func newTaskContext() *taskContext {
	// context
	ctx := &taskContext{}

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
