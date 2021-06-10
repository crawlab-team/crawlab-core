package interfaces

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/emirpasic/gods/lists/arraylist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelBaseService interface {
	GetModelId() (id ModelId)
	SetModelId(id ModelId)
	GetCol() (col *mongo.Col)
	SetCol(col *mongo.Col)
	GetById(id primitive.ObjectID) (res Model, err error)
	Get(query bson.M, opts *mongo.FindOptions) (res Model, err error)
	GetList(query bson.M, opts *mongo.FindOptions) (res arraylist.List, err error)
	DeleteById(id primitive.ObjectID) (err error)
	Delete(query bson.M) (err error)
	DeleteList(query bson.M) (err error)
	ForceDeleteList(query bson.M) (err error)
	UpdateById(id primitive.ObjectID, update bson.M) (err error)
	Update(query bson.M, update bson.M, fields []string) (err error)
	UpdateDoc(query bson.M, doc Model, fields []string) (err error)
	Insert(docs ...interface{}) (err error)
	Count(query bson.M) (total int, err error)
}

type ModelService interface {
	GetBaseService(id ModelId) (svc ModelBaseService)
}
