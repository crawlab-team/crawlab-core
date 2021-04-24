package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Setting struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Key   string             `json:"key" bson:"key"`
	Value string             `json:"value" bson:"value"`
}

func (s *Setting) Add() (err error) {
	if s.Id.IsZero() {
		s.Id = primitive.NewObjectID()
	}
	m := NewDelegate(ModelIdSetting, s)
	return m.Add()
}

func (s *Setting) Save() (err error) {
	m := NewDelegate(ModelIdSetting, s)
	return m.Save()
}

func (s *Setting) Delete() (err error) {
	m := NewDelegate(ModelIdSetting, s)
	return m.Delete()
}

func (s *Setting) GetArtifact() (a Artifact, err error) {
	d := NewDelegate(ModelIdSetting, s)
	return d.GetArtifact()
}

func (s *Setting) GetId() (id primitive.ObjectID) {
	return s.Id
}
