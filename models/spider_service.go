package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderServiceInterface interface {
	GetModelById(id primitive.ObjectID) (res Spider, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res Spider, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Spider, err error)
}

type spiderService struct {
	*CommonService
}

func (svc *spiderService) GetModelById(id primitive.ObjectID) (res Spider, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *spiderService) GetModel(query bson.M, opts *mongo.FindOptions) (res Spider, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *spiderService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Spider, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func NewSpiderService() (svc *spiderService) {
	return &spiderService{NewCommonService(ModelIdSpider)}
}

var SpiderService SpiderServiceInterface = NewSpiderService()
