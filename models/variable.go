package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Variable struct {
	Id     primitive.ObjectID `json:"_id" bson:"_id"`
	Key    string             `json:"key" bson:"key"`
	Value  string             `json:"value" bson:"value"`
	Remark string             `json:"remark" bson:"remark"`
}

func (v *Variable) Add() (err error) {
	if v.Id.IsZero() {
		v.Id = primitive.NewObjectID()
	}
	m := NewDelegate(interfaces.ModelIdVariable, v)
	return m.Add()
}

func (v *Variable) Save() (err error) {
	m := NewDelegate(interfaces.ModelIdVariable, v)
	return m.Save()
}

func (v *Variable) Delete() (err error) {
	m := NewDelegate(interfaces.ModelIdVariable, v)
	return m.Delete()
}

func (v *Variable) GetArtifact() (a interfaces.ModelArtifact, err error) {
	d := NewDelegate(interfaces.ModelIdVariable, v)
	return d.GetArtifact()
}

func (v *Variable) GetId() (id primitive.ObjectID) {
	return v.Id
}
