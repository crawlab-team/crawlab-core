package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Git struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Url      string             `json:"url" bson:"url"`
	AuthType string             `json:"auth_type" bson:"auth_type"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
}

func (t *Git) GetId() (id primitive.ObjectID) {
	return t.Id
}

func (t *Git) SetId(id primitive.ObjectID) {
	t.Id = id
}
