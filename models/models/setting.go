package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Setting struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Key   string             `json:"key" bson:"key"`
	Value string             `json:"value" bson:"value"`
}

func (s *Setting) GetId() (id primitive.ObjectID) {
	return s.Id
}

func (s *Setting) SetId(id primitive.ObjectID) {
	s.Id = id
}
