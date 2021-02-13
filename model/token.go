package model

import (
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
	m := NewDelegate(TokenColName, t)
	return m.Add()
}

func (t *Token) Save() (err error) {
	m := NewDelegate(TokenColName, t)
	return m.Save()
}

func (t *Token) Delete() (err error) {
	m := NewDelegate(TokenColName, t)
	return m.Delete()
}

func (t *Token) GetArtifact() (a Artifact, err error) {
	d := NewDelegate(TokenColName, t)
	return d.GetArtifact()
}

const TokenColName = "tokens"
