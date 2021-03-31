package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Node struct {
	Id   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
	//Status      string             `json:"status" bson:"status"`
	Ip          string `json:"ip" bson:"ip"`
	Port        string `json:"port" bson:"port"`
	Mac         string `json:"mac" bson:"mac"`
	Hostname    string `json:"hostname" bson:"hostname"`
	Description string `json:"description" bson:"description"`
	Key         string `json:"key" bson:"key"`
	IsMaster    bool   `json:"is_master" bson:"is_master"`
	Enabled     bool   `json:"enabled" bson:"enabled"`
	Active      bool   `json:"active" bson:"active"`
	Tag         string `json:"tag" bson:"tag"`

	Settings NodeSettings `json:"settings" bson:"settings"`
}

type NodeSettings struct {
	MaxRunners int `json:"max_runners" bson:"max_runners"`
}

func (n *Node) Add() (err error) {
	if n.Id.IsZero() {
		n.Id = primitive.NewObjectID()
	}
	m := NewDelegate(ModelColNameNode, n)
	return m.Add()
}

func (n *Node) Save() (err error) {
	m := NewDelegate(ModelColNameNode, n)
	return m.Save()
}

func (n *Node) Delete() (err error) {
	m := NewDelegate(ModelColNameNode, n)
	return m.Delete()
}

func (n *Node) GetArtifact() (a Artifact, err error) {
	m := NewDelegate(ModelColNameNode, n)
	return m.GetArtifact()
}
