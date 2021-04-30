package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
	Id   primitive.ObjectID `json:"_id" bson:"_id"`
	Tags []Tag              `json:"tags" bson:"-"`
}
