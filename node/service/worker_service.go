package service

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/grpc/client"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/utils"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"go.uber.org/dig"
	"io"
	"time"
)

type WorkerService struct {
	cfgSvc interfaces.NodeConfigService
	client interfaces.GrpcClient

	cfgPath           string
	stream            grpc.NodeService_StreamClient
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

	// subscribe to master
	svc.Subscribe()

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
	svc.Unsubscribe()
	log.Infof("worker[%s] service has stopped", svc.cfgSvc.GetNodeKey())
}

func (svc *WorkerService) Subscribe() {
	// register
	if err := backoff.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		data, err := json.Marshal(svc.cfgSvc.GetBasicNodeInfo())
		if err != nil {
			return trace.TraceError(err)
		}
		_, err = svc.client.GetNodeClient().Register(ctx, &grpc.Request{
			NodeKey: svc.cfgSvc.GetNodeKey(),
			Data:    data,
		})
		if err != nil {
			return err
		}
		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		_ = trace.TraceError(err)
		return
	}
	log.Infof("worker[%s] has been registered", svc.cfgSvc.GetNodeKey())

	// stream
	if err := backoff.Retry(func() error {
		var err error
		svc.stream, err = svc.client.GetNodeClient().Stream(context.Background())
		if err != nil {
			return trace.TraceError(err)
		}
		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		_ = trace.TraceError(err)
		return
	}

	// send connect stream message
	if err := backoff.Retry(func() error {
		return svc.stream.Send(&grpc.StreamMessage{
			Code:    grpc.StreamMessageCode_CONNECT,
			NodeKey: svc.cfgSvc.GetNodeKey(),
		})
	}, backoff.NewExponentialBackOff()); err != nil {
		_ = trace.TraceError(err)
		return
	}

	// log
	log.Infof("worker[%s] has subscribed to master", svc.cfgSvc.GetNodeKey())
}

func (svc *WorkerService) Unsubscribe() {
	if err := svc.stream.Send(&grpc.StreamMessage{
		Code:    grpc.StreamMessageCode_DISCONNECT,
		NodeKey: svc.cfgSvc.GetNodeKey(),
	}); err != nil {
		_ = trace.TraceError(err)
		return
	}
	log.Infof("worker[%s] has unsubscribed from master", svc.cfgSvc.GetNodeKey())
}

func (svc *WorkerService) Recv() {
	// worker client-side
	if svc.stream == nil {
		return
	}

	for {
		// receive message
		var msg *grpc.StreamMessage
		if err := backoff.Retry(func() (err error) {
			msg, err = svc.stream.Recv()
			if err == io.EOF {
				// no message
				return err
			} else if err != nil {
				// error
				return trace.TraceError(err)
			} else {
				// no error
				return nil
			}
		}, backoff.NewConstantBackOff(1*time.Second)); err != nil {
			return
		}

		// handle message according to code
		switch msg.Code {
		case grpc.StreamMessageCode_CONNECT:
			log.Infof("grpc stream has subscribed to master server")
			continue
		case grpc.StreamMessageCode_DISCONNECT:
			log.Infof("grpc stream has unsubscribed from master server")
			return
		case grpc.StreamMessageCode_PING:
			msg := GetStreamMessageWithData(grpc.StreamMessageCode_PING, svc.cfgSvc.GetBasicNodeInfo())
			if err := svc.stream.Send(msg); err != nil {
				_ = trace.TraceError(err)
				return
			}
		}
	}
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
