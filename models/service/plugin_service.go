package service

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypePlugin(d interface{}, err error) (res *models2.Plugin, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Plugin)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetPluginById(id primitive.ObjectID) (res *models2.Plugin, err error) {
	d, err := svc.NewBaseService(interfaces.ModelIdPlugin).GetById(id)
	return convertTypePlugin(d, err)
}

func (svc *Service) GetPlugin(query bson.M, opts *mongo.FindOptions) (res *models2.Plugin, err error) {
	d, err := svc.NewBaseService(interfaces.ModelIdPlugin).Get(query, opts)
	return convertTypePlugin(d, err)
}

func (svc *Service) GetPluginList(query bson.M, opts *mongo.FindOptions) (res []models2.Plugin, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdPlugin, query, opts, &res)
	return res, err
}

func (svc *Service) GetPluginByName(name string) (res *models2.Plugin, err error) {
	query := bson.M{"name": name}
	return svc.GetPlugin(query, nil)
}
