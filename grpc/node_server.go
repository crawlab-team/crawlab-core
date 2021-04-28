package grpc

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/models"
	node2 "github.com/crawlab-team/crawlab-core/node"
	"github.com/crawlab-team/crawlab-grpc"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type nodeServer struct {
	nodeSvc *node2.Service
	grpc.UnimplementedNodeServiceServer
}

// Register from handler/worker to master
func (svr nodeServer) Register(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	AllowMaster(svr.nodeSvc)

	// unmarshall data
	var node models.Node
	if req.Data != nil {
		if err := json.Unmarshal(req.Data, &node); err != nil {
			return HandleError(err)
		}

		if node.IsMaster {
			// error: cannot register master node
			return HandleError(errors.ErrorGrpcNotAllowed)
		}
	}

	// node key
	var nodeKey string
	if req.NodeKey != "" {
		nodeKey = req.NodeKey
	} else {
		nodeKey = node.Key
	}
	if nodeKey == "" {
		return HandleError(errors.ErrorModelMissingRequiredData)
	}

	// find in db
	nodeDb, err := models.NodeService.GetModelByKey(nodeKey, nil)
	if err == nil {
		if nodeDb.Status != constants.NodeStatusUnregistered {
			// error: already exists
			return HandleError(errors.ErrorModelAlreadyExists)
		} else if nodeDb.IsMaster {
			// error: cannot register master node
			return HandleError(errors.ErrorGrpcNotAllowed)
		} else {
			// register existing
			node = nodeDb
			node.Status = constants.NodeStatusRegistered
			if err := node.Save(); err != nil {
				return HandleError(err)
			}
		}
	} else if err == mongo2.ErrNoDocuments {
		// register new
		node.Key = nodeKey
		node.Status = constants.NodeStatusRegistered
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
func (svr nodeServer) SendHeartbeat(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	AllowMaster(svr.nodeSvc)

	// find in db
	node, err := models.NodeService.GetModelByKey(req.NodeKey, nil)
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

// Ping from master to worker
func (svr nodeServer) Ping(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	AllowWorker(svr.nodeSvc)

	return HandleSuccessWithData(svr.nodeSvc.GetNodeInfo())
}

func NewNodeServer(nodeSvc *node2.Service) (svr *nodeServer) {
	return &nodeServer{nodeSvc: nodeSvc}
}
