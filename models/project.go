package models

import (
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
	m := NewDelegate(ModelColNameProject, p)
	return m.Add()
}

func (p *Project) Save() (err error) {
	m := NewDelegate(ModelColNameProject, p)
	return m.Save()
}

func (p *Project) Delete() (err error) {
	m := NewDelegate(ModelColNameProject, p)
	return m.Delete()
}

func (p *Project) GetArtifact() (a Artifact, err error) {
	m := NewDelegate(ModelColNameProject, p)
	return m.GetArtifact()
}

func (p *Project) GetId() (id primitive.ObjectID) {
	return p.Id
}
