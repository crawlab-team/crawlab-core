package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeNode(d interface{}, err error) (res *Node, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*Node)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetNodeById(id primitive.ObjectID) (res *Node, err error) {
	d, err := MustGetService(interfaces.ModelIdNode).GetById(id)
	return convertTypeNode(d, err)
}

func (svc *Service) GetNode(query bson.M, opts *mongo.FindOptions) (res *Node, err error) {
	d, err := MustGetService(interfaces.ModelIdNode).Get(query, opts)
	return convertTypeNode(d, err)
}

func (svc *Service) GetNodeList(query bson.M, opts *mongo.FindOptions) (res []Node, err error) {
	err = getListSerializeTarget(interfaces.ModelIdNode, query, opts, &res)
	return res, err
}

func (svc *Service) GetNodeByKey(key string, opts *mongo.FindOptions) (res *Node, err error) {
	query := bson.M{"key": key}
	return svc.GetNode(query, opts)
}
