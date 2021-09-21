package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PluginService interface {
	Module
	SetFsPathBase(path string)
	SetMonitorInterval(interval time.Duration)
	SetPluginBaseUrl(baseUrl string)
	InstallPlugin(id primitive.ObjectID) (err error)
	UninstallPlugin(id primitive.ObjectID) (err error)
	RunPlugin(id primitive.ObjectID) (err error)
	StopPlugin(id primitive.ObjectID) (err error)
}
