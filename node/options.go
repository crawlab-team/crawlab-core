package node

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type ServiceOption func(svc interfaces.NodeService)

func WithConfigPath(path string) ServiceOption {
	return func(svc interfaces.NodeService) {
		svc.SetConfigPath(path)
	}
}
