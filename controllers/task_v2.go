package controllers

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/result"
	"github.com/crawlab-team/crawlab-core/spider/admin"
	"github.com/crawlab-team/crawlab-core/task/log"
	"github.com/crawlab-team/crawlab-core/task/scheduler"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/crawlab-db/generic"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

func PostTaskRun(c *gin.Context) {
	// task
	var t models.TaskV2
	if err := c.ShouldBindJSON(&t); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// validate spider id
	if t.SpiderId.IsZero() {
		HandleErrorBadRequest(c, errors.ErrorTaskEmptySpiderId)
		return
	}

	// spider
	s, err := service.NewModelServiceV2[models.SpiderV2]().GetById(t.SpiderId)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// options
	opts := &interfaces.SpiderRunOptions{
		Mode:     t.Mode,
		NodeIds:  t.NodeIds,
		Cmd:      t.Cmd,
		Param:    t.Param,
		Priority: t.Priority,
	}

	// user
	if u := GetUserFromContextV2(c); u != nil {
		opts.UserId = u.Id
	}

	// run
	adminSvc, err := admin.GetSpiderAdminServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	taskIds, err := adminSvc.Schedule(s.Id, opts)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, taskIds)

}

func PostTaskRestart(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// task
	t, err := service.NewModelServiceV2[models.TaskV2]().GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// options
	opts := &interfaces.SpiderRunOptions{
		Mode:     t.Mode,
		NodeIds:  t.NodeIds,
		Cmd:      t.Cmd,
		Param:    t.Param,
		Priority: t.Priority,
	}

	// user
	if u := GetUserFromContextV2(c); u != nil {
		opts.UserId = u.Id
	}

	// run
	adminSvc, err := admin.GetSpiderAdminServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	taskIds, err := adminSvc.Schedule(t.SpiderId, opts)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, taskIds)
}

func PostTaskCancel(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// task
	t, err := service.NewModelServiceV2[models.TaskV2]().GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// validate
	if !utils.IsCancellable(t.Status) {
		HandleErrorInternalServerError(c, errors.ErrorControllerNotCancellable)
		return
	}

	u := GetUserFromContextV2(c)

	// cancel
	schedulerSvc, err := scheduler.GetTaskSchedulerServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if err := schedulerSvc.Cancel(id, u.Id); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func GetTaskLogs(c *gin.Context) {
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

	// logs
	logDriver, err := log.GetFileLogDriver()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	logs, err := logDriver.Find(id.Hex(), "", (p.Page-1)*p.Size, p.Size)
	if err != nil {
		if strings.HasSuffix(err.Error(), "Status:404 Not Found") {
			HandleSuccess(c)
			return
		}
		HandleErrorInternalServerError(c, err)
		return
	}
	total, err := logDriver.Count(id.Hex(), "")
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithListData(c, logs, total)
}

func GetTaskData(c *gin.Context) {
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
	t, err := service.NewModelServiceV2[models.TaskV2]().GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// result service
	resultSvc, err := result.GetResultService(t.SpiderId)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// query
	query := generic.ListQuery{
		generic.ListQueryCondition{
			Key:   constants.TaskKey,
			Op:    generic.OpEqual,
			Value: t.Id,
		},
	}

	// list
	data, err := resultSvc.List(query, &generic.ListOptions{
		Skip:  (p.Page - 1) * p.Size,
		Limit: p.Size,
		Sort:  []generic.ListSort{{"_id", generic.SortDirectionDesc}},
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
