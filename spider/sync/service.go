package sync

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/spider/fs"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
)

type Service struct {
	// dependencies
	nodeCfgSvc interfaces.NodeConfigService

	// settings
	cfgPath string
}

func (svc *Service) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *Service) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *Service) SyncToFs(id primitive.ObjectID) (err error) {
	// validate node type
	if !svc.nodeCfgSvc.IsMaster() {
		return trace.TraceError(errors.ErrorSpiderForbidden)
	}

	if fsSvc, err := fs.GetSpiderFsService(id); err == nil {
		return fsSvc.GetFsService().SyncToFs()
	}
	return nil
}

func (svc *Service) SyncToWorkspace(id primitive.ObjectID) (err error) {
	if fsSvc, err := fs.GetSpiderFsService(id); err == nil {
		return fsSvc.GetFsService().SyncToWorkspace()
	}
	return nil
}

func NewSpiderSyncService(opts ...Option) (svc2 interfaces.SpiderSyncService, err error) {
	// service
	svc := &Service{}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(config.ProvideConfigService(svc.cfgPath)); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(nodeCfgSvc interfaces.NodeConfigService) {
		svc.nodeCfgSvc = nodeCfgSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}

func ProvideSpiderSyncService(path string, opts ...Option) func() (svc interfaces.SpiderSyncService, err error) {
	opts = append(opts, WithConfigPath(path))
	return func() (svc interfaces.SpiderSyncService, err error) {
		return NewSpiderSyncService(opts...)
	}
}
