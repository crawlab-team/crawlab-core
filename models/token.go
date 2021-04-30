package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Token string             `json:"token" bson:"token"`
}

func (t *Token) Add() (err error) {
	if t.Id.IsZero() {
		t.Id = primitive.NewObjectID()
	}
	m := NewDelegate(interfaces.ModelIdToken, t)
	return m.Add()
}

func (t *Token) Save() (err error) {
	m := NewDelegate(interfaces.ModelIdToken, t)
	return m.Save()
}

func (t *Token) Delete() (err error) {
	m := NewDelegate(interfaces.ModelIdToken, t)
	return m.Delete()
}

func (t *Token) GetArtifact() (a interfaces.ModelArtifact, err error) {
	d := NewDelegate(interfaces.ModelIdToken, t)
	return d.GetArtifact()
}

func (t *Token) GetId() (id primitive.ObjectID) {
	return t.Id
}
