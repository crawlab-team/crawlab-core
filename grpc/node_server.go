package grpc

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/crawlab-team/crawlab-grpc"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type nodeServer struct {
	grpc.UnimplementedNodeServiceServer
}

// Register from handler/worker to master
func (svc nodeServer) Register(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// unmarshall data
	var node models.Node
	if err := json.Unmarshal(req.Data, &node); err != nil {
		return HandleError(err)
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
			// already exists
			return HandleError(errors.ErrorModelAlreadyExists)
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

// Ping or send heartbeat from handler/worker to master
func (svc nodeServer) Ping(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// node key
	nodeKey := req.NodeKey

	// find in db
	node, err := models.NodeService.GetModelByKey(nodeKey, nil)
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

var NodeService nodeServer
