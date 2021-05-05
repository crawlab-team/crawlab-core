package service

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeJob(d interface{}, err error) (res *models2.Job, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Job)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetJobById(id primitive.ObjectID) (res *models2.Job, err error) {
	d, err := svc.NewBaseService(interfaces.ModelIdJob).GetById(id)
	return convertTypeJob(d, err)
}

func (svc *Service) GetJob(query bson.M, opts *mongo.FindOptions) (res *models2.Job, err error) {
	d, err := svc.NewBaseService(interfaces.ModelIdJob).Get(query, opts)
	return convertTypeJob(d, err)
}

func (svc *Service) GetJobList(query bson.M, opts *mongo.FindOptions) (res []models2.Job, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdJob, query, opts, &res)
	return res, err
}
