package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
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

type jobService struct {
	*Service
}

func (svc *jobService) GetById(id primitive.ObjectID) (res Job, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *jobService) Get(query bson.M, opts *mongo.FindOptions) (res Job, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *jobService) GetList(query bson.M, opts *mongo.FindOptions) (res []Job, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *jobService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *jobService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

var JobService = jobService{NewService(JobColName)}
