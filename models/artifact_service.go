package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeArtifact(d interface{}, err error) (res *Artifact, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*Artifact)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetArtifactById(id primitive.ObjectID) (res *Artifact, err error) {
	d, err := MustGetService(interfaces.ModelIdArtifact).GetById(id)
	return convertTypeArtifact(d, err)
}

func (svc *Service) GetArtifact(query bson.M, opts *mongo.FindOptions) (res *Artifact, err error) {
	d, err := MustGetService(interfaces.ModelIdArtifact).Get(query, opts)
	return convertTypeArtifact(d, err)
}

func (svc *Service) GetArtifactList(query bson.M, opts *mongo.FindOptions) (res []Artifact, err error) {
	err = svc.GetListSerializeTarget(query, opts, &res)
	return res, err
}
