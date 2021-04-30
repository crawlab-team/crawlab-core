package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	Id     primitive.ObjectID `bson:"_id" json:"_id"`
	TaskId primitive.ObjectID `bson:"task_id" json:"task_id"`
}

func (j *Job) Add() (err error) {
	if j.Id.IsZero() {
		j.Id = primitive.NewObjectID()
	}
	d := NewDelegate(interfaces.ModelIdJob, j)
	return d.Add()
}

func (j *Job) Save() (err error) {
	d := NewDelegate(interfaces.ModelIdJob, j)
	return d.Save()
}

func (j *Job) Delete() (err error) {
	d := NewDelegate(interfaces.ModelIdJob, j)
	return d.Delete()
}

func (j *Job) GetArtifact() (a interfaces.ModelArtifact, err error) {
	d := NewDelegate(interfaces.ModelIdJob, j)
	return d.GetArtifact()
}

func (j *Job) GetId() (id primitive.ObjectID) {
	return j.Id
}
