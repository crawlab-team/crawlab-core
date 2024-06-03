package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestModel struct {
	Id                     primitive.ObjectID `json:"_id" bson:"_id" collection:"testmodels"`
	BaseModelV2[TestModel] `bson:",inline"`
	Name                   string `json:"name" bson:"name"`
}
