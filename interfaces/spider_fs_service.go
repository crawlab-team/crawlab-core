package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type SpiderFsService interface {
	GetFsPath() (res string)
	GetWorkspacePath() (res string)
	GetRepoPath() (res string)
	SetId(id primitive.ObjectID)
	SetFsPathBase(path string)
	SetWorkspacePathBase(path string)
	SetRepoPathBase(path string)
	GetFsService() (fsSvc FsService)
}
