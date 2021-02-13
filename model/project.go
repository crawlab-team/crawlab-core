package model

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
