package grpc

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/crawlab-team/crawlab-grpc"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
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
	if _, err := models.NodeService.GetModel(bson.M{"key": nodeKey}, nil); err == nil {
		return HandleError(errors.ErrorModelAlreadyExists)
	} else if err != mongo2.ErrNoDocuments {
		return HandleError(err)
	}

	// register
	node.Key = nodeKey
	if err := node.Add(); err != nil {
		return HandleError(err)
	}

	return HandleSuccessWithData(node)
}

var NodeService nodeServer
