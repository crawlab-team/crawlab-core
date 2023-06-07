package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Environment struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Key   string             `json:"key" bson:"key"`
	Value string             `json:"value" bson:"value"`
}

func (e *Environment) GetId() (id primitive.ObjectID) {
	return e.Id
}

func (e *Environment) SetId(id primitive.ObjectID) {
	e.Id = id
}

type EnvironmentList []Environment

func (l *EnvironmentList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
