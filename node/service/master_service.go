package service

import (
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/grpc/server"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/utils"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"time"
)

type MasterService struct {
	modelSvc service.ModelService
	cfgSvc   interfaces.NodeConfigService
	server   interfaces.GrpcServer

	// settings variables
	cfgPath         string
	monitorInterval time.Duration
}

func (svc *MasterService) Init() (err error) {
	return svc.Register()
}

func (svc *MasterService) Start() {
	// start grpc server
	if err := svc.server.Start(); err != nil {
		panic(err)
	}

	// start monitoring worker nodes
	go svc.Monitor()

	// wait for quit signal
	svc.Wait()

	// stop
	svc.Stop()
}

func (svc *MasterService) Wait() {
	utils.DefaultWait()
}

func (svc *MasterService) Stop() {
	log.Info("worker node service has stopped")
}

func (svc *MasterService) Monitor() {
	for {
		if err := backoff.Retry(func() error {
			return svc.monitor()
		}, backoff.NewExponentialBackOff()); err != nil {
			_ = trace.TraceError(err)
		}

		time.Sleep(svc.monitorInterval)
	}
}

func (svc *MasterService) GetConfigService() (cfgSvc interfaces.NodeConfigService) {
	return svc.cfgSvc
}

func (svc *MasterService) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *MasterService) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *MasterService) SetMonitorInterval(duration time.Duration) {
	svc.monitorInterval = duration
}

func (svc *MasterService) Register() (err error) {
	nodeKey := svc.GetConfigService().GetNodeKey()
	node, err := svc.modelSvc.GetNodeByKey(nodeKey, nil)
	if err == mongo2.ErrNoDocuments {
		// not exists
		node := &models.Node{
			Key:      nodeKey,
			Name:     nodeKey,
			IsMaster: true,
			Status:   constants.NodeStatusOnline,
			Enabled:  true,
			Active:   true,
			ActiveTs: time.Now(),
		}
		nodeD := delegate.NewModelNodeDelegate(node)
		return nodeD.Add()
	} else if err == nil {
		// exists
		nodeD := delegate.NewModelNodeDelegate(node)
		return nodeD.UpdateStatusOnline()
	} else {
		// error
		return err
	}
}

func (svc *MasterService) monitor() (err error) {
	// all worker nodes
	nodes, err := svc.modelSvc.GetNodeList(bson.M{"is_master": false}, nil)
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			return nil
		}
		return trace.TraceError(err)
	}

	// error flag
	isErr := false

	// iterate all nodes
	for _, n := range nodes {
		// message stream
		stream, err := svc.server.GetSubscribe(n.GetKey())
		if err != nil {
			_ = trace.TraceError(err)
			isErr = true
			continue
		}

		// send stream message
		if err := stream.Send(&grpc.StreamMessage{
			Code:    grpc.StreamMessageCode_PING,
			NodeKey: n.GetKey(),
		}); err != nil {
			_ = trace.TraceError(err)
			isErr = true
			continue
		}

		// get message from inbound stream message channel
		inChMsg, err := svc.server.GetInboundStreamMessageChannel(n.GetKey())
		if err != nil {
			_ = trace.TraceError(err)
			isErr = true
			continue
		}
		msg := <-inChMsg

		// validate
		if msg.Code != grpc.StreamMessageCode_PING {
			_ = trace.TraceError(errors.ErrorNodeInvalidCode)
			isErr = true
			continue
		}
		var nodeInfo entity.NodeInfo
		if err := bson.Unmarshal(msg.Data, &nodeInfo); err != nil {
			_ = trace.TraceError(err)
			isErr = true
			continue
		}
		if nodeInfo.Key != n.GetKey() {
			_ = trace.TraceError(errors.ErrorNodeInvalidNodeKey)
			isErr = true
			continue
		}

		// update status
		nodeD := delegate.NewModelNodeDelegate(&n)
		if err := nodeD.UpdateStatus(true, time.Now(), constants.NodeStatusOnline); err != nil {
			_ = trace.TraceError(err)
			isErr = true
			continue
		}
	}

	if isErr {
		return trace.TraceError(errors.ErrorNodeMonitorError)
	}

	return nil
}

func NewMasterService(opts ...Option) (res interfaces.NodeMasterService, err error) {
	// master service
	svc := &MasterService{
		monitorInterval: 60 * time.Second,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Provide(config.ProvideConfigService(svc.cfgPath)); err != nil {
		return nil, err
	}
	if err := c.Provide(server.ProvideServer(svc.cfgPath)); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(cfgSvc interfaces.NodeConfigService, modelSvc service.ModelService, server interfaces.GrpcServer) {
		svc.cfgSvc = cfgSvc
		svc.modelSvc = modelSvc
		svc.server = server
	}); err != nil {
		return nil, err
	}

	// init
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

func ProvideMasterService(path string, opts ...Option) func() (interfaces.NodeMasterService, error) {
	return func() (interfaces.NodeMasterService, error) {
		opts = append(opts, WithConfigPath(path))
		return NewMasterService(opts...)
	}
}
