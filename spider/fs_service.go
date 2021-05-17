package spider

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/fs"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/go-trace"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"sync"
)

type FsService struct {
	// settings variables
	fsPathBase        string
	workspacePathBase string
	repoPathBase      string

	// dependencies
	modelSvc service.ModelService
	fsSvc    interfaces.FsService

	// internals
	id primitive.ObjectID
	s  *models.Spider
}

func (svc *FsService) GetFsPath(id primitive.ObjectID) (res string) {
	return fmt.Sprintf("%s/%s", viper.GetString("spider.fs"), id.Hex())
}

func (svc *FsService) GetWorkspacePath(id primitive.ObjectID) (res string) {
	return fmt.Sprintf("%s/%s", viper.GetString("spider.workspace"), id.Hex())
}

func (svc *FsService) GetRepoPath(id primitive.ObjectID) (res string) {
	return fmt.Sprintf("%s/%s", viper.GetString("spider.repo"), id.Hex())
}

func (svc *FsService) SetId(id primitive.ObjectID) {
	svc.id = id
}

func (svc *FsService) SetFsPathBase(path string) {
	svc.fsPathBase = path
}

func (svc *FsService) SetWorkspacePathBase(path string) {
	svc.workspacePathBase = path
}

func (svc *FsService) SetRepoPathBase(path string) {
	svc.repoPathBase = path
}

func (svc *FsService) GetFsService() (fsSvc interfaces.FsService) {
	return svc.fsSvc
}

func NewSpiderFsService(id primitive.ObjectID, opts ...interfaces.SpiderFsOption) (svc2 interfaces.SpiderFsService, err error) {
	// service
	svc := &FsService{
		fsPathBase:        fs.DefaultFsPath,
		workspacePathBase: fs.DefaultWorkspacePath,
		repoPathBase:      fs.DefaultRepoPath,
		id:                id,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// validate
	if svc.id.IsZero() {
		return nil, trace.TraceError(errors.ErrorSpiderMissingRequiredOption)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(modelSvc service.ModelService) {
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// spider
	svc.s, err = svc.modelSvc.GetSpiderById(svc.id)
	if err != nil {
		return nil, err
	}

	// fs service
	svc.fsSvc, err = fs.NewFsService(
		fs.WithFsPath(svc.GetFsPath(svc.id)),
		fs.WithWorkspacePath(svc.GetWorkspacePath(svc.id)),
		fs.WithRepoPath(svc.GetRepoPath(svc.id)),
	)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func ProvideSpiderFsService(id primitive.ObjectID, opts ...interfaces.SpiderFsOption) func() (svc interfaces.SpiderFsService, err error) {
	return func() (svc interfaces.SpiderFsService, err error) {
		return NewSpiderFsService(id, opts...)
	}
}

var spiderFsSvcCache = sync.Map{}

func GetSpiderFsService(id primitive.ObjectID, opts ...interfaces.SpiderFsOption) (svc interfaces.SpiderFsService, err error) {
	res, ok := spiderFsSvcCache.Load(id)
	if !ok {
		svc, err = NewSpiderFsService(id, opts...)
		if err != nil {
			return nil, err
		}
		spiderFsSvcCache.Store(id, svc)
		return svc, nil
	}

	svc, ok = res.(interfaces.SpiderFsService)
	if !ok {
		return nil, trace.TraceError(errors.ErrorFsInvalidType)
	}

	return svc, nil
}
