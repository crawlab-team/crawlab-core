package server

import (
	"context"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	grpc "github.com/crawlab-team/crawlab-grpc"
)

type ModelDelegateServer struct {
	grpc.UnimplementedModelDelegateServer
}

// Do and perform an RPC action of constants.Delegate
func (svr ModelDelegateServer) Do(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// bind message
	obj, msg, err := NewModelDelegateBinder(req).BindWithDelegateMessage()
	if err != nil {
		return HandleError(err)
	}

	// convert to model
	doc, ok := obj.(interfaces.Model)
	if !ok {
		return HandleError(errors.ErrorModelInvalidType)
	}

	// model delegate
	d := delegate.NewModelDelegate(doc)

	// declare artifact
	var a interfaces.ModelArtifact

	// apply method
	switch msg.GetMethod() {
	case interfaces.ModelDelegateMethodAdd:
		err = d.Add()
	case interfaces.ModelDelegateMethodSave:
		err = d.Save()
	case interfaces.ModelDelegateMethodDelete:
		err = d.Delete()
	case interfaces.ModelDelegateMethodGetArtifact:
		a, err = d.GetArtifact()
	case interfaces.ModelDelegateMethodRefresh:
		// TODO: implement
	}
	if err != nil {
		return HandleError(err)
	}

	return HandleSuccessWithData(a)
}

func NewModelDelegateServer() (svr *ModelDelegateServer) {
	return &ModelDelegateServer{}
}
