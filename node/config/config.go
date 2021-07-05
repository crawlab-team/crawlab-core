package config

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/utils"
)

type Config entity.NodeInfo

type Options struct {
	Key      string
	IsMaster bool
	AuthKey  string
}

var DefaultConfigOptions = &Options{
	Key:      utils.NewUUIDString(),
	IsMaster: false,
	AuthKey:  constants.DefaultGrpcAuthKey,
}

func NewConfig(opts *Options) (cfg *Config) {
	if opts == nil {
		opts = DefaultConfigOptions
	}
	if opts.Key == "" {
		opts.Key = utils.NewUUIDString()
	}
	if opts.AuthKey == "" {
		opts.AuthKey = constants.DefaultGrpcAuthKey
	}
	return &Config{
		Key:      opts.Key,
		IsMaster: opts.IsMaster,
		AuthKey:  opts.AuthKey,
	}
}
