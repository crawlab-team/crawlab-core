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
	vcs "github.com/crawlab-team/crawlab-vcs"
	"github.com/crawlab-team/go-trace"
	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"path"
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
		{
			Method:      http.MethodGet,
			Path:        "/:id/git",
			HandlerFunc: spiderCtx.getGit,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id/git/remote-refs",
			HandlerFunc: spiderCtx.getGitRemoteRefs,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/git/pull",
			HandlerFunc: spiderCtx.gitPull,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/git/commit",
			HandlerFunc: spiderCtx.gitCommit,
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
	// spider id
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	// options
	var opts interfaces.SpiderRunOptions
	if err := c.ShouldBindJSON(&opts); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// user
	if u := GetUserFromContext(c); u != nil {
		opts.UserId = u.GetId()
	}

	// schedule
	if err := ctx.adminSvc.Schedule(id, &opts); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) getGit(c *gin.Context) {
	// spider id
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	// spider fs service
	fsSvc, err := ctx.syncSvc.GetFsService(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// sync from remote to workspace
	if err := fsSvc.GetFsService().SyncToWorkspace(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// git client
	gitClient, err := ctx._getGitClient(id, fsSvc)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// current branch
	currentBranch, err := ctx._getCurrentBranch(gitClient)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// branches
	branches, err := gitClient.GetBranches()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if branches == nil || len(branches) == 0 && currentBranch != "" {
		branches = []vcs.GitRef{{Name: currentBranch}}
	}

	// tags
	tags, err := gitClient.GetTags()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// changes
	changes, err := gitClient.GetStatus()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// logs
	logs, err := gitClient.GetLogsWithRefs()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// ignore
	ignore, err := ctx._getGitIgnore(fsSvc)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// git
	_git, err := ctx.modelSvc.GetGitById(id)
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	// response
	res := bson.M{
		"current_branch": currentBranch,
		"branches":       branches,
		"tags":           tags,
		"changes":        changes,
		"logs":           logs,
		"ignore":         ignore,
		"git":            _git,
	}

	HandleSuccessWithData(c, res)
}

func (ctx *spiderContext) getGitRemoteRefs(c *gin.Context) {
	// spider id
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	// remote name
	remoteName := c.Query("remote")
	if remoteName == "" {
		remoteName = vcs.GitRemoteNameUpstream
	}

	// spider fs service
	fsSvc, err := ctx.syncSvc.GetFsService(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// sync from remote to workspace
	//if err := fsSvc.GetFsService().SyncToWorkspace(); err != nil {
	//	HandleErrorInternalServerError(c, err)
	//	return
	//}

	// git client
	gitClient, err := ctx._getGitClient(id, fsSvc)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// refs
	refs, err := gitClient.GetRemoteRefs(remoteName)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, refs)
}

func (ctx *spiderContext) gitPull(c *gin.Context) {
	// payload
	var payload entity.GitPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// spider id
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	// spider fs service
	fsSvc, err := ctx.syncSvc.GetFsService(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// git
	g, err := ctx.modelSvc.GetGitById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// git client
	gitClient, err := ctx._getGitClient(id, fsSvc)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// set remote
	r, err := gitClient.GetRemote(constants.GitRemoteNameUpstream)
	if err != nil {
		if err == git.ErrRemoteNotFound {
			// create upstream remote if not exists
			r, err = ctx._createGitRemote(gitClient, g)
			if err != nil {
				HandleErrorInternalServerError(c, err)
				return
			}
		} else {
			// error
			HandleErrorInternalServerError(c, err)
			return
		}
	} else {
		// re-create upstream remote if remote urls not matched
		if len(r.Config().URLs) == 0 || r.Config().URLs[0] != g.Url {
			// delete existing upstream remote
			if err := gitClient.DeleteRemote(constants.GitRemoteNameUpstream); err != nil {
				HandleErrorInternalServerError(c, err)
				return
			}

			// create upstream remote if not exists
			r, err = ctx._createGitRemote(gitClient, g)
			if err != nil {
				HandleErrorInternalServerError(c, err)
				return
			}
		}
	}

	// branch to pull
	var branch string
	if payload.Branch == "" {
		// by default current branch
		branch, err = gitClient.GetCurrentBranch()
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	} else {
		// payload branch
		branch = payload.Branch
	}

	// attempt to pull with target branch
	if err := ctx._gitPull(gitClient, constants.GitRemoteNameUpstream, branch); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// sync to fs
	if err := fsSvc.GetFsService().SyncToFs(interfaces.WithOnlyFromWorkspace()); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) gitCommit(c *gin.Context) {
	// payload
	var payload entity.GitPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// spider id
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	// spider fs service
	fsSvc, err := ctx.syncSvc.GetFsService(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// sync from remote to workspace
	if err := fsSvc.GetFsService().SyncToWorkspace(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// git client
	gitClient, err := ctx._getGitClient(id, fsSvc)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// add
	for _, p := range payload.Paths {
		if err := gitClient.Add(p); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	// commit
	if err := gitClient.Commit(payload.CommitMessage); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// push
	if err := gitClient.Push(
		vcs.WithRemoteNamePush(vcs.GitRemoteNameUpstream),
	); err != nil {
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
	if err := ctx._upsertDataCollection(c, s); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// save
	if err := delegate2.NewModelDelegate(s, GetUserFromContext(c)).Save(); err != nil {
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
	if err := ctx._upsertDataCollection(c, s); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// add
	if err := delegate2.NewModelDelegate(s, GetUserFromContext(c)).Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// add stat
	st := &models.SpiderStat{
		Id: s.GetId(),
	}
	if err := delegate2.NewModelDelegate(st, GetUserFromContext(c)).Add(); err != nil {
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

func (ctx *spiderContext) _upsertDataCollection(c *gin.Context, s *models.Spider) (err error) {
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
				if err := delegate2.NewModelDelegate(dc, GetUserFromContext(c)).Add(); err != nil {
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

func (ctx *spiderContext) _getGitIgnore(fsSvc interfaces.SpiderFsService) (ignore []string, err error) {
	filePath := path.Join(fsSvc.GetWorkspacePath(), ".gitignore")
	if !utils.Exists(filePath) {
		return nil, nil
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	ignore = strings.Split(string(data), "\n")
	return ignore, nil
}

func (ctx *spiderContext) _createGitRemote(gitClient *vcs.GitClient, g *models.Git) (r *git.Remote, err error) {
	r, err = gitClient.CreateRemote(&config.RemoteConfig{
		Name: constants.GitRemoteNameUpstream,
		URLs: []string{g.Url},
	})
	if err != nil {
		return nil, trace.TraceError(err)
	}
	return r, nil
}

func (ctx *spiderContext) _gitPull(gitClient *vcs.GitClient, remote, branch string) (err error) {
	// remote refs
	remoteRefs, err := gitClient.GetRemoteRefs(remote)
	if err != nil {
		return err
	}

	// ref
	var ref *plumbing.Reference
	for _, remoteRef := range remoteRefs {
		if remoteRef.Type == vcs.GitRefTypeBranch && remoteRef.Name == branch {
			ref = plumbing.NewHashReference(plumbing.NewBranchReferenceName(branch), plumbing.NewHash(remoteRef.Hash))
			break
		}
	}

	// reset
	_ = gitClient.Reset()

	// checkout to target branch
	_ = gitClient.CheckoutBranchWithRemoteFromRef(branch, remote, ref)

	// pull
	if err := gitClient.Pull(
		vcs.WithRemoteNamePull(remote),
		vcs.WithBranchNamePull(branch),
	); err != nil {
		return trace.TraceError(err)
	}

	// reset
	_ = gitClient.Reset()

	return nil
}

func (ctx *spiderContext) _getGitClient(id primitive.ObjectID, fsSvc interfaces.SpiderFsService) (gitClient *vcs.GitClient, err error) {
	// auth type
	authType := ""

	// git
	g, err := ctx.modelSvc.GetGitById(id)
	if err != nil {
		if err != mongo2.ErrNoDocuments {
			return nil, trace.TraceError(err)
		}
	} else {
		authType = g.AuthType
	}

	// git client
	gitClient = fsSvc.GetFsService().GetGitClient()

	switch authType {
	case constants.GitAuthTypeHttp:
		gitClient.SetAuthType(vcs.GitAuthTypeHTTP)
		gitClient.SetUsername(g.Username)
		gitClient.SetPassword(g.Password)
	case constants.GitAuthTypeSsh:
		gitClient.SetAuthType(vcs.GitAuthTypeSSH)
		gitClient.SetUsername(g.Username)
		gitClient.SetPrivateKey(g.Password)
	default:
		return gitClient, nil
	}

	return gitClient, nil
}

func (ctx *spiderContext) _getCurrentBranch(gitClient *vcs.GitClient) (currentBranch string, err error) {
	// current branch from repo
	currentBranch, err = gitClient.GetCurrentBranch()
	if err != nil {
		return "", err
	}

	// remote refs
	remoteRefs, err := gitClient.GetRemoteRefs(constants.GitRemoteNameUpstream)
	if err != nil {
		return currentBranch, err
	}

	// return if current branch is not default branch (master)
	if currentBranch != vcs.GitDefaultBranchName {
		return currentBranch, nil
	}

	// iterate remote refs and determine current branch
	var hasMain bool
	var mainRef *plumbing.Reference
	for _, remoteRef := range remoteRefs {
		// skip non-branch remote ref
		if remoteRef.Type != vcs.GitRefTypeBranch {
			continue
		}

		// return if current branch is in remote refs
		if remoteRef.Name == currentBranch {
			return currentBranch, nil
		}

		// set has main if any
		if remoteRef.Name == vcs.GitBranchNameMain {
			hasMain = true
			mainRef = plumbing.NewHashReference(plumbing.NewBranchReferenceName(remoteRef.Name), plumbing.NewHash(remoteRef.Hash))
		}
	}

	// error if no main branch found in remote refs
	if !hasMain {
		return "", trace.TraceError(errors.ErrorGitNoMainBranch)
	}

	// checkout to main branch
	if err := gitClient.CheckoutBranchWithRemoteFromRef(vcs.GitBranchNameMain, constants.GitRemoteNameUpstream, mainRef); err != nil {
		return "", trace.TraceError(err)
	}

	return vcs.GitBranchNameMain, nil
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
