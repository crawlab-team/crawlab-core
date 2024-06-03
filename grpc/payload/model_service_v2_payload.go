package payload

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelServiceV2Payload[T any] struct {
	Type        string             `json:"type"`
	Id          primitive.ObjectID `json:"_id"`
	Query       bson.M             `json:"query"`
	FindOptions mongo.FindOptions  `json:"find_options"`
	Model       T                  `json:"model"`
	Update      bson.M             `json:"update"`
	Models      []T                `json:"models"`
}
