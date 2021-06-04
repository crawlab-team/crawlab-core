package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Plugin struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Cmd         string             `json:"cmd" bson:"cmd"`
}

func (p *Plugin) GetId() (id primitive.ObjectID) {
	return p.Id
}

func (p *Plugin) SetId(id primitive.ObjectID) {
	p.Id = id
}
