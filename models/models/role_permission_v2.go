package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RolePermissionV2 struct {
	Id                            primitive.ObjectID `json:"_id" bson:"_id" collection:"role_permissions"`
	BaseModelV2[RolePermissionV2] `bson:",inline"`
	RoleId                        primitive.ObjectID `json:"role_id" bson:"role_id"`
	PermissionId                  primitive.ObjectID `json:"permission_id" bson:"permission_id"`
}
