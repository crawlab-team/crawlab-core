package server

import (
	"context"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"go.uber.org/dig"
)

type TaskServer struct {
	grpc.UnimplementedTaskServiceServer

	// dependencies
	modelSvc service.ModelService

	// internals
	server interfaces.GrpcServer
}

func (svr TaskServer) GetTaskInfo(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	panic("implement me")
}

func (svr TaskServer) SaveItem(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	panic("implement me")
}

func (svr TaskServer) SaveItems(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	panic("implement me")
}

func NewTaskServer() (res *TaskServer, err error) {
	// Task server
	svr := &TaskServer{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(modelSvc service.ModelService) {
		svr.modelSvc = modelSvc
	}); err != nil {
		return nil, err
	}

	return svr, nil
}
