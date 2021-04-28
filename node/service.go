package node

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"io/ioutil"
	"os"
	"path"
)

type service struct {
	cfg  *Config
	opts *ServiceOptions
}

func (svc *service) Init() (err error) {
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

	return nil
}

func (svc *service) GetNodeKey() (res string) {
	return svc.cfg.Key
}

func (svc *service) IsMaster() (res bool) {
	return svc.cfg.IsMaster
}

type ServiceOptions struct {
	ConfigPath string
}

var DefaultServiceOptions = &ServiceOptions{
	ConfigPath: defaultConfigPath,
}

func NewService(opts *ServiceOptions) (svc *service, err error) {
	if opts == nil {
		opts = DefaultServiceOptions
	}

	svc = &service{
		cfg:  &Config{},
		opts: opts,
	}

	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

var Service, _ = NewService(nil)
