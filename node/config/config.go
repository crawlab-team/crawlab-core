package config

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/utils"
)

type Config entity.NodeInfo

type Options struct {
	Key      string
	IsMaster bool
}

var DefaultConfigOptions = &Options{
	Key:      utils.NewUUIDString(),
	IsMaster: false,
}

func NewConfig(opts *Options) (cfg *Config) {
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
