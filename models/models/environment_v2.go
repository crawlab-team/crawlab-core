package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EnvironmentV2 struct {
	Id                         primitive.ObjectID `json:"_id" bson:"_id" collection:"environments"`
	BaseModelV2[EnvironmentV2] `bson:",inline"`
	Key                        string `json:"key" bson:"key"`
	Value                      string `json:"value" bson:"value"`
}
