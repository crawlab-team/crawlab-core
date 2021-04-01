package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserServiceInterface interface {
	GetModelById(id primitive.ObjectID) (res User, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res User, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []User, err error)
}

type userService struct {
	*CommonService
}

func NewUserService() (svc *userService) {
	return &userService{
		NewCommonService(ModelIdUser),
	}
}
func (svc *userService) GetModelById(id primitive.ObjectID) (res User, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *userService) GetModel(query bson.M, opts *mongo.FindOptions) (res User, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *userService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []User, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

var UserService = NewUserService()
