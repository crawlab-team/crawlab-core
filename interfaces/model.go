package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModelInterface interface {
	Add() (err error)
	Save() (err error)
	Delete() (err error)
	GetArtifact() (a ModelArtifact, err error)
	GetId() (id primitive.ObjectID)
}

type ModelId int

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

type BaseModelWithTagsInterface interface {
	BaseModelInterface
	SetTags(tags []Tag)
}
