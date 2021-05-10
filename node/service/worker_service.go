package service

import (
	"context"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/grpc/client"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/utils"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"go.uber.org/dig"
	"time"
)

type WorkerService struct {
	cfgSvc interfaces.NodeConfigService
	client interfaces.GrpcClient

	cfgPath           string
	stream            grpc.NodeService_SubscribeClient
	heartbeatInterval time.Duration
}

func (svc *WorkerService) Init() (err error) {
	// do nothing
	return nil
}

func (svc *WorkerService) Start() {
	// start grpc client
	if err := svc.client.Start(); err != nil {
		panic(err)
	}

	// start receiving stream messages
	go svc.Recv()

	// start sending heartbeat to master
	go svc.ReportStatus()

	// wait for quit signal
	svc.Wait()

	// stop
	svc.Stop()
}

func (svc *WorkerService) Wait() {
	utils.DefaultWait()
}

func (svc *WorkerService) Stop() {
	svc.unsubscribe()
	log.Infof("worker[%s] service has stopped", svc.cfgSvc.GetNodeKey())
}

func (svc *WorkerService) Recv() {
	msgCh := svc.client.GetMessageChannel()
	for {
		msg := <-msgCh

		if err := svc.handleStreamMessage(msg); err != nil {
			continue
		}
	}
}

func (svc *WorkerService) recv() (msg *grpc.StreamMessage, err error) {
	return svc.stream.Recv()
}

func (svc *WorkerService) handleStreamMessage(msg *grpc.StreamMessage) (err error) {
	switch msg.Code {
	case grpc.StreamMessageCode_PING:
		_, err = svc.client.GetNodeClient().SendHeartbeat(context.Background(), svc.client.NewRequest(svc.cfgSvc.GetBasicNodeInfo()))
	}
	if err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (svc *WorkerService) ReportStatus() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), svc.heartbeatInterval)
		_, err := svc.client.GetNodeClient().SendHeartbeat(ctx, &grpc.Request{
			NodeKey: svc.cfgSvc.GetNodeKey(),
		})
		if err != nil {
			_ = trace.TraceError(err)
		}
		cancel()
	}
}

func (svc *WorkerService) GetConfigService() (cfgSvc interfaces.NodeConfigService) {
	return svc.cfgSvc
}

func (svc *WorkerService) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *WorkerService) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *WorkerService) SetHeartbeatInterval(duration time.Duration) {
	svc.heartbeatInterval = duration
}

func (svc *WorkerService) unsubscribe() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	svc.GetConfigService().GetBasicNodeInfo()
	_, err := svc.client.GetNodeClient().Unsubscribe(ctx, svc.client.NewRequest(svc.GetConfigService().GetBasicNodeInfo()))
	if err != nil {
		trace.PrintError(err)
	}
}

func NewWorkerService(opts ...Option) (res *WorkerService, err error) {
	svc := &WorkerService{
		heartbeatInterval: 15 * time.Second,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(config.ProvideConfigService(svc.cfgPath)); err != nil {
		return nil, err
	}
	if err := c.Provide(client.ProvideClient(svc.cfgPath)); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(cfgSvc interfaces.NodeConfigService, client interfaces.GrpcClient) {
		svc.cfgSvc = cfgSvc
		svc.client = client
	}); err != nil {
		return nil, err
	}

	// init
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

func ProvideWorkerService(path string, opts ...Option) func() (interfaces.NodeWorkerService, error) {
	return func() (interfaces.NodeWorkerService, error) {
		opts = append(opts, WithConfigPath(path))
		return NewWorkerService(opts...)
	}
}
