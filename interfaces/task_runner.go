package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskRunner interface {
	Init() (err error)
	Run() (err error)
	Cancel() (err error)
	Dispose() (err error)
	SetLogDriverType(driverType string)
	GetTaskId() (id primitive.ObjectID)
}
