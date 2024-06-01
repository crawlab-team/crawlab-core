package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectV2 struct {
	Id                     primitive.ObjectID `json:"_id" bson:"_id" collection:"projects"`
	BaseModelV2[ProjectV2] `bson:",inline"`
	Name                   string `json:"name" bson:"name"`
	Description            string `json:"description" bson:"description"`
	Spiders                int    `json:"spiders" bson:"-"`
}
