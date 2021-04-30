package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobServiceInterface interface {
	GetModelById(id primitive.ObjectID) (res Job, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res Job, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Job, err error)
}

type jobService struct {
	*CommonService
}

func NewJobService() (svc *jobService) {
	return &jobService{
		NewCommonService(interfaces.ModelIdJob),
	}
}
func (svc *jobService) GetModelById(id primitive.ObjectID) (res Job, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *jobService) GetModel(query bson.M, opts *mongo.FindOptions) (res Job, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *jobService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Job, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

var JobService *jobService
