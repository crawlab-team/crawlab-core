package node

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/inject"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/utils"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"time"
)

type WorkerService struct {
	*ConfigService

	env     string
	grpcSvc interfaces.GrpcService

	stream grpc.NodeService_StreamClient
}

func (svc *WorkerService) Inject() (err error) {
	utils.MustResolveModule(svc.env, svc.grpcSvc)
	return nil
}

func (svc *WorkerService) Start() {
	svc.Subscribe()
	go svc.Recv()
	go svc.ReportStatus()
	svc.Wait()
	svc.Stop()
}

func (svc *WorkerService) Wait() {
	utils.DefaultWait()
}

func (svc *WorkerService) Stop() {
	svc.Unsubscribe()
	log.Info("worker node service has stopped")
}

func (svc *WorkerService) Subscribe() {
	// client
	c := inject.GrpcService.GetClient()

	// register
	if err := backoff.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		data, err := json.Marshal(svc.GetBasicNodeInfo())
		if err != nil {
			return trace.TraceError(err)
		}
		_, err = c.GetNodeClient().Register(ctx, &grpc.Request{
			NodeKey: svc.GetNodeKey(),
			Data:    data,
		})
		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		_ = trace.TraceError(err)
		return
	}

	// stream
	if err := backoff.Retry(func() error {
		var err error
		svc.stream, err = c.GetNodeClient().Stream(context.Background())
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
			NodeKey: svc.GetNodeKey(),
		})
	}, backoff.NewExponentialBackOff()); err != nil {
		_ = trace.TraceError(err)
		return
	}
}

func (svc *WorkerService) Unsubscribe() {
	if err := svc.stream.Send(&grpc.StreamMessage{
		Code:    grpc.StreamMessageCode_DISCONNECT,
		NodeKey: svc.GetNodeKey(),
	}); err != nil {
		_ = trace.TraceError(err)
		return
	}
}

func (svc *WorkerService) Recv() {
	// worker client-side
	if svc.stream == nil {
		return
	}

	for {
		// receive message
		msg, err := svc.stream.Recv()
		if err != nil {
			_ = trace.TraceError(err)
			return
		}

		// handle message according to code
		switch msg.Code {
		case grpc.StreamMessageCode_CONNECT:
			log.Infof("grpc stream has subscribed to master server")
		case grpc.StreamMessageCode_DISCONNECT:
			log.Infof("grpc stream has unsubscribed from master server")
			return
		case grpc.StreamMessageCode_PING:
			msg := GetStreamMessageWithData(grpc.StreamMessageCode_PING, svc.GetBasicNodeInfo())
			if err := svc.stream.Send(msg); err != nil {
				_ = trace.TraceError(err)
				return
			}
		}
	}
}

func (svc *WorkerService) ReportStatus() {
	c := inject.GrpcService.GetClient()
	for {
		// TODO: parameterize
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		_, err := c.GetNodeClient().SendHeartbeat(ctx, &grpc.Request{
			NodeKey: svc.GetNodeKey(),
		})
		if err != nil {
			_ = trace.TraceError(err)
		}
		cancel()
	}
}

func NewWorkerService(cfgSvc *ConfigService) (svc *WorkerService) {
	return &WorkerService{
		ConfigService: cfgSvc,
		stream:        nil,
	}
}
