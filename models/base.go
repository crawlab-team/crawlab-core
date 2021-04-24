package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ModelIdArtifact = iota
	ModelIdTag
	ModelIdNode
	ModelIdProject
	ModelIdSpider
	ModelIdTask
	ModelIdJob
	ModelIdSchedule
	ModelIdUser
	ModelIdSetting
	ModelIdToken
	ModelIdVariable
)

type ModelId int

const (
	ModelColNameArtifact = "artifacts"
	ModelColNameTag      = "tags"
	ModelColNameNode     = "nodes"
	ModelColNameProject  = "projects"
	ModelColNameSpider   = "spiders"
	ModelColNameTask     = "tasks"
	ModelColNameJob      = "jobs"
	ModelColNameSchedule = "schedules"
	ModelColNameUser     = "users"
	ModelColNameSetting  = "settings"
	ModelColNameToken    = "tokens"
	ModelColNameVariable = "variables"
)

type BaseModelInterface interface {
	Add() (err error)
	Save() (err error)
	Delete() (err error)
	GetArtifact() (a Artifact, err error)
	GetId() (id primitive.ObjectID)
}

type BaseModelWithTagsInterface interface {
	BaseModelInterface
	SetTags(tags []Tag)
}

type BaseModel struct {
	Id   primitive.ObjectID `json:"_id" bson:"_id"`
	Tags []Tag              `json:"tags" bson:"-"`
}
