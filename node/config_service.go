package node

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"io/ioutil"
	"os"
	"path"
)

type ConfigService struct {
	cfg  *Config
	path string
}

func (svc *ConfigService) Init() (err error) {
	// check config directory path
	configDirPath := path.Dir(svc.path)
	if !utils.Exists(configDirPath) {
		if err := os.MkdirAll(configDirPath, os.FileMode(0766)); err != nil {
			return trace.TraceError(err)
		}
	}

	if !utils.Exists(svc.path) {
		// not exists, set to default config
		// and create a config file for persistence
		svc.cfg = NewConfig(nil)
		data, err := json.Marshal(svc.cfg)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(svc.path, data, os.FileMode(0766)); err != nil {
			return err
		}
	} else {
		// exists, read and set to config
		data, err := ioutil.ReadFile(svc.path)
		if err != nil {
			return trace.TraceError(err)
		}
		if err := json.Unmarshal(data, svc.cfg); err != nil {
			return err
		}
	}

	return nil
}

func (svc *ConfigService) Reload() (err error) {
	return svc.Init()
}

func (svc *ConfigService) GetBasicNodeInfo() (res interfaces.Entity) {
	return &entity.NodeInfo{
		Key:      svc.GetNodeKey(),
		IsMaster: svc.IsMaster(),
	}
}

func (svc *ConfigService) GetNodeKey() (res string) {
	return svc.cfg.Key
}

func (svc *ConfigService) IsMaster() (res bool) {
	return svc.cfg.IsMaster
}

func (svc *ConfigService) SetConfigPath(path string) {
	svc.path = path
}

func NewConfigService(path string) (svc *ConfigService, err error) {
	svc = &ConfigService{
		cfg:  &Config{},
		path: path,
	}

	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}
