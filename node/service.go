package node

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/store"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"github.com/google/wire"
	"io/ioutil"
	"os"
	"path"
)

func init() {
	store.NodeServiceSet = wire.NewSet(provideService)
}

type Service struct {
	cfg  *Config
	opts *ServiceOptions
}

func (svc *Service) Init() (err error) {
	// check config directory path
	configDirPath := path.Dir(svc.opts.ConfigPath)
	if !utils.Exists(configDirPath) {
		if err := os.MkdirAll(configDirPath, 0763); err != nil {
			return trace.TraceError(err)
		}
	}

	if !utils.Exists(svc.opts.ConfigPath) {
		// not exists, set to default config
		// and create a config file for persistence
		svc.cfg = NewConfig(nil)
		data, err := json.Marshal(svc.cfg)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(svc.opts.ConfigPath, data, 0763); err != nil {
			return err
		}
	} else {
		// exists, read and set to config
		data, err := ioutil.ReadFile(svc.opts.ConfigPath)
		if err != nil {
			return trace.TraceError(err)
		}
		if err := json.Unmarshal(data, svc.cfg); err != nil {
			return err
		}
	}

	// register to ServiceMap
	ServiceMap.Store(svc.GetNodeKey(), svc)

	return nil
}

func (svc *Service) Reload() (err error) {
	return svc.Init()
}

func (svc *Service) GetNodeInfo() (res interfaces.Entity) {
	return &entity.NodeInfo{
		Key:      svc.GetNodeKey(),
		IsMaster: svc.IsMaster(),
	}
}

func (svc *Service) GetNodeKey() (res string) {
	return svc.cfg.Key
}

func (svc *Service) IsMaster() (res bool) {
	return svc.cfg.IsMaster
}

func NewService(opts *ServiceOptions) (svc *Service, err error) {
	if opts == nil {
		opts = &ServiceOptions{}
	}
	opts = opts.FillEmpty().(*ServiceOptions)

	svc = &Service{
		cfg:  &Config{},
		opts: opts,
	}

	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

func provideService() (svc *Service, err error) {
	return NewService(nil)
}
