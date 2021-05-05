package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Token string             `json:"token" bson:"token"`
}

func (t *Token) GetId() (id primitive.ObjectID) {
	return t.Id
}

func (t *Token) SetId(id primitive.ObjectID) {
	t.Id = id
}
