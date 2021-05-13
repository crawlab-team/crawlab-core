package spider

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceOptions struct {
}

type RunOptions struct {
	Mode       string
	NodeIds    []primitive.ObjectID
	Param      string
	ScheduleId primitive.ObjectID
}

type CloneOptions struct {
	Name string
}

type FsOption func(fsSvc interfaces.SpiderFsService)

func WithId(id primitive.ObjectID) FsOption {
	return func(fsSvc interfaces.SpiderFsService) {
		fsSvc.SetId(id)
	}
}
