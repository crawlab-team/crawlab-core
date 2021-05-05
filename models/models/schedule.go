package models

import (
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	SpiderId    primitive.ObjectID `json:"spider_id" bson:"spider_id"`
	Cron        string             `json:"cron" bson:"cron"`
	EntryId     cron.EntryID       `json:"entry_id" bson:"entry_id"`
	Param       string             `json:"param" bson:"param"`
	//RunType        string               `json:"run_type" bson:"run_type"`
	NodeIds        []primitive.ObjectID `json:"node_ids" bson:"node_ids"`
	Enabled        bool                 `json:"enabled" bson:"enabled"`
	Mode           string               `json:"mode" bson:"mode"`
	UserId         primitive.ObjectID   `json:"user_id" bson:"user_id"`
	ScrapySpider   string               `json:"scrapy_spider" bson:"scrapy_spider"`
	ScrapyLogLevel string               `json:"scrapy_log_level" bson:"scrapy_log_level"`
	Tags           []string             `json:"tags" bson:"-"`
}

func (s *Schedule) GetId() (id primitive.ObjectID) {
	return s.Id
}

func (s *Schedule) SetId(id primitive.ObjectID) {
	s.Id = id
}
