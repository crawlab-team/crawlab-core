package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PermissionV2 struct {
	Id                        primitive.ObjectID `json:"_id" bson:"_id" collection:"permissions"`
	BaseModelV2[PermissionV2] `bson:",inline"`
	Key                       string   `json:"key" bson:"key"`
	Name                      string   `json:"name" bson:"name"`
	Description               string   `json:"description" bson:"description"`
	Type                      string   `json:"type" bson:"type"`
	Target                    []string `json:"target" bson:"target"`
	Allow                     []string `json:"allow" bson:"allow"`
	Deny                      []string `json:"deny" bson:"deny"`
}
