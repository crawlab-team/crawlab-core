package services

import (
	"github.com/crawlab-team/crawlab-core/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeServiceInterface interface {
	GetCurrentNode() (n *model.Node, err error)
	GetAllNodeIds() (ids []primitive.ObjectID, err error)
}

type NodeServiceOptions struct {
}

func NewNodeService(opts *NodeServiceOptions) (svc *nodeService, err error) {
	svc = &nodeService{}
	return svc, nil
}

type nodeService struct {
}

func (svc *nodeService) GetCurrentNode() (n *model.Node, err error) {
	panic("implement me")
}

func (svc *nodeService) GetAllNodeIds() (ids []primitive.ObjectID, err error) {
	panic("implement me")
}

var NodeService, err = NewNodeService(nil)
