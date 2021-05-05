package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	Id     primitive.ObjectID `bson:"_id" json:"_id"`
	TaskId primitive.ObjectID `bson:"task_id" json:"task_id"`
}

func (j *Job) GetId() (id primitive.ObjectID) {
	return j.Id
}

func (j *Job) SetId(id primitive.ObjectID) {
	j.Id = id
}
