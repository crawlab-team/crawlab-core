package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenServiceInterface interface {
	GetModelById(id primitive.ObjectID) (res Token, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res Token, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Token, err error)
}

type tokenService struct {
	*CommonService
}

func NewTokenService() (svc *tokenService) {
	return &tokenService{
		NewCommonService(ModelIdToken),
	}
}
func (svc *tokenService) GetModelById(id primitive.ObjectID) (res Token, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *tokenService) GetModel(query bson.M, opts *mongo.FindOptions) (res Token, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *tokenService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Token, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

var TokenService *tokenService
