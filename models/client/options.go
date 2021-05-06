package client

import "github.com/crawlab-team/crawlab-core/interfaces"

type ModelDelegateOption func(delegate interfaces.GrpcClientModelDelegate)

type ModelServiceOption func(delegate interfaces.GrpcClientModelService)
