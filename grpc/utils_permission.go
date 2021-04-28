package grpc

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/node"
)

func AllowMaster(nodeSvc *node.Service) {
	if nodeSvc == nil {
		var err error
		nodeSvc, err = node.GetDefaultService()
		if err != nil {
			panic(err)
		}
	}
	if !nodeSvc.IsMaster() {
		panic(errors.ErrorGrpcNotAllowed)
	}
}

func AllowWorker(nodeSvc *node.Service) {
	if nodeSvc == nil {
		var err error
		nodeSvc, err = node.GetDefaultService()
		if err != nil {
			panic(err)
		}
	}
	if nodeSvc.IsMaster() {
		panic(errors.ErrorGrpcNotAllowed)
	}
}
