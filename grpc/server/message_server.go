package server

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"go.uber.org/dig"
	"io"
)

type MessageServer struct {
	grpc.UnimplementedMessageServiceServer

	// dependencies
	modelSvc service.ModelService
	cfgSvc   interfaces.NodeConfigService

	// internals
	server interfaces.GrpcServer
}

func (svr MessageServer) Connect(stream grpc.MessageService_ConnectServer) (err error) {
	finished := make(chan bool)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		switch msg.Code {
		case grpc.StreamMessageCode_CONNECT:
			svr.server.SetSubscribe(msg.Key, &entity.GrpcSubscribe{
				Stream:   stream,
				Finished: finished,
			})
		case grpc.StreamMessageCode_DISCONNECT:
			svr.server.DeleteSubscribe(msg.Key)
			return nil
		case grpc.StreamMessageCode_SEND:
			sub, err := svr.server.GetSubscribe(msg.To)
			if err != nil {
				return err
			}
			svr.redirectMessage(sub, msg)
		}
	}
}

func (svr MessageServer) redirectMessage(sub interfaces.GrpcSubscribe, msg *grpc.StreamMessage) {
	stream := sub.GetStream()
	if stream == nil {
		trace.PrintError(errors.ErrorGrpcStreamNotFound)
		return
	}
	if err := stream.Send(msg); err != nil {
		trace.PrintError(err)
		return
	}
}

func NewMessageServer(opts ...MessageServerOption) (res *MessageServer, err error) {
	// plugin server
	svr := &MessageServer{}

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

func ProvideMessageServer(server interfaces.GrpcServer, opts ...MessageServerOption) func() (res *MessageServer, err error) {
	return func() (*MessageServer, error) {
		opts = append(opts, WithServerMessageServerService(server))
		return NewMessageServer(opts...)
	}
}
