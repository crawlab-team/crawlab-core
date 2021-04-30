package grpc

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/store"
	"github.com/crawlab-team/crawlab-core/utils"
)

func initGrpc() (err error) {
	if store.GrpcService, err = NewService(nil); err != nil {
		return nil
	}
	return nil
}

func InitGrpc() (err error) {
	return utils.InitModule(interfaces.ModuleIdGrpc, initGrpc)
}
