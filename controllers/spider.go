package controllers

import (
	"bytes"
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	delegate2 "github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/spider/admin"
	"github.com/crawlab-team/crawlab-core/spider/sync"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"io"
	"math"
	"net/http"
	"strings"
)

var SpiderController *spiderController

func getSpiderActions() []Action {
	spiderCtx := newSpiderContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "/:id/files/list",
			HandlerFunc: spiderCtx.listDir,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id/files/get",
			HandlerFunc: spiderCtx.getFile,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id/files/info",
			HandlerFunc: spiderCtx.getFileInfo,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/files/save",
			HandlerFunc: spiderCtx.saveFile,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/files/save/dir",
			HandlerFunc: spiderCtx.saveDir,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/files/rename",
			HandlerFunc: spiderCtx.renameFile,
		},
		{
			Method:      http.MethodDelete,
			Path:        "/:id/files/delete",
			HandlerFunc: spiderCtx.delete,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/files/copy",
			HandlerFunc: spiderCtx.copyFile,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/run",
			HandlerFunc: spiderCtx.run,
		},
		//{
		//	Method:      http.MethodPost,
		//	Path:        "/:id/clone",
		//	HandlerFunc: spiderCtx.clone,
		//},
	}
}

type spiderController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *spiderContext
}

func (ctr *spiderController) Get(c *gin.Context) {
	ctr.ctx._get(c)
}

func (ctr *spiderController) Put(c *gin.Context) {
	s, err := ctr.ctx._put(c)
	if err != nil {
		return
	}
	HandleSuccessWithData(c, s)
}

func (ctr *spiderController) Post(c *gin.Context) {
	s, err := ctr.ctx._post(c)
	if err != nil {
		return
	}
	HandleSuccessWithData(c, s)
}

func (ctr *spiderController) GetList(c *gin.Context) {
	withStats := c.Query("stats")
	if withStats == "" {
		ctr.d.GetList(c)
		return
	}
	ctr.ctx._getListWithStats(c)
}

type spiderContext struct {
	modelSvc       service.ModelService
	modelSpiderSvc interfaces.ModelBaseService
	syncSvc        interfaces.SpiderSyncService
	adminSvc       interfaces.SpiderAdminService
}

func (ctx *spiderContext) listDir(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodGet)
	if err != nil {
		return
	}

	files, err := fsSvc.List(payload.Path)
	if err != nil {
		if err.Error() != "response status code: 404" {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	HandleSuccessWithData(c, files)
}

func (ctx *spiderContext) getFile(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodGet)
	if err != nil {
		return
	}

	data, err := fsSvc.GetFile(payload.Path)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	data = utils.TrimFileData(data)

	HandleSuccessWithData(c, string(data))
}

func (ctx *spiderContext) getFileInfo(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodGet)
	if err != nil {
		return
	}

	info, err := fsSvc.GetFileInfo(payload.Path)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, info)
}

func (ctx *spiderContext) saveFile(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodPost)
	if err != nil {
		return
	}

	data := utils.FillEmptyFileData([]byte(payload.Data))

	if err := fsSvc.Save(payload.Path, data); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) saveDir(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodPost)
	if err != nil {
		return
	}

	data := []byte("")
	path := fmt.Sprintf("%s/%s", payload.Path, constants.FsKeepFileName)

	if err := fsSvc.Save(path, data); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) renameFile(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodPost)
	if err != nil {
		return
	}

	if err := fsSvc.Rename(payload.Path, payload.NewPath); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) delete(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodPost)
	if err != nil {
		return
	}

	if err := fsSvc.Delete(payload.Path); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) copyFile(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodPost)
	if err != nil {
		return
	}

	if err := fsSvc.Copy(payload.Path, payload.NewPath); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) run(c *gin.Context) {
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	var opts interfaces.SpiderRunOptions
	if err := c.ShouldBindJSON(&opts); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	if err := ctx.adminSvc.Schedule(id, &opts); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) _get(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	s, err := ctx.modelSvc.GetSpiderById(id)
	if err == mongo2.ErrNoDocuments {
		HandleErrorNotFound(c, err)
		return
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// stat
	s.Stat, err = ctx.modelSvc.GetSpiderStatById(s.GetId())
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// data collection
	if !s.ColId.IsZero() {
		col, err := ctx.modelSvc.GetDataCollectionById(s.ColId)
		if err != nil {
			if err != mongo2.ErrNoDocuments {
				HandleErrorInternalServerError(c, err)
				return
			}
		} else {
			s.ColName = col.Name
		}
	}

	HandleSuccessWithData(c, s)
}

func (ctx *spiderContext) _post(c *gin.Context) (s *models.Spider, err error) {
	// bind
	s = &models.Spider{}
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return nil, err
	}

	// upsert data collection
	if err := ctx._upsertDataCollection(s); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// save
	if err := delegate2.NewModelDelegate(s).Save(); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	return s, nil
}

func (ctx *spiderContext) _put(c *gin.Context) (s *models.Spider, err error) {
	// bind
	s = &models.Spider{}
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return nil, err
	}

	// upsert data collection
	if err := ctx._upsertDataCollection(s); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// add
	if err := delegate2.NewModelDelegate(s).Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// add stat
	st := &models.SpiderStat{
		Id: s.GetId(),
	}
	if err := delegate2.NewModelDelegate(st).Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	return s, nil
}

