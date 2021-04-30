package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeTask(d interface{}, err error) (res *Task, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*Task)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetTaskById(id primitive.ObjectID) (res *Task, err error) {
	d, err := NewBaseService(interfaces.ModelIdTask).GetById(id)
	return convertTypeTask(d, err)
}

func (svc *Service) GetTask(query bson.M, opts *mongo.FindOptions) (res *Task, err error) {
	d, err := NewBaseService(interfaces.ModelIdTask).Get(query, opts)
	return convertTypeTask(d, err)
}

func (svc *Service) GetTaskList(query bson.M, opts *mongo.FindOptions) (res []Task, err error) {
	err = svc.GetListSerializeTarget(query, opts, &res)
	return res, err
}
