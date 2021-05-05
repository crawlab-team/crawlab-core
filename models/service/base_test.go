package service_test

import (
	"context"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func SetupTest(t *testing.T) {
	if err := mongo.InitMongo(); err != nil {
		panic(err)
	}
	t.Cleanup(CleanupTest)
}

func CleanupTest() {
	db := mongo.GetMongoDb("")
	names, err := db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	for _, n := range names {
		if err := db.Collection(n).Drop(context.Background()); err != nil {
			panic(err)
		}
	}
}
