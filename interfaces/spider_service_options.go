package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type RunOptions struct {
	Mode       string
	NodeIds    []primitive.ObjectID
	Param      string
	ScheduleId primitive.ObjectID
}

type CloneOptions struct {
	Name string
}
