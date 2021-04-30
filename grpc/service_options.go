package grpc

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node"
)

type ServiceOptions struct {
	NodeService *node.Service
	Local       *entity.Address
	Remotes     []*entity.Address
}

func (opts *ServiceOptions) FillEmpty() interfaces.Options {
	if opts.NodeService == nil {
		opts.NodeService, _ = node.NewService(nil)
	}
	if opts.Local == nil || opts.Local.IsEmpty() {
		opts.Local = entity.NewAddress(nil)
	}
	return opts
}
