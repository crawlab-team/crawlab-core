package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VariableServiceInterface interface {
	GetModelById(id primitive.ObjectID) (res Variable, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res Variable, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Variable, err error)
}

type variableService struct {
	*CommonService
}

func NewVariableService() (svc *variableService) {
	return &variableService{
		NewCommonService(ModelIdVariable),
	}
}
func (svc *variableService) GetModelById(id primitive.ObjectID) (res Variable, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *variableService) GetModel(query bson.M, opts *mongo.FindOptions) (res Variable, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *variableService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Variable, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

var VariableService *variableService
