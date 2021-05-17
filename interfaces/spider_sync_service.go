package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderSyncService interface {
	WithConfigPath
	SyncToFs(id primitive.ObjectID) (err error)
	SyncToWorkspace(id primitive.ObjectID) (err error)
}
