package server

import (
	"context"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"go.uber.org/dig"
)

type ModelBaseServiceServer struct {
	grpc.UnimplementedModelBaseServiceServer

	// dependencies
	modelSvc interfaces.ModelService
}

func (svr ModelBaseServiceServer) GetById(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		return svc.GetById(params.Id)
	})
}

func (svr ModelBaseServiceServer) Get(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		return svc.Get(params.Query, params.FindOptions)
	})
}

func (svr ModelBaseServiceServer) GetList(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		list, err := svc.GetList(params.Query, params.FindOptions)
		if err != nil {
			return nil, err
		}
		data, err := list.ToJSON()
		if err != nil {
			return nil, err
		}
		return data, nil
	})
}

func (svr ModelBaseServiceServer) DeleteById(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		err := svc.DeleteById(params.Id)
		return nil, err
	})
}

func (svr ModelBaseServiceServer) Delete(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		err := svc.Delete(params.Query)
		return nil, err
	})
}

func (svr ModelBaseServiceServer) DeleteList(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		err := svc.DeleteList(params.Query)
		return nil, err
	})
}

func (svr ModelBaseServiceServer) ForceDeleteList(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		err := svc.ForceDeleteList(params.Query)
		return nil, err
	})
}

func (svr ModelBaseServiceServer) UpdateById(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		err := svc.UpdateById(params.Id, params.Update)
		return nil, err
	})
}

func (svr ModelBaseServiceServer) Update(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		err := svc.Update(params.Query, params.Update, params.Fields)
		return nil, err
	})
}

func (svr ModelBaseServiceServer) UpdateDoc(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		err := svc.UpdateDoc(params.Query, params.Doc, params.Fields)
		return nil, err
	})
}

func (svr ModelBaseServiceServer) Insert(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		err := svc.Insert(params.Docs...)
		return nil, err
	})
}

func (svr ModelBaseServiceServer) Count(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return svr.handleRequest(req, func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error) {
		return svc.Count(params.Query)
	})
}

func (svr ModelBaseServiceServer) handleRequest(req *grpc.Request, handle handleBaseServiceRequest) (res *grpc.Response, err error) {
	params, msg, err := NewModelBaseServiceBinder(req).BindWithBaseServiceMessage()
	if err != nil {
		return HandleError(err)
	}
	svc := svr.modelSvc.GetBaseService(msg.GetModelId())
	d, err := handle(params, svc)
	if err != nil {
		return HandleError(err)
	}
	if d == nil {
		return HandleSuccess()
	}
	return HandleSuccessWithData(d)
}

type handleBaseServiceRequest func(params *entity.GrpcBaseServiceParams, svc interfaces.ModelBaseService) (interface{}, error)

func NewModelBaseServiceServer() (svr2 *ModelBaseServiceServer, err error) {
	svr := &ModelBaseServiceServer{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(modelSvc service.ModelService) {
		svr.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return svr, nil
}
