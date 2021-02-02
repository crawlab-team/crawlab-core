package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	Id     primitive.ObjectID `bson:"_id" json:"_id"`
	TaskId primitive.ObjectID `bson:"task_id" json:"task_id"`
	Sys    `bson:"_sys" json:"_sys"`
	BaseModelInterface
}

const JobColName = "jobs"

func (j *Job) Add() (err error) {
	d := NewDelegate(j.Id, JobColName, j, &j.Sys)
	return d.Add()
}

func (j *Job) Save() (err error) {
	d := NewDelegate(j.Id, JobColName, j, &j.Sys)
	return d.Save()
}

func (j *Job) Delete() (err error) {
	d := NewDelegate(j.Id, JobColName, j, &j.Sys)
	return d.Delete()
}
