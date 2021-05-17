package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskSchedulerService interface {
	TaskBaseService
	// Assign enqueue task
	Assign(t Task) (err error)
	// Fetch continuously dequeue task from task queue
	Fetch() (err error)
	// Run task via sending grpc stream message to handler(s)
	Run(taskId primitive.ObjectID) (err error)
	// Cancel task via sending grpc stream message to handler(s)
	Cancel(taskId primitive.ObjectID) (err error)
}
