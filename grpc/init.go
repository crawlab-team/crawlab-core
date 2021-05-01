package grpc

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/store"
	"github.com/crawlab-team/crawlab-core/utils"
)

func initGrpc() (err error) {
	// default grpc service
	if store.GrpcService, err = NewService(nil); err != nil {
		return nil
	}

	// grpc service store
	store.GrpcServiceStore = store.NewGrpcServiceStore()

	// set default grpc service into the store
	if err = store.GrpcServiceStore.Set("", store.GrpcService); err != nil {
		return err
	}

	return nil
}

func InitGrpc() (err error) {
	return utils.InitModule(interfaces.ModuleIdGrpc, initGrpc)
}
