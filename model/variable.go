package model

import (
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
	m := NewDelegate(VariableColName, v)
	return m.Add()
}

func (v *Variable) Save() (err error) {
	m := NewDelegate(VariableColName, v)
	return m.Save()
}

func (v *Variable) Delete() (err error) {
	m := NewDelegate(VariableColName, v)
	return m.Delete()
}

func (v *Variable) GetArtifact() (a Artifact, err error) {
	d := NewDelegate(VariableColName, v)
	return d.GetArtifact()
}

const VariableColName = "variables"
