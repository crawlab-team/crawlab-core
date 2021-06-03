package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id         primitive.ObjectID   `json:"_id" bson:"_id"`
	SpiderId   primitive.ObjectID   `json:"spider_id" bson:"spider_id"`
	Status     string               `json:"status" bson:"status"`
	NodeId     primitive.ObjectID   `json:"node_id" bson:"node_id"`
	Cmd        string               `json:"cmd" bson:"cmd"`
	Param      string               `json:"param" bson:"param"`
	Error      string               `json:"error" bson:"error"`
	Pid        int                  `json:"pid" bson:"pid"`
	ScheduleId primitive.ObjectID   `json:"schedule_id" bson:"schedule_id"` // Schedule.Id
	Type       string               `json:"type" bson:"type"`
	Mode       string               `json:"mode" bson:"mode"`           // running mode of Task
	NodeIds    []primitive.ObjectID `json:"node_ids" bson:"node_ids"`   // list of Node.Id
	NodeTags   []string             `json:"node_tags" bson:"node_tags"` // list of Node.Tag
	ParentId   primitive.ObjectID   `json:"parent_id" bson:"parent_id"` // parent Task.Id if it'Spider a sub-task
	Priority   int                  `json:"priority" bson:"priority"`
	Stat       *TaskStat            `json:"stat,omitempty" bson:"-"`
}

func (t *Task) GetId() (id primitive.ObjectID) {
	return t.Id
}

func (t *Task) SetId(id primitive.ObjectID) {
	t.Id = id
}

func (t *Task) GetNodeId() (id primitive.ObjectID) {
	return t.NodeId
}

func (t *Task) SetNodeId(id primitive.ObjectID) {
	t.NodeId = id
}

func (t *Task) GetNodeIds() (ids []primitive.ObjectID) {
	return t.NodeIds
}

func (t *Task) GetNodeTags() (nodeTags []string) {
	return t.NodeTags
}

func (t *Task) GetStatus() (status string) {
	return t.Status
}

func (t *Task) SetStatus(status string) {
	t.Status = status
}

func (t *Task) GetError() (error string) {
	return t.Error
}

func (t *Task) SetError(error string) {
	t.Error = error
}

func (t *Task) GetSpiderId() (id primitive.ObjectID) {
	return t.SpiderId
}

func (t *Task) GetType() (ty string) {
	return t.Type
}

func (t *Task) GetCmd() (cmd string) {
	return t.Cmd
}

func (t *Task) GetParam() (param string) {
	return t.Param
}

func (t *Task) GetPriority() (p int) {
	return t.Priority
}

type TaskDailyItem struct {
	Date               string  `json:"date" bson:"_id"`
	TaskCount          int     `json:"task_count" bson:"task_count"`
	AvgRuntimeDuration float64 `json:"avg_runtime_duration" bson:"avg_runtime_duration"`
}
