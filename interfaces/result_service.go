package interfaces

import (
	"github.com/crawlab-team/crawlab-db/generic"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResultService interface {
	GetId() (id primitive.ObjectID)
	SetId(id primitive.ObjectID)
	Insert(records ...interface{}) (err error)
	List(query generic.ListQuery, opts *generic.ListOptions) (results []Result, err error)
	Count(query generic.ListQuery) (n int, err error)
}
