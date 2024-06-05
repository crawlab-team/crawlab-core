package payload

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelServiceV2Payload[T any] struct {
	Type        string             `json:"type,omitempty"`
	Id          primitive.ObjectID `json:"_id,omitempty"`
	Query       bson.M             `json:"query,omitempty"`
	FindOptions *mongo.FindOptions `json:"find_options,omitempty"`
	Model       T                  `json:"model,omitempty"`
	Update      bson.M             `json:"update,omitempty"`
	Models      []T                `json:"models,omitempty"`
}
