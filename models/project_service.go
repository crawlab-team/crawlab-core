package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type projectService struct {
	*Service
}

func (svc *projectService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Project, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func NewProjectService() (svc *projectService) {
	return &projectService{NewService(ModelIdProject)}
}

var ProjectService = NewProjectService()
