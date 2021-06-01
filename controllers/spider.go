package controllers

import (
	"bytes"
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/spider/admin"
	"github.com/crawlab-team/crawlab-core/spider/sync"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"io"
	"net/http"
	"strings"
)

var SpiderController ListActionController

var SpiderActions = []Action{
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
}

type spiderContext struct {
	syncSvc  interfaces.SpiderSyncService
	adminSvc interfaces.SpiderAdminService
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
	fsSvc, err = spiderCtx.syncSvc.GetFsService(id)
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

var spiderCtx = newSpiderContext()

func newSpiderContext() *spiderContext {
	// context
	ctx := &spiderContext{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(sync.NewSpiderSyncService); err != nil {
		panic(err)
	}
	if err := c.Provide(admin.NewSpiderAdminService); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(
		syncSvc interfaces.SpiderSyncService,
		adminSvc interfaces.SpiderAdminService,
	) {
		ctx.syncSvc = syncSvc
		ctx.adminSvc = adminSvc
	}); err != nil {
		panic(err)
	}

	return ctx
}
