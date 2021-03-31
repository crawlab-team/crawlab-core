package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type nodeService struct {
	*Service
}

func (svc *nodeService) GetById(id primitive.ObjectID) (res Node, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *nodeService) Get(query bson.M, opts *mongo.FindOptions) (res Node, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *nodeService) GetList(query bson.M, opts *mongo.FindOptions) (res []Node, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *nodeService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *nodeService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *nodeService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}

func NewNodeService() (svc *nodeService) {
	return &nodeService{NewService(ModelIdNode)}
}

var NodeService = NewNodeService()
