package grpc

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type ServerOptions struct {
	NodeService interfaces.NodeService
	Address     *entity.Address
}

func (opts *ServerOptions) FillEmpty() interfaces.Options {
	if opts.Address.IsEmpty() {
		opts.Address = entity.NewAddress(nil)
	}
	return opts
}
