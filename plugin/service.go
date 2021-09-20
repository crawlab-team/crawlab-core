package plugin

import (
	"encoding/json"
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	errors2 "github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

type Service struct {
	// settings variables
	fsPathBase string

	// dependencies
	modelSvc service.ModelService
}

func (svc *Service) Init() (err error) {
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

func (svc *Service) SetFsPathBase(path string) {
	svc.fsPathBase = path
}

func (svc *Service) InstallPlugin(id primitive.ObjectID) (err error) {
	// plugin
	p, err := svc.modelSvc.GetPluginById(id)
	if err != nil {
		return err
	}

	// install url type
	installUrlType := svc.getInstallUrlType(p)

	// install
	switch installUrlType {
	case constants.PluginInstallUrlTypePluginName:
		return svc.installPluginName(p)
	case constants.PluginInstallUrlTypeGithub:
		return svc.installGithub(p)
	case constants.PluginInstallUrlTypeGitee:
		return svc.installGitee(p)
	case constants.PluginInstallUrlTypeFile:
		return svc.installFile(p)
	case constants.PluginInstallUrlTypeGeneralUrl:
		return svc.installGeneralUrl(p)
	default:
		return trace.TraceError(errors2.ErrorPluginNotImplemented)
	}
}

func (svc *Service) UninstallPlugin(id primitive.ObjectID) (err error) {
	// TODO: implement
	panic("implement me")
}

func (svc *Service) RunPlugin(id primitive.ObjectID) (err error) {
	// TODO: implement
	panic("implement me")
}

func (svc *Service) StopPlugin(id primitive.ObjectID) (err error) {
	// TODO: implement
	panic("implement me")
}

func (svc *Service) getInstallUrlType(p interfaces.Plugin) (installUrlType string) {
	if p.GetName() != "" {
		return constants.PluginInstallUrlTypePluginName
	}

	url := p.GetInstallUrl()
	if strings.Contains(url, "github.com") {
		return constants.PluginInstallUrlTypeGithub
	} else if strings.Contains(url, "gitee.com") {
		return constants.PluginInstallUrlTypeGitee
	} else if strings.HasPrefix(url, "file:///") {
		return constants.PluginInstallUrlTypeFile
	} else {
		return constants.PluginInstallUrlTypeGeneralUrl
	}
}

func (svc *Service) installPluginName(p interfaces.Plugin) (err error) {
	p.SetInstallUrl(fmt.Sprintf("https://github.com/crawlab-team/plugin-%s", p.GetName()))
	return svc.installGithub(p)
}

func (svc *Service) installGithub(p interfaces.Plugin) (err error) {
	// TODO: implement
	panic("not implemented")
}

func (svc *Service) installGitee(p interfaces.Plugin) (err error) {
	// TODO: implement
	panic("not implemented")
}

func (svc *Service) installFile(p interfaces.Plugin) (err error) {
	// plugin path
	pluginPath := strings.Replace(p.GetInstallUrl(), "file://", "", 1)
	if !utils.Exists(pluginPath) {
		return trace.TraceError(errors2.ErrorPluginPathNotExists)
	}

	// plugin.json
	pluginJsonPath := filepath.Join(pluginPath, "plugin.json")
	if !utils.Exists(pluginJsonPath) {
		return trace.TraceError(errors2.ErrorPluginPluginJsonNotExists)
	}
	pluginJsonData, err := ioutil.ReadFile(pluginJsonPath)
	if err != nil {
		return trace.TraceError(err)
	}
	var _p models.Plugin
	if err := json.Unmarshal(pluginJsonData, &_p); err != nil {
		return trace.TraceError(err)
	}

	// sync to fs
	fsSvc, err := GetPluginFsService(p.GetId())
	if err != nil {
		return err
	}
	if err := fsSvc.GetFsService().SyncToFs(interfaces.WithOnlyFromWorkspace()); err != nil {
		return err
	}

	// fill plugin data and save to db
	_p.SetId(p.GetId())
	_p.SetInstallUrl(p.GetInstallUrl())
	if err := delegate.NewModelDelegate(&_p).Save(); err != nil {
		return err
	}

	return nil
}

func (svc *Service) installGeneralUrl(p interfaces.Plugin) (err error) {
	// TODO: implement
	panic("not implemented")
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

	// initialize
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

var store = sync.Map{}

func GetPluginService(path string, opts ...Option) (svc interfaces.PluginService, err error) {
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
