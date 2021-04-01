package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectServiceInterface interface {
	GetModelById(id primitive.ObjectID) (res Project, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res Project, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Project, err error)
}

type projectService struct {
	*CommonService
}

func (svc *projectService) GetModelById(id primitive.ObjectID) (res Project, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *projectService) GetModel(query bson.M, opts *mongo.FindOptions) (res Project, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *projectService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Project, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func NewProjectService() (svc *projectService) {
	return &projectService{NewCommonService(ModelIdProject)}
}

var ProjectService ProjectServiceInterface = NewProjectService()
