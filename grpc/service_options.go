package grpc

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type ServiceOptions struct {
	NodeServiceKey string
	Local          *entity.Address
	Remote         *entity.Address
}

func (opts *ServiceOptions) FillEmpty() interfaces.Options {
	if opts.Local == nil || opts.Local.IsEmpty() {
		opts.Local = entity.NewAddress(nil)
	}
	return opts
}
