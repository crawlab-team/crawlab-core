package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type SpiderFsService interface {
	GetFsPath(id primitive.ObjectID) (res string)
	GetWorkspacePath(id primitive.ObjectID) (res string)
	GetRepoPath(id primitive.ObjectID) (res string)
	SetId(id primitive.ObjectID)
}
