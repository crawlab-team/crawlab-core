package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	Role     string             `json:"role" bson:"role"`
	Email    string             `json:"email" bson:"email"`
	Setting  UserSetting        `json:"setting" bson:"setting"`
}

type UserSetting struct {
	NotificationTrigger  string   `json:"notification_trigger" bson:"notification_trigger"`
	DingTalkRobotWebhook string   `json:"ding_talk_robot_webhook" bson:"ding_talk_robot_webhook"`
	WechatRobotWebhook   string   `json:"wechat_robot_webhook" bson:"wechat_robot_webhook"`
	EnabledNotifications []string `json:"enabled_notifications" bson:"enabled_notifications"`
	ErrorRegexPattern    string   `json:"error_regex_pattern" bson:"error_regex_pattern"`
	MaxErrorLog          int      `json:"max_error_log" bson:"max_error_log"`
	LogExpireDuration    int64    `json:"log_expire_duration" bson:"log_expire_duration"`
}

func (u *User) Add() (err error) {
	if u.Id.IsZero() {
		u.Id = primitive.NewObjectID()
	}
	m := NewDelegate(ModelIdUser, u)
	return m.Add()
}

func (u *User) Save() (err error) {
	m := NewDelegate(ModelIdUser, u)
	return m.Save()
}

func (u *User) Delete() (err error) {
	m := NewDelegate(ModelIdUser, u)
	return m.Delete()
}

func (u *User) GetArtifact() (a Artifact, err error) {
	d := NewDelegate(ModelIdUser, u)
	return d.GetArtifact()
}

func (u *User) GetId() (id primitive.ObjectID) {
	return u.Id
}

const UserContextKey = "user"
