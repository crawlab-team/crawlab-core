package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
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

type tokenService struct {
	*Service
}

func (svc *tokenService) GetById(id primitive.ObjectID) (res Token, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *tokenService) Get(query bson.M, opts *mongo.FindOptions) (res Token, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *tokenService) GetList(query bson.M, opts *mongo.FindOptions) (res []Token, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *tokenService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *tokenService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *tokenService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}

var TokenService = tokenService{NewService(TokenColName)}
