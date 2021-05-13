package interfaces

import (
	"github.com/crawlab-team/crawlab-core/spider"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderService interface {
	Run(id primitive.ObjectID, opts *spider.RunOptions) (err error)
	Clone(id primitive.ObjectID, opts *spider.CloneOptions) (err error)
	Delete(id primitive.ObjectID) (err error)
	Sync(id primitive.ObjectID) (err error)
	GetFs(id primitive.ObjectID) (fsSvc *spider.FsService, err error)
}
