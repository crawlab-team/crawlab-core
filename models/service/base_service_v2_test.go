package service_test

import (
	"context"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"

	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

type TestModel struct {
	Id                            primitive.ObjectID `bson:"_id,omitempty" collection:"testmodels"`
	models.BaseModelV2[TestModel] `bson:",inline"`
	Name                          string `bson:"name"`
}

func setupTestDB() {
	viper.Set("mongo.db", "testdb")
}

func teardownTestDB() {
	db := mongo.GetMongoDb("testdb")
	db.Drop(context.Background())
}

func TestModelServiceV2_GetById(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Arrange
	svc := service.NewModelServiceV2[TestModel]()
	id := primitive.NewObjectID()
	testModel := TestModel{Id: id, Name: "Test Name"}
	testModel.SetCreatedAt(time.Now())
	testModel.SetUpdatedAt(time.Now())
	_, err := svc.InsertOne(testModel)
	assert.Nil(t, err)

	// Act
	result, err := svc.GetById(id)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Name", result.Name)
	assert.NotNil(t, result.GetCreatedAt())
	assert.NotNil(t, result.GetUpdatedAt())
}

func TestModelServiceV2_InsertOne(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Arrange
	svc := service.NewModelServiceV2[TestModel]()
	testModel := TestModel{Name: "Test Name"}

	// Act
	id, err := svc.InsertOne(testModel)

	// Assert
	assert.Nil(t, err)
	assert.NotEqual(t, primitive.NilObjectID, id)
}

func TestModelServiceV2_UpdateById(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Arrange
	svc := service.NewModelServiceV2[TestModel]()
	id := primitive.NewObjectID()
	testModel := TestModel{Id: id, Name: "Old Name"}
	_, err := svc.InsertOne(testModel)
	assert.Nil(t, err)

	// Act
	update := bson.M{"$set": bson.M{"name": "New Name"}}
	err = svc.UpdateById(id, update)

	// Assert
	assert.Nil(t, err)
	result, err := svc.GetById(id)
	assert.Nil(t, err)
	assert.Equal(t, "New Name", result.Name)
}

func TestModelServiceV2_DeleteById(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Arrange
	svc := service.NewModelServiceV2[TestModel]()
	id := primitive.NewObjectID()
	testModel := TestModel{Id: id, Name: "Test Name"}
	_, err := svc.InsertOne(testModel)
	assert.Nil(t, err)

	// Act
	err = svc.DeleteById(id)

	// Assert
	assert.Nil(t, err)
	result, err := svc.GetById(id)
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestModelServiceV2_GetList(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Arrange
	svc := service.NewModelServiceV2[TestModel]()
	testModels := []TestModel{
		{Name: "Name1"},
		{Name: "Name2"},
	}
	_, err := svc.InsertMany(testModels)
	assert.Nil(t, err)

	// Act
	results, err := svc.GetList(bson.M{}, nil)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 2, len(results))
}

func TestModelServiceV2_Count(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Arrange
	svc := service.NewModelServiceV2[TestModel]()
	testModels := []TestModel{
		{Name: "Name1"},
		{Name: "Name2"},
	}
	_, err := svc.InsertMany(testModels)
	assert.Nil(t, err)

	// Act
	total, err := svc.Count(bson.M{})

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 2, total)
}
