package node

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
)

func NewService(opts *ServiceOptions) (svc interfaces.NodeService, err error) {
	// validate options
	if opts == nil {
		opts = &ServiceOptions{}
	}
	opts = opts.FillEmpty().(*ServiceOptions)

	// config service
	cfgSvc, err := NewConfigService(opts)
	if err != nil {
		return nil, err
	}

	// construct service given the node type
	if cfgSvc.IsMaster() {
		svc = NewMasterService(cfgSvc)
	} else {
		svc = NewWorkerService(cfgSvc)
	}

	// start service
	svc.Start()

	return svc, nil
}
