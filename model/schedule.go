package model

import (
	"github.com/crawlab-team/crawlab-core/lib/cron"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	Id             primitive.ObjectID   `json:"_id" bson:"_id"`
	Name           string               `json:"name" bson:"name"`
	Description    string               `json:"description" bson:"description"`
	SpiderId       primitive.ObjectID   `json:"spider_id" bson:"spider_id"`
	Cron           string               `json:"cron" bson:"cron"`
	EntryId        cron.EntryID         `json:"entry_id" bson:"entry_id"`
	Param          string               `json:"param" bson:"param"`
	RunType        string               `json:"run_type" bson:"run_type"`
	NodeIds        []primitive.ObjectID `json:"node_ids" bson:"node_ids"`
	Status         string               `json:"status" bson:"status"`
	Enabled        bool                 `json:"enabled" bson:"enabled"`
	UserId         primitive.ObjectID   `json:"user_id" bson:"user_id"`
	ScrapySpider   string               `json:"scrapy_spider" bson:"scrapy_spider"`
	ScrapyLogLevel string               `json:"scrapy_log_level" bson:"scrapy_log_level"`
}

func (s *Schedule) Add() (err error) {
	if s.Id.IsZero() {
		s.Id = primitive.NewObjectID()
	}
	m := NewDelegate(ScheduleColName, s)
	return m.Add()
}

func (s *Schedule) Save() (err error) {
	m := NewDelegate(ScheduleColName, s)
	return m.Save()
}

func (s *Schedule) Delete() (err error) {
	m := NewDelegate(ScheduleColName, s)
	return m.Delete()
}

func (s *Schedule) GetArtifact() (a Artifact, err error) {
	d := NewDelegate(ScheduleColName, s)
	return d.GetArtifact()
}

const ScheduleColName = "schedules"
