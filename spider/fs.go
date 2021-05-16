package spider

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/fs"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
)

type FsService struct {
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

func NewSpiderFsService(opts ...interfaces.FsOption) (svc2 interfaces.SpiderFsService, err error) {
	// service
	svc := &FsService{}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// validate
	if svc.id.IsZero() {
		return nil, errors.ErrorSpiderMissingRequiredOption
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Provide(fs.NewFsService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(modelSvc service.ModelService) {
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, err
	}

	// spider
	s, err := svc.modelSvc.GetSpiderById(svc.id)
	if err != nil {
		return nil, err
	}

	// fs service
	fsSvc, err := fs.NewFsService(&fs.FileSystemServiceOptions{
		IsMaster:      viper.GetBool("server.master"),
		FsPath:        svc.GetFsPath(opts.Id),
		WorkspacePath: svc.GetWorkspacePath(opts.Id),
		RepoPath:      svc.GetRepoPath(opts.Id),
	})
	if err != nil {
		return nil, err
	}

	// assign
	svc.fileSystemService = fsSvc

	return svc, nil
}

func ProvideSpiderFsService() {

}
