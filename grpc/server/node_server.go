package server

import (
	"context"
	"encoding/json"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	grpc2 "github.com/crawlab-team/crawlab-core/grpc"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type NodeServer struct {
	grpc.UnimplementedNodeServiceServer
	interfaces.GrpcNodeServer

	nodeSvc  interfaces.NodeMasterService
	modelSvc service.ModelService
}

// Register from handler/worker to master
func (svr NodeServer) Register(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// unmarshall data
	var nodeInfo entity.NodeInfo
	if req.Data != nil {
		if err := json.Unmarshal(req.Data, &nodeInfo); err != nil {
			return grpc2.HandleError(err)
		}

		if nodeInfo.IsMaster {
			// error: cannot register master node
			return grpc2.HandleError(errors.ErrorGrpcNotAllowed)
		}
	}

	// node key
	var nodeKey string
	if req.NodeKey != "" {
		nodeKey = req.NodeKey
	} else {
		nodeKey = nodeInfo.Key
	}
	if nodeKey == "" {
		return grpc2.HandleError(errors.ErrorModelMissingRequiredData)
	}

	// find in db
	node, err := svr.modelSvc.GetNodeByKey(nodeKey, nil)
	if err == nil {
		if node.Status != constants.NodeStatusUnregistered {
			// error: already exists
			return grpc2.HandleError(errors.ErrorModelAlreadyExists)
		} else if node.IsMaster {
			// error: cannot register master node
			return grpc2.HandleError(errors.ErrorGrpcNotAllowed)
		} else {
			// register existing
			node.Status = constants.NodeStatusRegistered
			nodeD := delegate.NewModelDelegate(node)
			if err := nodeD.Save(); err != nil {
				return grpc2.HandleError(err)
			}
		}
	} else if err == mongo2.ErrNoDocuments {
		// register new
		node := &models2.Node{
			Key:         nodeKey,
			Name:        nodeInfo.Name,
			Ip:          nodeInfo.Ip,
			Hostname:    nodeInfo.Hostname,
			Description: nodeInfo.Description,
			Status:      constants.NodeStatusRegistered,
		}
		if node.Name == "" {
			node.Name = nodeKey
		}
		nodeD := delegate.NewModelDelegate(node)
		if err := nodeD.Add(); err != nil {
			return grpc2.HandleError(err)
		}
	} else {
		// error
		return grpc2.HandleError(err)
	}

	return grpc2.HandleSuccessWithData(node)
}

// SendHeartbeat from handler/worker to master
func (svr NodeServer) SendHeartbeat(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// find in db
	node, err := svr.modelSvc.GetNodeByKey(req.NodeKey, nil)
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			return grpc2.HandleError(errors.ErrorGrpcClientNotExists)
		}
		return grpc2.HandleError(err)
	}

	// validate status
	if node.Status == constants.NodeStatusUnregistered {
		return grpc2.HandleError(errors.ErrorNodeUnregistered)
	}

	// update status
	nodeD := delegate.NewModelNodeDelegate(node)
	if err := nodeD.UpdateStatus(true, time.Now(), constants.NodeStatusOnline); err != nil {
		return grpc2.HandleError(err)
	}

	return grpc2.HandleSuccessWithData(node)
}

func (svr NodeServer) Stream(stream grpc.NodeService_StreamServer) (err error) {
	// master server-side
	for {
		// receive stream message
		msg, err := stream.Recv()
		if err != nil {
			return trace.TraceError(err)
		}

		switch msg.Code {
		case grpc.StreamMessageCode_CONNECT:
			err = backoff.Retry(func() error {
				// validate node status
				if !svr.IsValidNode(msg.NodeKey) {
					return errors.ErrorNodeInvalidStatus
				}

				// set stream into map
				svr.nodeSvc.SetStream(msg.NodeKey, stream)

				// send ack
				if err := stream.Send(&grpc.StreamMessage{
					Code:    grpc.StreamMessageCode_DISCONNECT,
					NodeKey: svr.nodeSvc.GetNodeKey(),
				}); err != nil {
					return trace.TraceError(err)
				}

				// start to listen and handle server-side stream msg
				outChMsg, err := svr.nodeSvc.GetOutboundStreamMessageChannel(msg.NodeKey)
				if err != nil {
					return trace.TraceError(err)
				}
				go svr.HandleSendStreamMessage(stream, outChMsg)

				return nil
			}, backoff.NewExponentialBackOff())
			if err != nil {
				return trace.TraceError(err)
			}

		case grpc.StreamMessageCode_DISCONNECT:
			// delete stream
			svr.nodeSvc.DeleteStream(msg.NodeKey)

			// send ack
			return stream.Send(&grpc.StreamMessage{
				Code:    grpc.StreamMessageCode_DISCONNECT,
				NodeKey: svr.nodeSvc.GetNodeKey(),
			})

		default:
			// send stream message to inbound channel
			inChMsg, err := svr.nodeSvc.GetInboundStreamMessageChannel(msg.NodeKey)
			if err != nil {
				_ = trace.TraceError(err)
				continue
			}
			inChMsg <- msg
		}
	}
}

func (svr NodeServer) HandleSendStreamMessage(stream grpc.NodeService_StreamServer, chMsg chan *grpc.StreamMessage) {
	for {
		msg := <-chMsg
		if err := stream.Send(msg); err != nil {
			_ = trace.TraceError(err)
			return
		}
	}
}

func (svr NodeServer) IsValidNode(nodeKey string) (res bool) {
	node, err := svr.modelSvc.GetNodeByKey(nodeKey, nil)
	if err != nil {
		return false
	}
	if node.Status != constants.NodeStatusOnline {
		return false
	}
	if !node.Active {
		return false
	}
	return true
}

func NewNodeServer(nodeSvc interfaces.NodeMasterService) (svr *NodeServer) {
	var modelSvc service.ModelService
	utils.MustResolveModule("", modelSvc)

	return &NodeServer{
		nodeSvc:  nodeSvc,
		modelSvc: modelSvc,
	}
}
