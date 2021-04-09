package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScheduleServiceInterface interface {
	GetModelById(id primitive.ObjectID) (res Schedule, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res Schedule, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Schedule, err error)
}

type scheduleService struct {
	*CommonService
}

func NewScheduleService() (svc *scheduleService) {
	return &scheduleService{
		NewCommonService(ModelIdSchedule),
	}
}
func (svc *scheduleService) GetModelById(id primitive.ObjectID) (res Schedule, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *scheduleService) GetModel(query bson.M, opts *mongo.FindOptions) (res Schedule, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *scheduleService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Schedule, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

var ScheduleService *scheduleService
