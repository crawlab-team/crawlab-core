package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Plugin struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Type        string             `json:"type" bson:"type"`
	Proto       string             `json:"proto" bson:"proto"`
	Active      bool               `json:"active" bson:"active"`
	Endpoint    string             `json:"endpoint" bson:"endpoint"`
	Cmd         string             `json:"cmd" bson:"cmd"`
}

func (p *Plugin) GetId() (id primitive.ObjectID) {
	return p.Id
}

func (p *Plugin) SetId(id primitive.ObjectID) {
	p.Id = id
}
