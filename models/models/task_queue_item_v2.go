package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskQueueItemV2 struct {
	Id                           primitive.ObjectID `json:"_id" bson:"_id" collection:"task_queue"`
	BaseModelV2[TaskQueueItemV2] `bson:",inline"`
	Priority                     int                `json:"p" bson:"p"`
	NodeId                       primitive.ObjectID `json:"nid,omitempty" bson:"nid,omitempty"`
}
