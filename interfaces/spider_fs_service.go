package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type SpiderFsService interface {
	SetId(id primitive.ObjectID)
	GetFsPath() (res string)
	GetWorkspacePath() (res string)
	GetRepoPath() (res string)
	SetFsPathBase(path string)
	SetWorkspacePathBase(path string)
	SetRepoPathBase(path string)
	GetFsService() (fsSvc FsService)
}
