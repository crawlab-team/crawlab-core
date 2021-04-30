package node

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type ServiceOptions struct {
	ConfigPath string
}

func (opts *ServiceOptions) FillEmpty() interfaces.Options {
	if opts.ConfigPath == "" {
		opts.ConfigPath = DefaultConfigPath
	}
	return opts
}
