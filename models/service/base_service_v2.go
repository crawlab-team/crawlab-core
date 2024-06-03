package service

import (
	"fmt"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"sync"
)

var (
	instanceMap = make(map[string]interface{})
	onceMap     = make(map[string]*sync.Once)
	mu          sync.Mutex
)

type ModelServiceV2[T any] struct {
	col *mongo.Col
}

func (svc *ModelServiceV2[T]) GetById(id primitive.ObjectID) (model *T, err error) {
	var result T
	err = svc.col.FindId(id).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (svc *ModelServiceV2[T]) Get(query bson.M, options *mongo.FindOptions) (model *T, err error) {
	var result T
	err = svc.col.Find(query, options).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (svc *ModelServiceV2[T]) GetList(query bson.M, options *mongo.FindOptions) (models []T, err error) {
	var result []T
	err = svc.col.Find(query, options).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (svc *ModelServiceV2[T]) DeleteById(id primitive.ObjectID) (err error) {
	return svc.col.DeleteId(id)
}

func (svc *ModelServiceV2[T]) Delete(query bson.M) (err error) {
	_, err = svc.col.GetCollection().DeleteOne(svc.col.GetContext(), query)
	return err
}

func (svc *ModelServiceV2[T]) DeleteList(query bson.M) (err error) {
	return svc.col.Delete(query)
}

func (svc *ModelServiceV2[T]) UpdateById(id primitive.ObjectID, update bson.M) (err error) {
	return svc.col.UpdateId(id, update)
}

func (svc *ModelServiceV2[T]) UpdateOne(query bson.M, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateOne(svc.col.GetContext(), query, update)
	return err
}

func (svc *ModelServiceV2[T]) UpdateMany(query bson.M, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateMany(svc.col.GetContext(), query, update)
	return err
}

func (svc *ModelServiceV2[T]) ReplaceById(id primitive.ObjectID, model T) (err error) {
	_, err = svc.col.GetCollection().ReplaceOne(svc.col.GetContext(), bson.M{"_id": id}, model)
	return err
}

func (svc *ModelServiceV2[T]) Replace(query bson.M, model T) (err error) {
	_, err = svc.col.GetCollection().ReplaceOne(svc.col.GetContext(), query, model)
	return err
}

func (svc *ModelServiceV2[T]) InsertOne(model T) (id primitive.ObjectID, err error) {
	res, err := svc.col.GetCollection().InsertOne(svc.col.GetContext(), model)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (svc *ModelServiceV2[T]) InsertMany(models []T) (ids []primitive.ObjectID, err error) {
	var _models []interface{}
	for _, model := range models {
		_models = append(_models, model)
	}
	res, err := svc.col.GetCollection().InsertMany(svc.col.GetContext(), _models)
	if err != nil {
		return nil, err
	}
	for _, v := range res.InsertedIDs {
		ids = append(ids, v.(primitive.ObjectID))
	}
	return ids, nil
}

func (svc *ModelServiceV2[T]) Count(query bson.M) (total int, err error) {
	return svc.col.Count(query)
}

func (svc *ModelServiceV2[T]) GetCol() (col *mongo.Col) {
	return svc.col
}

func getCollectionName[T any]() string {
	var instance T
	t := reflect.TypeOf(instance)
	field, _ := t.FieldByName("Id")
	return field.Tag.Get("collection")
}

// NewModelServiceV2 return singleton instance of ModelServiceV2
func NewModelServiceV2[T any]() *ModelServiceV2[T] {
	typeName := fmt.Sprintf("%T", *new(T))

	mu.Lock()
	defer mu.Unlock()

	if _, exists := onceMap[typeName]; !exists {
		onceMap[typeName] = &sync.Once{}
	}

	var instance *ModelServiceV2[T]

	onceMap[typeName].Do(func() {
		collectionName := getCollectionName[T]()
		collection := mongo.GetMongoCol(collectionName)
		instance = &ModelServiceV2[T]{col: collection}
		instanceMap[typeName] = instance
	})

	return instanceMap[typeName].(*ModelServiceV2[T])
}
