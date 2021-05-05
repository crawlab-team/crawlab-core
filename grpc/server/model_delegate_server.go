package server

import (
	"context"
	"github.com/crawlab-team/crawlab-core/errors"
	grpc2 "github.com/crawlab-team/crawlab-core/grpc"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	grpc "github.com/crawlab-team/crawlab-grpc"
)

type ModelDelegateServer struct {
	nodeSvc interfaces.NodeMasterService
	grpc.UnimplementedModelDelegateServiceServer
}

// Do and perform a RPC action of constants.Delegate
func (svr ModelDelegateServer) Do(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// bind message
	obj, msg, err := NewModelDelegateBinder(req).BindWithDelegateMessage()
	if err != nil {
		return grpc2.HandleError(err)
	}

	// convert to model
	doc, ok := obj.(interfaces.Model)
	if !ok {
		return grpc2.HandleError(errors.ErrorModelInvalidType)
	}

	// model delegate
	d := delegate.NewModelDelegate(doc)

	// apply method
	switch msg.GetMethod() {
	case interfaces.ModelDelegateMethodAdd:
		err = d.Add()
	case interfaces.ModelDelegateMethodSave:
		err = d.Save()
	case interfaces.ModelDelegateMethodDelete:
		err = d.Delete()
	case interfaces.ModelDelegateMethodGetArtifact:
		err = errors.ErrorGrpcNotAllowed
	}
	if err != nil {
		return grpc2.HandleError(err)
	}

	return grpc2.HandleSuccess()
}

func NewModelDelegateServer(nodeSvc interfaces.NodeMasterService) (svr *ModelDelegateServer) {
	return &ModelDelegateServer{nodeSvc: nodeSvc}
}
