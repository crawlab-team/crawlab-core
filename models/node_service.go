package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeServiceInterface interface {
	GetModelById(id primitive.ObjectID) (res Node, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res Node, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Node, err error)
}

type nodeService struct {
	*CommonService
}

func (svc *nodeService) GetModelById(id primitive.ObjectID) (res Node, err error) {
	d, err := svc.GetById(id)
	if err != nil {
		return res, err
	}
	return *d.(*Node), err
}

func (svc *nodeService) GetModel(query bson.M, opts *mongo.FindOptions) (res Node, err error) {
	d, err := svc.Get(query, opts)
	if err != nil {
		return res, err
	}
	return *d.(*Node), err
}

func (svc *nodeService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Node, err error) {
	err = svc.GetListSerializeTarget(query, opts, &res)
	return res, err
}

func (svc *nodeService) GetModelByKey(key string, opts *mongo.FindOptions) (res Node, err error) {
	query := bson.M{"key": key}
	return svc.GetModel(query, opts)
}

func NewNodeService() (svc *nodeService) {
	return &nodeService{NewCommonService(ModelIdNode)}
}

var NodeService *nodeService
