package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EnvironmentV2 struct {
	any                        `collection:"environments"`
	BaseModelV2[EnvironmentV2] `bson:",inline"`
	Key                        string `json:"key" bson:"key"`
	Value                      string `json:"value" bson:"value"`
}
