package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleV2 struct {
	Id                  primitive.ObjectID `json:"_id" bson:"_id" collection:"roles"`
	BaseModelV2[RoleV2] `bson:",inline"`
	Key                 string `json:"key" bson:"key"`
	Name                string `json:"name" bson:"name"`
	Description         string `json:"description" bson:"description"`
}
