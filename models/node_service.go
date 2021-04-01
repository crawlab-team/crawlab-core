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
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *nodeService) GetModel(query bson.M, opts *mongo.FindOptions) (res Node, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *nodeService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Node, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func NewNodeService() (svc *nodeService) {
	return &nodeService{NewCommonService(ModelIdNode)}
}

var NodeService = NewNodeService()
