package fs

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/fs"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/go-trace"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"path/filepath"
)

type Service struct {
	// settings
	cfgPath           string
	fsPathBase        string
	workspacePathBase string
	repoPathBase      string

	// dependencies
	fsSvc    interfaces.FsService
	modelSvc service.ModelService

	// internals
	id primitive.ObjectID
	s  interfaces.Spider
	t  interfaces.Task
}

func (svc *Service) Init() (err error) {
	// task
	svc.t, err = svc.modelSvc.GetTaskById(svc.id)
	if err != nil {
		return err
	}

	// spider
	svc.s, err = svc.modelSvc.GetSpiderById(svc.t.GetSpiderId())
	if err != nil {
		return err
	}

	// fs service
	var fsOpts []fs.Option
	if svc.cfgPath != "" {
		fsOpts = append(fsOpts, fs.WithConfigPath(svc.cfgPath))
	}
	if svc.GetFsPath() != "" {
		fsOpts = append(fsOpts, fs.WithFsPath(svc.GetFsPath()))
	}
	if svc.GetWorkspacePath() != "" {
		fsOpts = append(fsOpts, fs.WithWorkspacePath(svc.GetWorkspacePath()))
	}
	if svc.repoPathBase != "" {
		fsOpts = append(fsOpts, fs.WithRepoPath(svc.GetRepoPath()))
	}

	// fs service
	svc.fsSvc, err = fs.NewFsService(fsOpts...)
	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) SetId(id primitive.ObjectID) {
	svc.id = id
}

func (svc *Service) GetFsService() (fsSvc interfaces.FsService) {
	return svc.fsSvc
}

func (svc *Service) GetWorkspacePath() (res string) {
	return filepath.Join(svc.workspacePathBase, svc.s.GetId().Hex(), svc.t.GetId().Hex())
}

func (svc *Service) SetFsPathBase(path string) {
	svc.fsPathBase = path
}

func (svc *Service) SetRepoPathBase(path string) {
	svc.repoPathBase = path
}

func (svc *Service) SetWorkspacePathBase(path string) {
	svc.workspacePathBase = path
}

func (svc *Service) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *Service) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *Service) GetFsPath() (res string) {
	return fmt.Sprintf("%s/%s", svc.fsPathBase, svc.s.GetId().Hex())
}

func (svc *Service) GetRepoPath() (res string) {
	return fmt.Sprintf("%s/%s", svc.repoPathBase, svc.id.Hex())
}

func NewTaskFsService(id primitive.ObjectID, opts ...Option) (svc2 interfaces.TaskFsService, err error) {
	// service
	svc := &Service{
		fsPathBase:        fs.DefaultFsPath,
		workspacePathBase: fs.DefaultWorkspacePath,
		repoPathBase:      fs.DefaultRepoPath,
		id:                id,
	}

	// workspace path
	if viper.GetString("fs.workspace.path") != "" {
		svc.workspacePathBase = viper.GetString("fs.workspace.path")
	}

	// repo path
	if viper.GetString("fs.repo.path") != "" {
		svc.repoPathBase = viper.GetString("fs.repo.path")
	}

	// workspace path
	if viper.GetString("fs.workspace.path") != "" {
		svc.workspacePathBase = viper.GetString("fs.workspace.path")
	}

	// repo path
	if viper.GetString("fs.repo.path") != "" {
		svc.repoPathBase = viper.GetString("fs.repo.path")
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// validate
	if svc.id.IsZero() {
		return nil, trace.TraceError(errors.ErrorSpiderMissingRequiredOption)
	}

	// model service
	svc.modelSvc, err = service.GetService()
	if err != nil {
		return nil, err
	}

	// init
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}
