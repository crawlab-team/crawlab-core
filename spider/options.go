package spider

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WithId(id primitive.ObjectID) interfaces.FsOption {
	return func(fsSvc interfaces.SpiderFsService) {
		fsSvc.SetId(id)
	}
}
