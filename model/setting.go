package model

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
	m := NewDelegate(SettingColName, s)
	return m.Add()
}

func (s *Setting) Save() (err error) {
	m := NewDelegate(SettingColName, s)
	return m.Save()
}

func (s *Setting) Delete() (err error) {
	m := NewDelegate(SettingColName, s)
	return m.Delete()
}

func (s *Setting) GetArtifact() (a Artifact, err error) {
	d := NewDelegate(SettingColName, s)
	return d.GetArtifact()
}

const SettingColName = "settings"
