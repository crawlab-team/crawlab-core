package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
}

func (r *Role) GetId() (id primitive.ObjectID) {
	return r.Id
}

func (r *Role) SetId(id primitive.ObjectID) {
	r.Id = id
}

type RoleList []Role

func (l *RoleList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
