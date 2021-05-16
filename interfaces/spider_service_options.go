package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type ServiceOptions struct {
}

type RunOptions struct {
	Mode       string
	NodeIds    []primitive.ObjectID
	Param      string
	ScheduleId primitive.ObjectID
}

type CloneOptions struct {
	Name string
}

type FsOption func(fsSvc SpiderFsService)
