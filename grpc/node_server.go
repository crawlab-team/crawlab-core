package grpc

import (
	"context"
	"encoding/json"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type NodeServer struct {
	grpc.UnimplementedNodeServiceServer
	interfaces.GrpcNodeServer

	nodeSvc interfaces.NodeMasterService
}

// Register from handler/worker to master
func (svr NodeServer) Register(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// unmarshall data
	var nodeInfo entity.NodeInfo
	if req.Data != nil {
		if err := json.Unmarshal(req.Data, &nodeInfo); err != nil {
			return HandleError(err)
		}

		if nodeInfo.IsMaster {
			// error: cannot register master node
			return HandleError(errors.ErrorGrpcNotAllowed)
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
		return HandleError(errors.ErrorModelMissingRequiredData)
	}

	// find in db
	node, err := models.MustGetRootService().GetNodeByKey(nodeKey, nil)
	if err == nil {
		if node.Status != constants.NodeStatusUnregistered {
			// error: already exists
			return HandleError(errors.ErrorModelAlreadyExists)
		} else if node.IsMaster {
			// error: cannot register master node
			return HandleError(errors.ErrorGrpcNotAllowed)
		} else {
			// register existing
			node.Status = constants.NodeStatusRegistered
			if err := node.Save(); err != nil {
				return HandleError(err)
			}
		}
	} else if err == mongo2.ErrNoDocuments {
		// register new
		node := models.Node{
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
		if err := node.Add(); err != nil {
			return HandleError(err)
		}
	} else {
		// error
		return HandleError(err)
	}

	return HandleSuccessWithData(node)
}

// SendHeartbeat from handler/worker to master
func (svr NodeServer) SendHeartbeat(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// find in db
	node, err := models.MustGetRootService().GetNodeByKey(req.NodeKey, nil)
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			return HandleError(errors.ErrorGrpcClientNotExists)
		}
		return HandleError(err)
	}

	// validate status
	if node.Status == constants.NodeStatusUnregistered {
		return HandleError(errors.ErrorNodeUnregistered)
	}

	// update status
	node.Status = constants.NodeStatusOnline
	node.Active = true
	node.ActiveTs = time.Now()
	if err := node.Save(); err != nil {
		return HandleError(err)
	}

	return HandleSuccessWithData(node)
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
	node, err := models.MustGetRootService().GetNodeByKey(nodeKey, nil)
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
	return &NodeServer{nodeSvc: nodeSvc}
}
