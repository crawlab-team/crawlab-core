package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Artifact struct {
	Id          primitive.ObjectID   `bson:"_id" json:"_id"`
	Col         string               `bson:"_col" json:"_col"`
	Del         bool                 `bson:"_del" json:"_del"`
	TagIds      []primitive.ObjectID `bson:"_tid" json:"_tid"`
	ArtifactSys `bson:"_sys" json:"_sys"`
	Obj         interface{} `bson:"_obj" json:"_obj"`
}

type ArtifactSys struct {
	CreateTs  time.Time          `json:"create_ts" bson:"create_ts"`
	CreateUid primitive.ObjectID `json:"create_uid" bson:"create_uid"`
	UpdateTs  time.Time          `json:"update_ts" bson:"update_ts"`
	UpdateUid primitive.ObjectID `json:"update_uid" bson:"update_uid"`
	DeleteTs  time.Time          `json:"delete_ts" bson:"delete_ts"`
	DeleteUid primitive.ObjectID `json:"delete_uid" bson:"delete_uid"`
}

func (a *Artifact) GetTags() (res interface{}, err error) {
	if a.TagIds == nil || len(a.TagIds) == 0 {
		return res, nil
	}
	query := bson.M{
		"_id": bson.M{
			"$in": a.TagIds,
		},
	}
	return TagService.GetModelList(query, nil)
}

func (a *Artifact) UpdateTags(tagNames []string) (err error) {
	return TagService.UpdateById(a.Id, tagNames)
}
