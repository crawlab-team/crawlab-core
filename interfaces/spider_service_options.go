package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type SpiderRunOptions struct {
	Mode       string
	NodeIds    []primitive.ObjectID
	Param      string
	ScheduleId primitive.ObjectID
}

type SpiderCloneOptions struct {
	Name string
}
