package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
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

type settingService struct {
	*Service
}

func (svc *settingService) GetById(id primitive.ObjectID) (res Setting, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *settingService) Get(query bson.M, opts *mongo.FindOptions) (res Setting, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *settingService) GetList(query bson.M, opts *mongo.FindOptions) (res []Setting, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *settingService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *settingService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *settingService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}

var SettingService = settingService{NewService(SettingColName)}
