package client

import "github.com/crawlab-team/crawlab-core/interfaces"

type ModelDelegateOption func(delegate interfaces.GrpcClientModelDelegate)

func WithConfigPath(path string) ModelDelegateOption {
	return func(d interfaces.GrpcClientModelDelegate) {
		d.SetConfigPath(path)
	}
}

type ModelServiceOption func(delegate interfaces.GrpcClientModelService)
