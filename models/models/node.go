package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Node struct {
	Id               primitive.ObjectID `json:"_id" bson:"_id"`
	Key              string             `json:"key" bson:"k"`
	Name             string             `json:"name" bson:"n"`
	Ip               string             `json:"ip" bson:"ip"`
	Port             string             `json:"port" bson:"p"`
	Mac              string             `json:"mac" bson:"m"`
	Hostname         string             `json:"hostname" bson:"h"`
	Description      string             `json:"description" bson:"d"`
	IsMaster         bool               `json:"is_master" bson:"im"`
	Status           string             `json:"status" bson:"s"`
	Enabled          bool               `json:"enabled" bson:"en"`
	Active           bool               `json:"active" bson:"a"`
	ActiveTs         time.Time          `json:"active_ts" bson:"at"`
	AvailableRunners int                `json:"available_runners" bson:"ar"`
	MaxRunners       int                `json:"max_runners" bson:"mr"`
	Tags             []Tag              `json:"tags" bson:"-"`
}

func (n *Node) GetId() (id primitive.ObjectID) {
	return n.Id
}

func (n *Node) SetId(id primitive.ObjectID) {
	n.Id = id
}

func (n *Node) GetTags() (tags []interfaces.Tag) {
	return convertTagsToInterfaces(n.Tags)
}

func (n *Node) SetTags(tags []interfaces.Tag) {
	n.Tags = convertInterfacesToTags(tags)
}

func (n *Node) GetKey() (key string) {
	return n.Key
}

func (n *Node) SetActive(active bool) {
	n.Active = active
}

func (n *Node) SetActiveTs(activeTs time.Time) {
	n.ActiveTs = activeTs
}

func (n *Node) GetStatus() (status string) {
	return n.Status
}

func (n *Node) SetStatus(status string) {
	n.Status = status
}

func (n *Node) GetEnabled() (enabled bool) {
	return n.Enabled
}

func (n *Node) SetEnabled(enabled bool) {
	n.Enabled = enabled
}

func (n *Node) GetAvailableRunners() (runners int) {
	return n.AvailableRunners
}

func (n *Node) SetAvailableRunners(runners int) {
	n.AvailableRunners = runners
}

func (n *Node) GetMaxRunners() (runners int) {
	return n.MaxRunners
}

func (n *Node) SetMaxRunners(runners int) {
	n.MaxRunners = runners
}

func (n *Node) IncrementAvailableRunners() {
	n.AvailableRunners++
}

func (n *Node) DecrementAvailableRunners() {
	n.AvailableRunners--
}
