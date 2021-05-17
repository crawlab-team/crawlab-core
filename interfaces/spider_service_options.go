package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type ServiceOption func(svc SpiderService)

type RunOptions struct {
	Mode       string
	NodeIds    []primitive.ObjectID
	Param      string
	ScheduleId primitive.ObjectID
}

type CloneOptions struct {
	Name string
}

type SpiderFsOption func(fsSvc SpiderFsService)

func WithFsPathBase(path string) SpiderFsOption {
	return func(fsSvc SpiderFsService) {
		fsSvc.SetFsPathBase(path)
	}
}

func WithWorkspacePathBase(path string) SpiderFsOption {
	return func(fsSvc SpiderFsService) {
		fsSvc.SetWorkspacePathBase(path)
	}
}

func WithRepoPathBase(path string) SpiderFsOption {
	return func(fsSvc SpiderFsService) {
		fsSvc.SetRepoPathBase(path)
	}
}
