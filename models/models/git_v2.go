package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GitV2 struct {
	Id                 primitive.ObjectID `json:"_id" bson:"_id" collection:"gits"`
	BaseModelV2[GitV2] `bson:",inline"`
	Url                string `json:"url" bson:"url"`
	AuthType           string `json:"auth_type" bson:"auth_type"`
	Username           string `json:"username" bson:"username"`
	Password           string `json:"password" bson:"password"`
	CurrentBranch      string `json:"current_branch" bson:"current_branch"`
	AutoPull           bool   `json:"auto_pull" bson:"auto_pull"`
}
