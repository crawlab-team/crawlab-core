package plugin

import (
	"github.com/crawlab-team/crawlab-core/fs"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"os"
	"sync"
)

type Service struct {
	// settings variables
	dirPath    string
	fsPathBase string

	// dependencies
	modelSvc service.ModelService
	fsSvc    interfaces.FsService
}

func (svc *Service) Init() (err error) {
	// initialize directory
	if err := svc.initDir(); err != nil {
		return err
	}

	// initialize plugins
	if err := svc.initPlugins(); err != nil {
		return err
	}

	return nil
}

func (svc *Service) Start() {
	// do nothing
}

func (svc *Service) Wait() {
	utils.DefaultWait()
}

func (svc *Service) Stop() {
	// do nothing
}

func (svc *Service) SetDirPath(path string) {
	svc.dirPath = path
}

func (svc *Service) SetFsPathBase(path string) {
	svc.fsPathBase = path
}

func (svc *Service) initDir() (err error) {
	_, err = os.Stat(svc.dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(svc.dirPath, os.FileMode(0766)); err != nil {
				return trace.TraceError(err)
			}
		}
	}
	return nil
}

func (svc *Service) initPlugins() (err error) {
	for _, f := range utils.ListDir(svc.dirPath) {
		// skip non-directory
		if !f.IsDir() {
			continue
		}

		// add to db
		p := &models.Plugin{
			Name: f.Name(),
		}
		if err := svc._addPluginToDb(p); err != nil {
			return err
		}

		// sync to fs
		if err := svc.fsSvc.SyncToFs(interfaces.WithOnlyFromWorkspace()); err != nil {
			return err
		}
	}
	return nil
}

func (svc *Service) _addPluginToDb(p *models.Plugin) (err error) {
	_, err = svc.modelSvc.GetPluginByName(p.Name)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			// not exists, add new
			return delegate.NewModelDelegate(p).Add()
		}
		return err
	} else {
		// exists
		return nil
	}
}

func NewPluginService(opts ...Option) (svc2 interfaces.PluginService, err error) {
	// service
	svc := &Service{
		dirPath:    DefaultPluginDirPath,
		fsPathBase: DefaultPluginFsPathBase,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
	) {
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, err
	}

	// fs service
	var fsOpts []fs.Option
	if svc.fsPathBase != "" {
		fsOpts = append(fsOpts, fs.WithFsPath(svc.fsPathBase))
	}
	if svc.dirPath != "" {
		fsOpts = append(fsOpts, fs.WithWorkspacePath(svc.dirPath))
	}
	svc.fsSvc, err = fs.NewFsService(fsOpts...)
	if err != nil {
		return nil, err
	}

	// initialize
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

var store = sync.Map{}

func GetPluginService(path string, opts ...Option) (svc interfaces.PluginService, err error) {
	if path == "" {
		path = DefaultPluginDirPath
	}
	opts = append(opts, WithDirPath(path))
	res, ok := store.Load(path)
	if ok {
		svc, ok = res.(interfaces.PluginService)
		if ok {
			return svc, nil
		}
	}
	svc, err = NewPluginService(opts...)
	if err != nil {
		return nil, err
	}
	store.Store(path, svc)
	return svc, nil
}

func ProvideGetPluginService(path string, opts ...Option) func() (svr interfaces.PluginService, err error) {
	return func() (svr interfaces.PluginService, err error) {
		return GetPluginService(path, opts...)
	}
}
