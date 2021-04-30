package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArtifactServiceInterface interface {
	GetModelById(id primitive.ObjectID) (res Artifact, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res Artifact, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Artifact, err error)
}

type artifactService struct {
	*CommonService
}

func NewArtifactService() (svc *artifactService) {
	return &artifactService{
		NewCommonService(interfaces.ModelIdArtifact),
	}
}

func (svc *artifactService) GetModelById(id primitive.ObjectID) (res Artifact, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *artifactService) GetModel(query bson.M, opts *mongo.FindOptions) (res Artifact, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *artifactService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Artifact, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

var ArtifactService *artifactService
