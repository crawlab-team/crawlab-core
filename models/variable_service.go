package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeVariable(d interface{}, err error) (res *Variable, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*Variable)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetVariableById(id primitive.ObjectID) (res *Variable, err error) {
	d, err := MustGetService(interfaces.ModelIdVariable).GetById(id)
	return convertTypeVariable(d, err)
}

func (svc *Service) GetVariable(query bson.M, opts *mongo.FindOptions) (res *Variable, err error) {
	d, err := MustGetService(interfaces.ModelIdVariable).Get(query, opts)
	return convertTypeVariable(d, err)
}

func (svc *Service) GetVariableList(query bson.M, opts *mongo.FindOptions) (res []Variable, err error) {
	err = svc.GetListSerializeTarget(query, opts, &res)
	return res, err
}

func (svc *Service) GetVariableByKey(key string, opts *mongo.FindOptions) (res *Variable, err error) {
	query := bson.M{"key": key}
	return svc.GetVariable(query, opts)
}
