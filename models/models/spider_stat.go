package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SpiderStat struct {
	Id                     primitive.ObjectID `json:"_id" bson:"_id"`
	LastTaskId             primitive.ObjectID `json:"last_task_id" bson:"ltid"`
	LastTask               *Task              `json:"last_task,omitempty" bson:"-"`
	Tasks                  int                `json:"tasks" bson:"t"`
	Results                int                `json:"results" bson:"r"`
	WaitDuration           int64              `json:"wait_duration" bson:"wd"`           // in second
	RuntimeDuration        int64              `json:"runtime_duration" bson:"rd"`        // in second
	TotalDuration          int64              `json:"total_duration" bson:"td"`          // in second
	AverageWaitDuration    int64              `json:"average_wait_duration" bson:"-"`    // in second
	AverageRuntimeDuration int64              `json:"average_runtime_duration" bson:"-"` // in second
	AverageTotalDuration   int64              `json:"average_total_duration" bson:"-"`   // in second
}

func (s *SpiderStat) GetId() (id primitive.ObjectID) {
	return s.Id
}

func (s *SpiderStat) SetId(id primitive.ObjectID) {
	s.Id = id
}
