package node

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/utils"
)

type Config entity.NodeInfo

type ConfigOptions struct {
	Key      string
	IsMaster bool
}

var DefaultConfigOptions = &ConfigOptions{
	Key:      utils.NewUUIDString(),
	IsMaster: false,
}

func NewConfig(opts *ConfigOptions) (cfg *Config) {
	if opts == nil {
		opts = DefaultConfigOptions
	}
	if opts.Key == "" {
		opts.Key = utils.NewUUIDString()
	}
	return &Config{
		Key:      opts.Key,
		IsMaster: opts.IsMaster,
	}
}
