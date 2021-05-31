package server

import (
	"context"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/task/stats"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"go.uber.org/dig"
)

type TaskServer struct {
	grpc.UnimplementedTaskServiceServer

	// dependencies
	modelSvc service.ModelService
	cfgSvc   interfaces.NodeConfigService
	statsSvc interfaces.TaskStatsService

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

func (svr TaskServer) FetchTask(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	panic("implement me")
}

func NewTaskServer(opts ...TaskServerOption) (res *TaskServer, err error) {
	// task server
	svr := &TaskServer{}

	// apply options
	for _, opt := range opts {
		opt(svr)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Provide(stats.ProvideGetTaskStatsService(svr.server.GetConfigPath())); err != nil {
		return nil, err
	}
	if err := c.Provide(config.ProvideConfigService(svr.server.GetConfigPath())); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		statsSvc interfaces.TaskStatsService,
		cfgSvc interfaces.NodeConfigService,
	) {
		svr.modelSvc = modelSvc
		svr.statsSvc = statsSvc
		svr.cfgSvc = cfgSvc
	}); err != nil {
		return nil, err
	}

	return svr, nil
}

func ProvideTaskServer(server interfaces.GrpcServer, opts ...TaskServerOption) func() (res *TaskServer, err error) {
	return func() (*TaskServer, error) {
		opts = append(opts, WithServerTaskServerService(server))
		return NewTaskServer(opts...)
	}
}
