package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type PluginService interface {
	Module
	SetFsPathBase(path string)
	InstallPlugin(id primitive.ObjectID) (err error)
	UninstallPlugin(id primitive.ObjectID) (err error)
	RunPlugin(id primitive.ObjectID) (err error)
	StopPlugin(id primitive.ObjectID) (err error)
}
