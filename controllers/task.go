package controllers

import (
	"github.com/crawlab-team/crawlab-core/config"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/result"
	"github.com/crawlab-team/crawlab-core/spider/admin"
	"github.com/crawlab-team/crawlab-core/task/scheduler"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/crawlab-db/mongo"
	clog "github.com/crawlab-team/crawlab-log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"net/http"
	"strings"
	"sync"
)

var TaskController *taskController

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
		Method:      http.MethodPost,
		Path:        "/:id/cancel",
		HandlerFunc: taskCtx.cancel,
	},
	{
		Method:      http.MethodGet,
		Path:        "/:id/logs",
		HandlerFunc: taskCtx.getLogs,
	},
	{
		Method:      http.MethodGet,
		Path:        "/:id/data",
		HandlerFunc: taskCtx.getData,
	},
}

type taskController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *taskContext
}

func (ctr *taskController) GetList(c *gin.Context) {
	withStats := c.Query("stats")
	if withStats == "" {
		ctr.d.GetList(c)
		return
	}
	ctr.ctx.getListWithStats(c)
}

type taskContext struct {
	modelSvc     service.ModelService
	modelTaskSvc interfaces.ModelBaseService
	adminSvc     interfaces.SpiderAdminService
	schedulerSvc interfaces.TaskSchedulerService
	l            clog.Driver

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
		Mode:     t.Mode,
		NodeIds:  t.NodeIds,
		Param:    t.Param,
		Priority: t.Priority,
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
		Mode:     t.Mode,
		NodeIds:  t.NodeIds,
		Param:    t.Param,
		Priority: t.Priority,
	}

	// run
	if err := ctx.adminSvc.Schedule(t.SpiderId, opts); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *taskContext) cancel(c *gin.Context) {
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

	// validate
	if !utils.IsCancellable(t.Status) {
		HandleErrorInternalServerError(c, errors.ErrorControllerNotCancellable)
		return
	}

	// run
	if err := ctx.schedulerSvc.Cancel(id); err != nil {
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
		if strings.HasSuffix(err.Error(), "Status:404 Not Found") {
			HandleSuccess(c)
			return
		}
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

func (ctx *taskContext) getListWithStats(c *gin.Context) {
	// params
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)
	sort := MustGetSortOption(c)

	// get list
	list, err := ctx.modelTaskSvc.GetList(query, &mongo.FindOptions{
		Sort:  sort,
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			HandleErrorNotFound(c, err)
		} else {
			HandleErrorInternalServerError(c, err)
		}
		return
	}

	// check empty list
	if len(list.Values()) == 0 {
		HandleSuccessWithListData(c, nil, 0)
		return
	}

	// ids
	var ids []primitive.ObjectID
	for _, d := range list.Values() {
		t := d.(interfaces.Model)
		ids = append(ids, t.GetId())
	}

	// total count
	total, err := ctx.modelTaskSvc.Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// stat list
	query = bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	stats, err := ctx.modelSvc.GetTaskStatList(query, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// cache stat list to dict
	dict := map[primitive.ObjectID]models.TaskStat{}
	for _, s := range stats {
		dict[s.GetId()] = s
	}

	// iterate list again
	var data []interface{}
	for _, d := range list.Values() {
		t := d.(*models.Task)
		s, ok := dict[t.GetId()]
		if ok {
			t.Stat = &s
		}
		data = append(data, *t)
	}

	// response
	HandleSuccessWithListData(c, data, total)
}

func (ctx *taskContext) getData(c *gin.Context) {
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

	// task
	t, err := ctx.modelSvc.GetTaskById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// spider
	s, err := ctx.modelSvc.GetSpiderById(t.SpiderId)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// result service
	resultSvc, err := result.GetResultService(s.ColId)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// query
	query := bson.M{
		"_tid": id,
	}

	// list
	data, err := resultSvc.GetList(query, &mongo.FindOptions{
		Skip:  (p.Page - 1) * p.Size,
		Limit: p.Size,
		Sort:  bson.D{{"_id", -1}},
	})
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// total
	total, err := resultSvc.Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithListData(c, data, total)
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
	if err := c.Provide(scheduler.ProvideGetTaskSchedulerService(config.DefaultConfigPath)); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		adminSvc interfaces.SpiderAdminService,
		schedulerSvc interfaces.TaskSchedulerService,
	) {
		ctx.modelSvc = modelSvc
		ctx.adminSvc = adminSvc
		ctx.schedulerSvc = schedulerSvc
	}); err != nil {
		panic(err)
	}

	// model task service
	ctx.modelTaskSvc = ctx.modelSvc.GetBaseService(interfaces.ModelIdTask)

	return ctx
}

func newTaskController() *taskController {
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdTask, modelSvc.GetBaseService(interfaces.ModelIdTask), TaskActions)
	d := NewListPostActionControllerDelegate(ControllerIdTask, modelSvc.GetBaseService(interfaces.ModelIdTask), TaskActions)
	ctx := newTaskContext()

	return &taskController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
