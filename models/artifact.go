package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Artifact struct {
	Id          primitive.ObjectID   `bson:"_id" json:"_id"`
	Col         string               `bson:"_col" json:"_col"`
	Del         bool                 `bson:"_del" json:"_del"`
	TagIds      []primitive.ObjectID `bson:"_tid" json:"_tid"`
	ArtifactSys `bson:"_sys" json:"_sys"`
	Obj         interface{} `bson:"_obj" json:"_obj"`
}

func (a *Artifact) Add() (err error) {
	if a.Id.IsZero() {
		a.Id = primitive.NewObjectID()
	}
	m := NewDelegate(interfaces.ModelIdArtifact, a)
	return m.Add()
}

func (a *Artifact) Save() (err error) {
	m := NewDelegate(interfaces.ModelIdArtifact, a)
	return m.Save()
}

func (a *Artifact) Delete() (err error) {
	m := NewDelegate(interfaces.ModelIdArtifact, a)
	return m.Delete()
}

func (a *Artifact) GetArtifact() (res interfaces.ModelArtifact, err error) {
	return nil, errors.ErrorModelNotAllowed
}

func (a *Artifact) GetId() (id primitive.ObjectID) {
	return a.Id
}

func (a *Artifact) GetTags() (res []interfaces.Tag, err error) {
	if a.TagIds == nil || len(a.TagIds) == 0 {
		return res, nil
	}
	query := bson.M{
		"_id": bson.M{
			"$in": a.TagIds,
		},
	}
	svc, err := GetRootService()
	if err != nil {
		return nil, err
	}
	tags, err := svc.GetTagList(query, nil)
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		res = append(res, &tag)
	}
	return res, nil
}

func (a *Artifact) UpdateTags(tagNames []string) (err error) {
	svc, err := GetRootService()
	if err != nil {
		return err
	}
	return svc.UpdateById(a.Id, tagNames)
}
