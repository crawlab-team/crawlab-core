package server

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"go.uber.org/dig"
)

type PluginServer struct {
	grpc.UnimplementedPluginServiceServer

	// dependencies
	modelSvc service.ModelService
	cfgSvc   interfaces.NodeConfigService

	// internals
	server interfaces.GrpcServer
}

func (svr PluginServer) Register(ctx context.Context, in *grpc.PluginRequest) (res *grpc.Response, err error) {
	panic("implement me")
}

func (svr PluginServer) Subscribe(request *grpc.PluginRequest, stream grpc.PluginService_SubscribeServer) (err error) {
	log.Infof("master received subscribe request from plugin[%s]", request.Name)

	// finished channel
	finished := make(chan bool)

	// set subscribe
	svr.server.SetSubscribe("plugin:"+request.Name, &entity.GrpcSubscribe{
		Stream:   stream,
		Finished: finished,
	})
	ctx := stream.Context()

	log.Infof("master subscribed plugin[%s]", request.Name)

	// Keep this scope alive because once this scope exits - the stream is closed
	for {
		select {
		case <-finished:
			log.Infof("closing stream for plugin[%s]", request.Name)
			return nil
		case <-ctx.Done():
			log.Infof("plugin[%s] has disconnected", request.Name)
			return nil
		}

	}
}

func (svr PluginServer) deserialize(msg *grpc.StreamMessage) (data entity.StreamMessageTaskData, err error) {
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return data, trace.TraceError(err)
	}
	if data.TaskId.IsZero() {
		return data, trace.TraceError(errors.ErrorGrpcInvalidType)
	}
	return data, nil
}

func NewPluginServer(opts ...PluginServerOption) (res *PluginServer, err error) {
	// plugin server
	svr := &PluginServer{}

	// apply options
	for _, opt := range opts {
		opt(svr)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Provide(config.ProvideConfigService(svr.server.GetConfigPath())); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		cfgSvc interfaces.NodeConfigService,
	) {
		svr.modelSvc = modelSvc
		svr.cfgSvc = cfgSvc
	}); err != nil {
		return nil, err
	}

	return svr, nil
}

func ProvidePluginServer(server interfaces.GrpcServer, opts ...PluginServerOption) func() (res *PluginServer, err error) {
	return func() (*PluginServer, error) {
		opts = append(opts, WithServerPluginServerService(server))
		return NewPluginServer(opts...)
	}
}
