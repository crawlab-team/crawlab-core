package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskServiceInterface interface {
	GetModelById(id primitive.ObjectID) (res Task, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res Task, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Task, err error)
}

type taskService struct {
	*CommonService
}

func (svc *taskService) GetModelById(id primitive.ObjectID) (res Task, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *taskService) GetModel(query bson.M, opts *mongo.FindOptions) (res Task, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *taskService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Task, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func NewTaskService() (svc *taskService) {
	return &taskService{NewCommonService(ModelIdTask)}
}

var TaskService TaskServiceInterface = NewTaskService()
