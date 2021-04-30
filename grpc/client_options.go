package grpc

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type ClientOptions struct {
	Address        *entity.Address
	TimeoutSeconds int
}

func (opts *ClientOptions) FillEmpty() (res interfaces.Options) {
	if opts.Address == nil || opts.Address.IsEmpty() {
		opts.Address = entity.NewAddress(nil)
	}
	if opts.TimeoutSeconds == 0 {
		opts.TimeoutSeconds = 30
	}
	return opts
}
