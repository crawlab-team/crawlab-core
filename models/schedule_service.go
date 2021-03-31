package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type scheduleService struct {
	*Service
}

func (svc *scheduleService) GetById(id primitive.ObjectID) (res Schedule, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *scheduleService) Get(query bson.M, opts *mongo.FindOptions) (res Schedule, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *scheduleService) GetList(query bson.M, opts *mongo.FindOptions) (res []Schedule, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *scheduleService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *scheduleService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *scheduleService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}

func NewScheduleService() (svc *scheduleService) {
	return &scheduleService{NewService(ModelIdSchedule)}
}

var ScheduleService = NewScheduleService()
