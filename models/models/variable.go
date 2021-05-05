package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Variable struct {
	Id     primitive.ObjectID `json:"_id" bson:"_id"`
	Key    string             `json:"key" bson:"key"`
	Value  string             `json:"value" bson:"value"`
	Remark string             `json:"remark" bson:"remark"`
}

func (v *Variable) GetId() (id primitive.ObjectID) {
	return v.Id
}

func (v *Variable) SetId(id primitive.ObjectID) {
	v.Id = id
}
