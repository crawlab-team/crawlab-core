package models

type UserV2 struct {
	BaseModelV2[UserV2] `bson:",inline"`
}
