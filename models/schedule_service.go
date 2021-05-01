package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeSchedule(d interface{}, err error) (res *Schedule, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*Schedule)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetScheduleById(id primitive.ObjectID) (res *Schedule, err error) {
	d, err := MustGetService(interfaces.ModelIdSchedule).GetById(id)
	return convertTypeSchedule(d, err)
}

func (svc *Service) GetSchedule(query bson.M, opts *mongo.FindOptions) (res *Schedule, err error) {
	d, err := MustGetService(interfaces.ModelIdSchedule).Get(query, opts)
	return convertTypeSchedule(d, err)
}

func (svc *Service) GetScheduleList(query bson.M, opts *mongo.FindOptions) (res []Schedule, err error) {
	err = svc.GetListSerializeTarget(query, opts, &res)
	return res, err
}
