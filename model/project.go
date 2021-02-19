package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Tags        []string           `json:"tags" bson:"tags"`
}

func (p *Project) Add() (err error) {
	if p.Id.IsZero() {
		p.Id = primitive.NewObjectID()
	}
	m := NewDelegate(ProjectColName, p)
	return m.Add()
}

func (p *Project) Save() (err error) {
	m := NewDelegate(ProjectColName, p)
	return m.Save()
}

func (p *Project) Delete() (err error) {
	m := NewDelegate(ProjectColName, p)
	return m.Delete()
}

func (p *Project) GetArtifact() (a Artifact, err error) {
	d := NewDelegate(ProjectColName, p)
	return d.GetArtifact()
}

const ProjectColName = "projects"

type projectService struct {
	*Service
}

func (svc *projectService) GetById(id primitive.ObjectID) (res Project, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *projectService) Get(query bson.M, opts *mongo.FindOptions) (res Project, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *projectService) GetList(query bson.M, opts *mongo.FindOptions) (res []Project, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *projectService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *projectService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *projectService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}

var ProjectService = projectService{NewService(ProjectColName)}
