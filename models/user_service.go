package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeUser(d interface{}, err error) (res *User, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*User)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetUserById(id primitive.ObjectID) (res *User, err error) {
	d, err := MustGetService(interfaces.ModelIdUser).GetById(id)
	return convertTypeUser(d, err)
}

func (svc *Service) GetUser(query bson.M, opts *mongo.FindOptions) (res *User, err error) {
	d, err := MustGetService(interfaces.ModelIdUser).Get(query, opts)
	return convertTypeUser(d, err)
}

func (svc *Service) GetUserList(query bson.M, opts *mongo.FindOptions) (res []User, err error) {
	err = getListSerializeTarget(interfaces.ModelIdUser, query, opts, &res)
	return res, err
}

func (svc *Service) GetUserByUsername(username string, opts *mongo.FindOptions) (res *User, err error) {
	query := bson.M{"username": username}
	return svc.GetUser(query, opts)
}
