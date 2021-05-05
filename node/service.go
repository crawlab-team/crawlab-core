package node

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
)

func NewService(opts ...ServiceOption) (svc interfaces.NodeService, err error) {
	// default config service
	cfgSvc := &ConfigService{
		path: DefaultConfigPath,
	}

	// construct service given the node type
	if cfgSvc.IsMaster() {
		svc, err = NewMasterService(cfgSvc)
	} else {
		svc = NewWorkerService(cfgSvc)
	}

	// apply option
	for _, opt := range opts {
		opt(svc)
	}

	// start service
	go svc.Start()

	return svc, nil
}
