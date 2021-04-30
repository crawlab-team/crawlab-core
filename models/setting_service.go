package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeSetting(d interface{}, err error) (res *Setting, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*Setting)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetSettingById(id primitive.ObjectID) (res *Setting, err error) {
	d, err := NewBaseService(interfaces.ModelIdSetting).GetById(id)
	return convertTypeSetting(d, err)
}

func (svc *Service) GetSetting(query bson.M, opts *mongo.FindOptions) (res *Setting, err error) {
	d, err := NewBaseService(interfaces.ModelIdSetting).Get(query, opts)
	return convertTypeSetting(d, err)
}

func (svc *Service) GetSettingList(query bson.M, opts *mongo.FindOptions) (res []Setting, err error) {
	err = svc.GetListSerializeTarget(query, opts, &res)
	return res, err
}

func (svc *Service) GetSettingByKey(key string, opts *mongo.FindOptions) (res *Setting, err error) {
	query := bson.M{"key": key}
	return svc.GetSetting(query, opts)
}
