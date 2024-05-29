package models

import (
	"context"
	"errors"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

type ModelV2[T any] interface {
	GetId() (id primitive.ObjectID)
	SetId(id primitive.ObjectID)
	Save(ctx context.Context) (err error)
	Delete(ctx context.Context) (err error)
}

type BaseModelV2[T any] struct {
	Id primitive.ObjectID `json:"_id" bson:"_id"`
}

func (m *BaseModelV2[T]) GetId() (id primitive.ObjectID) {
	return m.Id
}

func (m *BaseModelV2[T]) SetId(id primitive.ObjectID) {
	m.Id = id
}

func (m *BaseModelV2[T]) Save(ctx context.Context) (err error) {
	collectionName, err := GetCollectionName(new(T))
	if err != nil {
		return err
	}
	collection := mongo.GetMongoCol(collectionName)

	if m.Id.IsZero() {
		m.Id = primitive.NewObjectID()
		res, err := collection.GetCollection().InsertOne(ctx, m)
		if err != nil {
			return err
		}
		m.Id = res.InsertedID.(primitive.ObjectID)
	} else {
		_, err = collection.GetCollection().ReplaceOne(ctx, bson.M{"_id": m.Id}, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *BaseModelV2[T]) Delete(ctx context.Context) (err error) {
	collectionName, err := GetCollectionName(new(T))
	if err != nil {
		return err
	}
	collection := mongo.GetMongoCol(collectionName)
	_, err = collection.GetCollection().DeleteOne(ctx, bson.M{"_id": m.Id})
	return err
}

func GetCollectionName(model interface{}) (string, error) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	collectionTag := t.Field(0).Tag.Get("collection")
	if collectionTag == "" {
		return "", errors.New("collection name not specified in struct tag")
	}
	return collectionTag, nil
}
