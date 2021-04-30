package store

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/google/wire"
)

func initializeGrpcService() (svc interfaces.GrpcService, err error) {
	wire.Build(GrpcServiceSet)
	return svc, nil
}

var GrpcService interfaces.GrpcService

var GrpcServiceSet wire.ProviderSet
