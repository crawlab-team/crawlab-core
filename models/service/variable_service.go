package service

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeVariable(d interface{}, err error) (res *models2.Variable, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Variable)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetVariableById(id primitive.ObjectID) (res *models2.Variable, err error) {
	d, err := svc.NewBaseService(interfaces.ModelIdVariable).GetById(id)
	return convertTypeVariable(d, err)
}

func (svc *Service) GetVariable(query bson.M, opts *mongo.FindOptions) (res *models2.Variable, err error) {
	d, err := svc.NewBaseService(interfaces.ModelIdVariable).Get(query, opts)
	return convertTypeVariable(d, err)
}

func (svc *Service) GetVariableList(query bson.M, opts *mongo.FindOptions) (res []models2.Variable, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdVariable, query, opts, &res)
	return res, err
}

func (svc *Service) GetVariableByKey(key string, opts *mongo.FindOptions) (res *models2.Variable, err error) {
	query := bson.M{"key": key}
	return svc.GetVariable(query, opts)
}
