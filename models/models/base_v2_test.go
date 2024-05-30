package models_test

import (
	"context"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"testing"

	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestModel struct {
	models.BaseModelV2[TestModel] `bson:",inline" collection:"test_collection"`
	Name                          string `json:"name" bson:"name"`
}

func setup() {
	viper.Set("mongo.db", "testdb")
}

func teardown() {
	// Clean up the database after tests
	mongo.GetMongoCol("test_collection").GetCollection().Drop(context.Background())
}

func TestBaseModelV2_Save(t *testing.T) {
	setup()
	defer teardown()

	ctx := context.Background()

	t.Run("Insert New Document", func(t *testing.T) {
		model := &TestModel{Name: "Test"}
		err := model.Save(ctx)
		assert.NoError(t, err)
		assert.False(t, model.GetId().IsZero())

		collection := mongo.GetMongoCol("test_collection")
		var result TestModel
		err = collection.GetCollection().FindOne(ctx, bson.M{"_id": model.GetId()}).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, "Test", result.Name)
	})

	t.Run("Update Existing Document", func(t *testing.T) {
		model := &TestModel{Name: "Test"}
		model.Save(ctx)

		model.Name = "Updated Test"
		err := model.Save(ctx)
		assert.NoError(t, err)

		collection := mongo.GetMongoCol("test_collection")
		var result TestModel
		err = collection.GetCollection().FindOne(ctx, bson.M{"_id": model.GetId()}).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Test", result.Name)
	})
}

func TestBaseModelV2_Delete(t *testing.T) {
	setup()
	defer teardown()

	ctx := context.Background()

	t.Run("Delete Existing Document", func(t *testing.T) {
		model := &TestModel{Name: "Test"}
		model.Save(ctx)

		err := model.Delete(ctx)
		assert.NoError(t, err)

		collection := mongo.GetMongoCol("test_collection")
		count, err := collection.GetCollection().CountDocuments(ctx, bson.M{"_id": model.GetId()})
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Delete Non-Existent Document", func(t *testing.T) {
		model := &TestModel{}
		model.SetId(primitive.NewObjectID())
		err := model.Delete(ctx)
		assert.NoError(t, err)
	})
}

func TestGetCollectionName(t *testing.T) {
	t.Run("Get Collection Name from Struct Tag", func(t *testing.T) {
		model := &TestModel{}
		name, err := models.GetCollectionName(model)
		assert.NoError(t, err)
		assert.Equal(t, "test_collection", name)
	})

	t.Run("Missing Collection Name Tag", func(t *testing.T) {
		model := &struct {
			models.BaseModelV2[TestModel] `bson:",inline"`
		}{}
		_, err := models.GetCollectionName(model)
		assert.Error(t, err)
	})
}
