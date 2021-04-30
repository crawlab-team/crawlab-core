package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Color       string             `json:"color" bson:"color"`
	Description string             `json:"description" bson:"description"`
	Col         string             `json:"col" bson:"col"`
}

func (p *Tag) Add() (err error) {
	if p.Id.IsZero() {
		p.Id = primitive.NewObjectID()
	}
	m := NewDelegate(interfaces.ModelIdTag, p)
	return m.Add()
}

func (p *Tag) Save() (err error) {
	m := NewDelegate(interfaces.ModelIdTag, p)
	return m.Save()
}

func (p *Tag) Delete() (err error) {
	m := NewDelegate(interfaces.ModelIdTag, p)
	return m.Delete()
}

func (p *Tag) GetArtifact() (a interfaces.ModelArtifact, err error) {
	m := NewDelegate(interfaces.ModelIdTag, p)
	return m.GetArtifact()
}

func (p *Tag) GetId() (id primitive.ObjectID) {
	return p.Id
}
