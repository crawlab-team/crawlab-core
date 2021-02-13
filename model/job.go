package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	Id     primitive.ObjectID `bson:"_id" json:"_id"`
	TaskId primitive.ObjectID `bson:"task_id" json:"task_id"`
}

const JobColName = "jobs"

func (j *Job) Add() (err error) {
	if j.Id.IsZero() {
		j.Id = primitive.NewObjectID()
	}
	d := NewDelegate(JobColName, j)
	return d.Add()
}

func (j *Job) Save() (err error) {
	d := NewDelegate(JobColName, j)
	return d.Save()
}

func (j *Job) Delete() (err error) {
	d := NewDelegate(JobColName, j)
	return d.Delete()
}

func (j *Job) GetArtifact() (a Artifact, err error) {
	d := NewDelegate(JobColName, j)
	return d.GetArtifact()
}
