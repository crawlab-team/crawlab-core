package entity

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrpcBaseServiceParams struct {
	Query       bson.M             `json:"q"`
	Id          primitive.ObjectID `json:"id"`
	Update      bson.M             `json:"u"`
	Fields      []string           `json:"f"`
	FindOptions *mongo.FindOptions `json:"o"`
	Docs        []interfaces.Model `json:"d"`
}

func (params *GrpcBaseServiceParams) Value() interface{} {
	return params
}
