package service_test

import (
	"context"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func SetupTest(t *testing.T) {
	if err := mongo.InitMongo(); err != nil {
		panic(err)
	}
	t.Cleanup(CleanupTest)
}

func CleanupTest() {
	db := mongo.GetMongoDb("")
	names, _ := db.ListCollectionNames(context.Background(), nil)
	for _, n := range names {
		_, _ = db.Collection(n).DeleteMany(context.Background(), bson.M{})
	}

	// avoid caching
	time.Sleep(200 * time.Millisecond)
}
