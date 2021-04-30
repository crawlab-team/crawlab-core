package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeSpider(d interface{}, err error) (res *Spider, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*Spider)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetSpiderById(id primitive.ObjectID) (res *Spider, err error) {
	d, err := NewBaseService(interfaces.ModelIdSpider).GetById(id)
	return convertTypeSpider(d, err)
}

func (svc *Service) GetSpider(query bson.M, opts *mongo.FindOptions) (res *Spider, err error) {
	d, err := NewBaseService(interfaces.ModelIdSpider).Get(query, opts)
	return convertTypeSpider(d, err)
}

func (svc *Service) GetSpiderList(query bson.M, opts *mongo.FindOptions) (res []Spider, err error) {
	err = svc.GetListSerializeTarget(query, opts, &res)
	return res, err
}
