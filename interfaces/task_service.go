package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskService interface {
	Init() (err error)
	Close()
	Assign(t Task) (err error)
	Fetch() (t Task, err error)
	Run(taskId primitive.ObjectID) (err error)
	Cancel(taskId primitive.ObjectID) (err error)
	FindLogs(id primitive.ObjectID, pattern string, skip, size int) (lines []string, err error)
	SetMaxRunners(maxRunners int)
	SetPollWaitDuration(duration time.Duration)
}
