package models

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataCollectionV2 struct {
	Id                          primitive.ObjectID `json:"_id" bson:"_id" collection:"data_collections"`
	BaseModelV2[DataCollection] `bson:",inline"`
	Name                        string             `json:"name" bson:"name"`
	Fields                      []entity.DataField `json:"fields" bson:"fields"`
	Dedup                       struct {
		Enabled bool     `json:"enabled" bson:"enabled"`
		Keys    []string `json:"keys" bson:"keys"`
		Type    string   `json:"type" bson:"type"`
	} `json:"dedup" bson:"dedup"`
}
