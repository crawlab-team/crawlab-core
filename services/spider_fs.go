package services

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderFsServiceOptions struct {
	Id primitive.ObjectID
}

func NewSpiderFsService(opts *SpiderFsServiceOptions) (svc *spiderFsService, err error) {
	// validate options
	if opts.Id.IsZero() {
		return nil, constants.ErrMissingId
	}

	// spider
	s, err := models.SpiderService.GetById(opts.Id)
	if err != nil {
		return nil, err
	}

	// spider fs service
	svc = &spiderFsService{
		s:    &s,
		opts: opts,
	}

	// fs service
	fsSvc, err := NewFileSystemService(&FileSystemServiceOptions{
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
	s    *models.Spider
	opts *SpiderFsServiceOptions
	*fileSystemService
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
