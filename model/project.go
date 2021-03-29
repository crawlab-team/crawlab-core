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
	m := NewDelegate(ProjectColName, p)
	return m.GetArtifact()
}

const ProjectColName = "projects"

type projectService struct {
	*Service
	PublicServiceInterface
}

func (svc *projectService) GetById(id primitive.ObjectID) (res interface{}, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *projectService) Get(query bson.M, opts *mongo.FindOptions) (res interface{}, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *projectService) GetList(query bson.M, opts *mongo.FindOptions) (res []interface{}, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *projectService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *projectService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *projectService) UpdateList(query bson.M, doc Project) (err error) {
	update := svc.getUpdate(doc)
	return svc.update(query, update)
}

func (svc *projectService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}

func (svc *projectService) getUpdate(doc Project) (update bson.M) {
	update = bson.M{}
	if doc.Name != "" {
		update["name"] = doc.Name
	}
	if doc.Description != "" {
		update["description"] = doc.Description
	}
	if doc.Tags != nil {
		update["tags"] = doc.Tags
	}
	return update
}

var ProjectService = projectService{NewService(ProjectColName), nil}
