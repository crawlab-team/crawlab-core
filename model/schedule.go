package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
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

type scheduleService struct {
	*Service
}

func (svc *scheduleService) GetById(id primitive.ObjectID) (res Schedule, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *scheduleService) Get(query bson.M, opts *mongo.FindOptions) (res Schedule, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *scheduleService) GetList(query bson.M, opts *mongo.FindOptions) (res []Schedule, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *scheduleService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *scheduleService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *scheduleService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}

var ScheduleService = scheduleService{NewService(ScheduleColName)}
