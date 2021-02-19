package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Node struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Status      string             `json:"status" bson:"status"`
	Ip          string             `json:"ip" bson:"ip"`
	Port        string             `json:"port" bson:"port"`
	Mac         string             `json:"mac" bson:"mac"`
	Hostname    string             `json:"hostname" bson:"hostname"`
	Description string             `json:"description" bson:"description"`
	Key         string             `json:"key" bson:"key"`
	IsMaster    bool               `json:"is_master" bson:"is_master"`
	Enabled     bool               `json:"enabled" bson:"enabled"`
	Active      bool               `json:"active" bson:"active"`

	Settings NodeSettings `json:"settings" bson:"settings"`
}

type NodeSettings struct {
	MaxRunners int `json:"max_runners" bson:"max_runners"`
}

func (n *Node) Add() (err error) {
	if n.Id.IsZero() {
		n.Id = primitive.NewObjectID()
	}
	m := NewDelegate(NodeColName, n)
	return m.Add()
}

func (n *Node) Save() (err error) {
	m := NewDelegate(NodeColName, n)
	return m.Save()
}

func (n *Node) Delete() (err error) {
	m := NewDelegate(NodeColName, n)
	return m.Delete()
}

func (n *Node) GetArtifact() (a Artifact, err error) {
	m := NewDelegate(NodeColName, n)
	return m.GetArtifact()
}

const NodeColName = "nodes"

type nodeService struct {
	*Service
}

func (svc *nodeService) GetById(id primitive.ObjectID) (res Node, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *nodeService) Get(query bson.M, opts *mongo.FindOptions) (res Node, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *nodeService) GetList(query bson.M, opts *mongo.FindOptions) (res []Node, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *nodeService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *nodeService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *nodeService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}

var NodeService = nodeService{NewService(NodeColName)}
