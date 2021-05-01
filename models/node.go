package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Node struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Key         string             `json:"key" bson:"key"`
	Name        string             `json:"name" bson:"name"`
	Ip          string             `json:"ip" bson:"ip"`
	Port        string             `json:"port" bson:"port"`
	Mac         string             `json:"mac" bson:"mac"`
	Hostname    string             `json:"hostname" bson:"hostname"`
	Description string             `json:"description" bson:"description"`
	IsMaster    bool               `json:"is_master" bson:"is_master"`
	Status      string             `json:"status" bson:"status"`
	Enabled     bool               `json:"enabled" bson:"enabled"`
	Active      bool               `json:"active" bson:"active"`
	ActiveTs    time.Time          `json:"active_ts" bson:"active_ts"`
	Settings    NodeSettings       `json:"settings" bson:"settings"`
	Tags        []interfaces.Tag   `json:"tags" bson:"-"`
}

type NodeSettings struct {
	MaxRunners int `json:"max_runners" bson:"max_runners"`
}

func (n *Node) Add() (err error) {
	if n.Id.IsZero() {
		n.Id = primitive.NewObjectID()
	}
	m := NewDelegate(interfaces.ModelIdNode, n)
	return m.Add()
}

func (n *Node) Save() (err error) {
	m := NewDelegate(interfaces.ModelIdNode, n)
	return m.Save()
}

func (n *Node) Delete() (err error) {
	m := NewDelegate(interfaces.ModelIdNode, n)
	return m.Delete()
}

func (n *Node) GetArtifact() (a interfaces.ModelArtifact, err error) {
	m := NewDelegate(interfaces.ModelIdNode, n)
	return m.GetArtifact()
}

func (n *Node) GetId() (id primitive.ObjectID) {
	return n.Id
}

func (n *Node) SetTags(tags []interfaces.Tag) {
	n.Tags = tags
}
