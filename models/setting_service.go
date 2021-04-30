package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SettingServiceInterface interface {
	GetModelById(id primitive.ObjectID) (res Setting, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res Setting, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Setting, err error)
}

type settingService struct {
	*CommonService
}

func NewSettingService() (svc *settingService) {
	return &settingService{
		NewCommonService(interfaces.ModelIdSetting),
	}
}
func (svc *settingService) GetModelById(id primitive.ObjectID) (res Setting, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *settingService) GetModel(query bson.M, opts *mongo.FindOptions) (res Setting, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *settingService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Setting, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

var SettingService *settingService
