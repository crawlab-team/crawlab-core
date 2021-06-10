package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DataCollection struct {
	Id   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

func (dc *DataCollection) GetId() (id primitive.ObjectID) {
	return dc.Id
}

func (dc *DataCollection) SetId(id primitive.ObjectID) {
	dc.Id = id
}
