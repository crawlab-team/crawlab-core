package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type SpiderRunOptions struct {
	Mode       string
	NodeIds    []primitive.ObjectID
	Param      string
	ScheduleId primitive.ObjectID
	Priority   int
}

type SpiderCloneOptions struct {
	Name string
}
