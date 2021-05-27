package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskQueueItem struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Priority int                `json:"p" bson:"p"`
}

func (t *TaskQueueItem) GetId() (id primitive.ObjectID) {
	return t.Id
}

func (t *TaskQueueItem) SetId(id primitive.ObjectID) {
	t.Id = id
}
