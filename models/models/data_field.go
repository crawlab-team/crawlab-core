package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DataField struct {
	Id   primitive.ObjectID `json:"_id" bson:"_id"`
	Key  string             `json:"key" bson:"key"`
	Name string             `json:"name" bson:"name"`
}
