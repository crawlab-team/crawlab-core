package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskStat struct {
	Id              primitive.ObjectID `json:"_id" bson:"_id"`
	CreateTs        time.Time          `json:"create_ts" bson:"cts"`
	StartTs         time.Time          `json:"start_ts" bson:"sts"`
	EndTs           time.Time          `json:"end_ts" bson:"ets"`
	WaitDuration    float64            `json:"wait_duration" bson:"wd"`
	RuntimeDuration float64            `json:"runtime_duration" bson:"rd"`
	TotalDuration   float64            `json:"total_duration" bson:"td"`
	ResultCount     int                `json:"result_count" bson:"rc"`
	ErrorLogCount   int                `json:"error_log_count" bson:"elc"`
}

func (ts *TaskStat) GetId() (id primitive.ObjectID) {
	return ts.Id
}

func (ts *TaskStat) SetId(id primitive.ObjectID) {
	ts.Id = id
}
