package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskManagerService interface {
	TaskBaseService
	// Enqueue task into the task queue
	Enqueue(t Task) (err error)
	// Cancel task via sending grpc stream message to handler(s)
	Cancel(taskId primitive.ObjectID) (err error)
}
