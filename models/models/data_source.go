package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DataSource struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Type        string             `json:"type" bson:"type"`
	Description string             `json:"description" bson:"description"`
	Url         string             `json:"url" bson:"url"`
	Host        string             `json:"host" bson:"host"`
	Port        string             `json:"port" bson:"port"`
	Database    string             `json:"database" bson:"database"`
	Username    string             `json:"username" bson:"username"`
	Password    string             `json:"password,omitempty" bson:"-"`
	Extra       map[string]string  `json:"extra,omitempty" bson:"extra,omitempty"`
}

func (ds *DataSource) GetId() (id primitive.ObjectID) {
	return ds.Id
}

func (ds *DataSource) SetId(id primitive.ObjectID) {
	ds.Id = id
}
