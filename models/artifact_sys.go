package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ArtifactSys struct {
	CreateTs  time.Time          `json:"create_ts" bson:"create_ts"`
	CreateUid primitive.ObjectID `json:"create_uid" bson:"create_uid"`
	UpdateTs  time.Time          `json:"update_ts" bson:"update_ts"`
	UpdateUid primitive.ObjectID `json:"update_uid" bson:"update_uid"`
	DeleteTs  time.Time          `json:"delete_ts" bson:"delete_ts"`
	DeleteUid primitive.ObjectID `json:"delete_uid" bson:"delete_uid"`
}
