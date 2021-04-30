package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
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
	m := NewDelegate(interfaces.ModelIdSetting, s)
	return m.Add()
}

func (s *Setting) Save() (err error) {
	m := NewDelegate(interfaces.ModelIdSetting, s)
	return m.Save()
}

func (s *Setting) Delete() (err error) {
	m := NewDelegate(interfaces.ModelIdSetting, s)
	return m.Delete()
}

func (s *Setting) GetArtifact() (a interfaces.ModelArtifact, err error) {
	d := NewDelegate(interfaces.ModelIdSetting, s)
	return d.GetArtifact()
}

func (s *Setting) GetId() (id primitive.ObjectID) {
	return s.Id
}
