package services

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/fs"
	"github.com/crawlab-team/crawlab-core/models"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderFsServiceOptions struct {
	Id primitive.ObjectID
}

func NewSpiderFsService(opts *SpiderFsServiceOptions) (svc *spiderFsService, err error) {
	// validate options
	if opts.Id.IsZero() {
		return nil, errors.ErrorModelMissingId
	}

	// spider
	s, err := models.MustGetRootService().GetSpiderById(opts.Id)
	if err != nil {
		return nil, err
	}

	// spider fs service
	svc = &spiderFsService{
		s:    s,
		opts: opts,
	}

	// fs service
	fsSvc, err := fs.NewFileSystemService(&fs.FileSystemServiceOptions{
		IsMaster:      viper.GetBool("server.master"),
		FsPath:        svc.getFsPath(opts.Id),
		WorkspacePath: svc.getWorkspacePath(opts.Id),
		RepoPath:      svc.getRepoPath(opts.Id),
	})
	if err != nil {
		return nil, err
	}

	// assign
	svc.fileSystemService = fsSvc

	return svc, nil
}

type spiderFsService struct {
	s    *models2.Spider
	opts *SpiderFsServiceOptions
	*fs.fileSystemService
}

func (svc *spiderFsService) getFsPath(id primitive.ObjectID) (res string) {
	return fmt.Sprintf("%s/%s", viper.GetString("spider.fs"), id.Hex())
}

func (svc *spiderFsService) getWorkspacePath(id primitive.ObjectID) (res string) {
	return fmt.Sprintf("%s/%s", viper.GetString("spider.workspace"), id.Hex())
}

func (svc *spiderFsService) getRepoPath(id primitive.ObjectID) (res string) {
	return fmt.Sprintf("%s/%s", viper.GetString("spider.repo"), id.Hex())
}
