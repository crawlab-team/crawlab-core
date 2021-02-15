package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
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

type variableService struct {
	*Service
}

func (svc *variableService) GetById(id primitive.ObjectID) (res Variable, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *variableService) Get(query bson.M, opts *mongo.FindOptions) (res Variable, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *variableService) GetList(query bson.M, opts *mongo.FindOptions) (res []Variable, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *variableService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *variableService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

var VariableService = variableService{NewService(VariableColName)}
