package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Permission struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Type        string             `bson:"type" json:"type"`
	Target      string             `bson:"target" json:"target"`
	Filter      string             `bson:"filter" json:"filter"`
}

func (p *Permission) GetId() (id primitive.ObjectID) {
	return p.Id
}

func (p *Permission) SetId(id primitive.ObjectID) {
	p.Id = id
}

type PermissionList []Permission

func (l *PermissionList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
