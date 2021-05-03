package grpc

import (
	"context"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	grpc "github.com/crawlab-team/crawlab-grpc"
)

type ModelDelegateServer struct {
	nodeSvc interfaces.NodeMasterService
	grpc.UnimplementedModelDelegateServiceServer
}

// Do and perform a RPC action of constants.Delegate
func (svr ModelDelegateServer) Do(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	obj, msg, err := NewDelegateBinder(req).BindWithDelegateMessage()
	if err != nil {
		return HandleError(err)
	}
	doc, ok := obj.(interfaces.BaseModelInterface)
	if !ok {
		return HandleError(errors.ErrorModelInvalidType)
	}
	switch msg.GetMethod() {
	case interfaces.ModelDelegateMethodAdd:
		err = doc.Add()
	case interfaces.ModelDelegateMethodSave:
		err = doc.Save()
	case interfaces.ModelDelegateMethodDelete:
		err = doc.Delete()
	case interfaces.ModelDelegateMethodGetArtifact:
		err = errors.ErrorGrpcNotAllowed
	}
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func NewModelDelegateServer(nodeSvc interfaces.NodeMasterService) (svr *ModelDelegateServer) {
	return &ModelDelegateServer{nodeSvc: nodeSvc}
}
