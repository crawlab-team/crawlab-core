package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeToken(d interface{}, err error) (res *Token, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*Token)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetTokenById(id primitive.ObjectID) (res *Token, err error) {
	d, err := NewBaseService(interfaces.ModelIdToken).GetById(id)
	return convertTypeToken(d, err)
}

func (svc *Service) GetToken(query bson.M, opts *mongo.FindOptions) (res *Token, err error) {
	d, err := NewBaseService(interfaces.ModelIdToken).Get(query, opts)
	return convertTypeToken(d, err)
}

func (svc *Service) GetTokenList(query bson.M, opts *mongo.FindOptions) (res []Token, err error) {
	err = svc.GetListSerializeTarget(query, opts, &res)
	return res, err
}
