package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Password struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Password string             `json:"password" bson:"p"`
}

func (p *Password) GetId() (id primitive.ObjectID) {
	return p.Id
}

func (p *Password) SetId(id primitive.ObjectID) {
	p.Id = id
}
