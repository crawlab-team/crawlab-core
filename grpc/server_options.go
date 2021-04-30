package grpc

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	node2 "github.com/crawlab-team/crawlab-core/node"
)

type ServerOptions struct {
	NodeService *node2.Service
	Address     *entity.Address
}

func (opts *ServerOptions) FillEmpty() interfaces.Options {
	if opts.Address.IsEmpty() {
		opts.Address = entity.NewAddress(nil)
	}
	return opts
}
