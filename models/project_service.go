package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
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
	d, err := svc.GetById(id)
	return *d.(*Project), err
}

func (svc *projectService) GetModel(query bson.M, opts *mongo.FindOptions) (res Project, err error) {
	d, err := svc.Get(query, opts)
	return *d.(*Project), err
}

func (svc *projectService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Project, err error) {
	err = svc.GetListSerializeTarget(query, opts, &res)
	return res, err
}

func NewProjectService() (svc *projectService) {
	return &projectService{NewCommonService(interfaces.ModelIdProject)}
}

var ProjectService *projectService