func (ctx *spiderContext) _getListWithStats(c *gin.Context) {
	// params
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)
	sort := MustGetSortOption(c)

	// get list
	list, err := ctx.modelSpiderSvc.GetList(query, &mongo.FindOptions{
		Sort:  sort,
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if err.Error() == mongo2.ErrNoDocuments.Error() {
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
		s := d.(*models.Spider)
		ids = append(ids, s.GetId())
	}

	// total count
	total, err := ctx.modelSpiderSvc.Count(query)
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
	stats, err := ctx.modelSvc.GetSpiderStatList(query, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// cache stat list to dict
	dict := map[primitive.ObjectID]models.SpiderStat{}
	var tids []primitive.ObjectID
	for _, st := range stats {
		if st.Tasks > 0 {
			taskCount := int64(st.Tasks)
			st.AverageWaitDuration = int64(math.Round(float64(st.WaitDuration) / float64(taskCount)))
			st.AverageRuntimeDuration = int64(math.Round(float64(st.RuntimeDuration) / float64(taskCount)))
			st.AverageTotalDuration = int64(math.Round(float64(st.TotalDuration) / float64(taskCount)))
		}
		dict[st.GetId()] = st

		if !st.LastTaskId.IsZero() {
			tids = append(tids, st.LastTaskId)
		}
	}

	// task list and stats
	var tasks []models.Task
	dictTask := map[primitive.ObjectID]models.Task{}
	dictTaskStat := map[primitive.ObjectID]models.TaskStat{}
	if len(tids) > 0 {
		// task list
		queryTask := bson.M{
			"_id": bson.M{
				"$in": tids,
			},
		}
		tasks, err = ctx.modelSvc.GetTaskList(queryTask, nil)
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}

		// task stats list
		taskStats, err := ctx.modelSvc.GetTaskStatList(queryTask, nil)
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}

		// cache task stats to dict
		for _, st := range taskStats {
			dictTaskStat[st.GetId()] = st
		}

		// cache task list to dict
		for _, t := range tasks {
			st, ok := dictTaskStat[t.GetId()]
			if ok {
				t.Stat = &st
			}
			dictTask[t.GetSpiderId()] = t
		}
	}

	// iterate list again
	var data []interface{}
	for _, d := range list.Values() {
		s := d.(*models.Spider)

		// spider stat
		st, ok := dict[s.GetId()]
		if ok {
			s.Stat = &st

			// last task
			t, ok := dictTask[s.GetId()]
			if ok {
				s.Stat.LastTask = &t
			}
		}

		// add to list
		data = append(data, *s)
	}

	// response
	HandleSuccessWithListData(c, data, total)
}

func (ctx *spiderContext) _processFileRequest(c *gin.Context, method string) (id primitive.ObjectID, payload entity.FileRequestPayload, fsSvc interfaces.SpiderFsService, err error) {
	// id
	id, err = primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// payload
	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// multipart/form-data
		payload, err = ctx._getFileRequestMultipartPayload(c)
		if err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
	} else {
		// query or application/json
		switch method {
		case http.MethodGet:
			err = c.ShouldBindQuery(&payload)
		default:
			err = c.ShouldBindJSON(&payload)
		}
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	// fs service
	fsSvc, err = newSpiderContext().syncSvc.GetFsService(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	return
}

func (ctx *spiderContext) _getFileRequestMultipartPayload(c *gin.Context) (payload entity.FileRequestPayload, err error) {
	fh, err := c.FormFile("file")
	if err != nil {
		return
	}
	f, err := fh.Open()
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, f); err != nil {
		return
	}
	payload.Path = c.PostForm("path")
	payload.Data = buf.String()
	return
}

func (ctx *spiderContext) _processActionRequest(c *gin.Context) (id primitive.ObjectID, err error) {
	// id
	id, err = primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	return
}

func (ctx *spiderContext) _upsertDataCollection(s *models.Spider) (err error) {
	if s.ColId.IsZero() {
		// validate
		if s.ColName == "" {
			return trace.TraceError(errors.ErrorControllerMissingRequestFields)
		}
		// no id
		dc, err := ctx.modelSvc.GetDataCollectionByName(s.ColName, nil)
		if err != nil {
			if err == mongo2.ErrNoDocuments {
				// not exists, add new
				dc = &models.DataCollection{Name: s.ColName}
				if err := delegate2.NewModelDelegate(dc).Add(); err != nil {
					return err
				}
			} else {
				// error
				return err
			}
		}
		s.ColId = dc.Id

		// create index
		_ = mongo.GetMongoCol(dc.Name).CreateIndex(mongo2.IndexModel{Keys: bson.M{"_tid": 1}})
	} else {
		// with id
		dc, err := ctx.modelSvc.GetDataCollectionById(s.ColId)
		if err != nil {
			return err
		}
		s.ColId = dc.Id
	}
	return nil
}

var _spiderCtx *spiderContext

func newSpiderContext() *spiderContext {
	if _spiderCtx != nil {
		return _spiderCtx
	}

	// context
	ctx := &spiderContext{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		panic(err)
	}
	if err := c.Provide(sync.NewSpiderSyncService); err != nil {
		panic(err)
	}
	if err := c.Provide(admin.NewSpiderAdminService); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		syncSvc interfaces.SpiderSyncService,
		adminSvc interfaces.SpiderAdminService,
	) {
		ctx.modelSvc = modelSvc
		ctx.syncSvc = syncSvc
		ctx.adminSvc = adminSvc
	}); err != nil {
		panic(err)
	}

	// model spider service
	ctx.modelSpiderSvc = ctx.modelSvc.GetBaseService(interfaces.ModelIdSpider)

	_spiderCtx = ctx

	return ctx
}

func newSpiderController() *spiderController {
	actions := getSpiderActions()
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdSpider, modelSvc.GetBaseService(interfaces.ModelIdSpider), actions)
	d := NewListPostActionControllerDelegate(ControllerIdSpider, modelSvc.GetBaseService(interfaces.ModelIdSpider), actions)
	ctx := newSpiderContext()

	return &spiderController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
