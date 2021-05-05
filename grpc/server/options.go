package server

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type Option func(svr interfaces.GrpcServer)

func WithAddress(address interfaces.Address) Option {
	return func(svr interfaces.GrpcServer) {
		svr.SetAddress(address)
	}
}
