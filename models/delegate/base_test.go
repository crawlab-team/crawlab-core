package delegate_test

import (
	"context"
	"github.com/crawlab-team/crawlab-db/mongo"
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
	names, _ := db.ListCollectionNames(context.Background(), nil)
	for _, n := range names {
		_ = db.Collection(n).Drop(context.Background())
	}
}
