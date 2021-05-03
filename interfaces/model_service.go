package interfaces

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/emirpasic/gods/lists/arraylist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelService interface {
	GetById(id primitive.ObjectID) (res interface{}, err error)
	Get(query bson.M, opts *mongo.FindOptions) (res interface{}, err error)
	GetList(query bson.M, opts *mongo.FindOptions) (res arraylist.List, err error)
	DeleteById(id primitive.ObjectID) (err error)
	Delete(query bson.M) (err error)
	DeleteList(query bson.M) (err error)
	UpdateById(id primitive.ObjectID, update interface{}) (err error)
	Update(query bson.M, update interface{}, fields []string) (err error)
	Count(query bson.M) (total int, err error)
}
