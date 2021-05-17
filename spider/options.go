package spider

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WithId(id primitive.ObjectID) interfaces.SpiderFsOption {
	return func(fsSvc interfaces.SpiderFsService) {
		fsSvc.SetId(id)
	}
}

func WithFsPathBase(path string) interfaces.SpiderFsOption {
	return func(svc interfaces.SpiderFsService) {
		svc.SetFsPathBase(path)
	}
}

func WithWorkspacePathBase(path string) interfaces.SpiderFsOption {
	return func(svc interfaces.SpiderFsService) {
		svc.SetWorkspacePathBase(path)
	}
}

func WithRepoPathBase(path string) interfaces.SpiderFsOption {
	return func(svc interfaces.SpiderFsService) {
		svc.SetRepoPathBase(path)
	}
}
