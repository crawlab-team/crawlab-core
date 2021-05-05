package node

import (
	"fmt"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/utils"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/dig"
	"sync"
	"time"
)

type MasterService struct {
	*ConfigService

	env      string
	modelSvc service.ModelService

	chMsgMap  sync.Map
	streamMap sync.Map
}

func (svc *MasterService) Inject() (err error) {
	utils.MustResolveModule(svc.env, svc.modelSvc)
	return nil
}

func (svc *MasterService) Start() {
	go svc.Monitor()
	svc.Wait()
	svc.Stop()
}

func (svc *MasterService) Wait() {
	utils.DefaultWait()
}

func (svc *MasterService) Stop() {
	log.Info("worker node service has stopped")
}

func (svc *MasterService) getStreamMessageChannel(prefix string, key string) (chMsg chan *grpc.StreamMessage, err error) {
	_key := fmt.Sprintf("%s:%s", prefix, key)
	res, ok := svc.chMsgMap.Load(_key)
	if !ok {
		chMsg := make(chan *grpc.StreamMessage)
		svc.chMsgMap.Store(_key, chMsg)
		return chMsg, nil
	}

	chMsg, ok = res.(chan *grpc.StreamMessage)
	if !ok {
		return nil, errors.ErrorNodeInvalidType
	}
	return chMsg, nil
}

func (svc *MasterService) GetInboundStreamMessageChannel(key string) (chMsg chan *grpc.StreamMessage, err error) {
	return svc.getStreamMessageChannel("in", key)
}

func (svc *MasterService) GetOutboundStreamMessageChannel(key string) (chMsg chan *grpc.StreamMessage, err error) {
	return svc.getStreamMessageChannel("out", key)
}

func (svc *MasterService) monitor() (err error) {
	// all nodes
	nodes, err := svc.modelSvc.GetNodeList(nil, nil)
	if err != nil {
		return trace.TraceError(err)
	}

	// error flag
	isErr := false

	// iterate all nodes
	for _, n := range nodes {
		// message stream
		stream, err := svc.GetStream(n.GetKey())
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
		inChMsg, err := svc.GetInboundStreamMessageChannel(n.GetKey())
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
		if err := n.UpdateStatus(true, time.Now(), constants.NodeStatusOnline); err != nil {
			_ = trace.TraceError(err)
			isErr = true
			continue
		}
	}

	if !isErr {
		return trace.TraceError(errors.ErrorNodeMonitorError)
	}

	return nil
}

func (svc *MasterService) Monitor() {
	for {
		if err := backoff.Retry(func() error {
			return svc.monitor()
		}, backoff.NewExponentialBackOff()); err != nil {
			_ = trace.TraceError(err)
		}

		// TODO: parameterize
		time.Sleep(60 * time.Second)
	}
}

func (svc *MasterService) GetStream(key string) (stream grpc.NodeService_StreamServer, err error) {
	res, ok := svc.streamMap.Load(key)
	if !ok {
		return nil, errors.ErrorNodeStreamNotFound
	}
	stream, ok = res.(grpc.NodeService_StreamServer)
	if !ok {
		return nil, errors.ErrorNodeInvalidType
	}
	return stream, nil
}

func (svc *MasterService) SetStream(key string, stream grpc.NodeService_StreamServer) {
	svc.streamMap.Store(key, stream)
}

func (svc *MasterService) DeleteStream(key string) {
	svc.streamMap.Delete(key)
}

func NewMasterService(cfgSvs *ConfigService) (svc *MasterService, err error) {
	svc = &MasterService{
		ConfigService: cfgSvs,
		chMsgMap:      sync.Map{},
		streamMap:     sync.Map{},
	}
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(modelSvc service.ModelService) {
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, err
	}
	return svc, nil
}
