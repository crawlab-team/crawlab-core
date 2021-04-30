package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeProject(d interface{}, err error) (res *Project, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*Project)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetProjectById(id primitive.ObjectID) (res *Project, err error) {
	d, err := NewBaseService(interfaces.ModelIdProject).GetById(id)
	return convertTypeProject(d, err)
}

func (svc *Service) GetProject(query bson.M, opts *mongo.FindOptions) (res *Project, err error) {
	d, err := NewBaseService(interfaces.ModelIdProject).Get(query, opts)
	return convertTypeProject(d, err)
}

func (svc *Service) GetProjectList(query bson.M, opts *mongo.FindOptions) (res []Project, err error) {
	err = svc.GetListSerializeTarget(query, opts, &res)
	return res, err
}
