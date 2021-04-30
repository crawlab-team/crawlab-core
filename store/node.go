package store

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/google/wire"
)

func initializeNodeService() (svc interfaces.NodeService, err error) {
	wire.Build(NodeServiceSet)
	return svc, nil
}

var NodeService interfaces.NodeService

var NodeServiceSet wire.ProviderSet
