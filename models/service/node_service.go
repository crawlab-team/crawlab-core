package service

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeNode(d interface{}, err error) (res *models2.Node, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Node)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetNodeById(id primitive.ObjectID) (res *models2.Node, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdNode).GetById(id)
	return convertTypeNode(d, err)
}

func (svc *Service) GetNode(query bson.M, opts *mongo.FindOptions) (res *models2.Node, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdNode).Get(query, opts)
	return convertTypeNode(d, err)
}

func (svc *Service) GetNodeList(query bson.M, opts *mongo.FindOptions) (res []models2.Node, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdNode, query, opts, &res)
	return res, err
}

func (svc *Service) GetNodeByKey(key string, opts *mongo.FindOptions) (res *models2.Node, err error) {
	query := bson.M{"k": key}
	return svc.GetNode(query, opts)
}
