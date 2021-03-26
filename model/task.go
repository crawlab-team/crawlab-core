package model

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	Id              primitive.ObjectID   `json:"_id" bson:"_id"`
	SpiderId        primitive.ObjectID   `json:"spider_id" bson:"spider_id"`
	StartTs         time.Time            `json:"start_ts" bson:"start_ts"`
	FinishTs        time.Time            `json:"finish_ts" bson:"finish_ts"`
	Status          string               `json:"status" bson:"status"`
	NodeId          primitive.ObjectID   `json:"node_id" bson:"node_id"`
	LogPath         string               `json:"log_path" bson:"log_path"`
	Cmd             string               `json:"cmd" bson:"cmd"`
	Param           string               `json:"param" bson:"param"`
	Error           string               `json:"error" bson:"error"`
	ResultCount     int                  `json:"result_count" bson:"result_count"`
	ErrorLogCount   int                  `json:"error_log_count" bson:"error_log_count"`
	WaitDuration    float64              `json:"wait_duration" bson:"wait_duration"`
	RuntimeDuration float64              `json:"runtime_duration" bson:"runtime_duration"`
	TotalDuration   float64              `json:"total_duration" bson:"total_duration"`
	Pid             int                  `json:"pid" bson:"pid"`
	RunType         string               `json:"run_type" bson:"run_type"`       // deprecated
	ScheduleId      primitive.ObjectID   `json:"schedule_id" bson:"schedule_id"` // Schedule.Id
	Type            string               `json:"type" bson:"type"`
	Mode            string               `json:"mode" bson:"mode"`           // running mode of Task
	NodeIds         []primitive.ObjectID `json:"node_ids" bson:"node_ids"`   // list of Node.Id
	NodeTags        []string             `json:"node_tags" bson:"node_tags"` // list of Node.Tag
	ParentId        primitive.ObjectID   `json:"parent_id" bson:"parent_id"` // parent Task.Id if it's a sub-task
}

func (t *Task) Add() (err error) {
	if t.Id.IsZero() {
		t.Id = primitive.NewObjectID()
	}
	m := NewDelegate(TaskColName, t)
	return m.Add()
}

func (t *Task) Save() (err error) {
	m := NewDelegate(TaskColName, t)
	return m.Save()
}

func (t *Task) Delete() (err error) {
	m := NewDelegate(TaskColName, t)
	return m.Delete()
}

func (t *Task) GetArtifact() (a Artifact, err error) {
	d := NewDelegate(TaskColName, t)
	return d.GetArtifact()
}

type TaskDailyItem struct {
	Date               string  `json:"date" bson:"_id"`
	TaskCount          int     `json:"task_count" bson:"task_count"`
	AvgRuntimeDuration float64 `json:"avg_runtime_duration" bson:"avg_runtime_duration"`
}

const TaskColName = "tasks"

type taskService struct {
	*Service
}

func (svc *taskService) GetById(id primitive.ObjectID) (res Task, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *taskService) Get(query bson.M, opts *mongo.FindOptions) (res Task, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *taskService) GetList(query bson.M, opts *mongo.FindOptions) (res []Task, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *taskService) DeleteById(id primitive.ObjectID) (err error) {
	return svc.deleteId(id)
}

func (svc *taskService) DeleteList(query bson.M) (err error) {
	return svc.delete(query)
}

func (svc *taskService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}

var TaskService = taskService{NewService(TaskColName)}
