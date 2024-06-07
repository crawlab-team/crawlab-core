package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VariableV2 struct {
	Id                      primitive.ObjectID `json:"_id" bson:"_id" collection:"variables"`
	BaseModelV2[VariableV2] `bson:",inline"`
	Key                     string `json:"key" bson:"key"`
	Value                   string `json:"value" bson:"value"`
	Remark                  string `json:"remark" bson:"remark"`
}
