package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestModel is a simple struct to be used as a model in tests
type TestModel struct {
	Id   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" collection:"test_collection"`
	Name string             `bson:"name" json:"name"`
}

// SetupTestDB sets up the test database
func SetupTestDB() {
	viper.Set("mongo.db", "testdb")
}

// SetupRouter sets up the gin router for testing
func SetupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// CleanupTestDB cleans up the test database
func CleanupTestDB() {
	mongo.GetMongoDb("testdb").Drop(context.Background())
}

func TestBaseControllerV2_GetById(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	// Insert a test document
	id := primitive.NewObjectID()
	_, err := mongo.GetMongoCol("test_collection").Insert(bson.M{"_id": id, "name": "test"})
	assert.NoError(t, err)

	// Initialize the controller
	ctr := controllers.NewControllerV2[TestModel]()

	// Set up the router
	router := SetupRouter()
	router.GET("/testmodels/:id", ctr.GetById)

	// Create a test request
	req, _ := http.NewRequest("GET", "/testmodels/"+id.Hex(), nil)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.Response[TestModel]
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test", response.Data.Name)
}

func TestBaseControllerV2_Post(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	// Initialize the controller
	ctr := controllers.NewControllerV2[TestModel]()

	// Set up the router
	router := SetupRouter()
	router.POST("/testmodels", ctr.Post)

	// Create a test request
	testModel := TestModel{Name: "test"}
	jsonValue, _ := json.Marshal(testModel)
	req, _ := http.NewRequest("POST", "/testmodels", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.Response[TestModel]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test", response.Data.Name)

	// Check if the document was inserted into the database
	var result TestModel
	err = mongo.GetMongoCol("test_collection").Find(bson.M{"_id": response.Data.Id}, nil).One(&result)
	assert.NoError(t, err)
	assert.Equal(t, "test", result.Name)
}

func TestBaseControllerV2_DeleteById(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	// Insert a test document
	id := primitive.NewObjectID()
	_, err := mongo.GetMongoCol("test_collection").Insert(bson.M{"_id": id, "name": "test"})
	assert.NoError(t, err)

	// Initialize the controller
	ctr := controllers.NewControllerV2[TestModel]()

	// Set up the router
	router := SetupRouter()
	router.DELETE("/testmodels/:id", ctr.DeleteById)

	// Create a test request
	req, _ := http.NewRequest("DELETE", "/testmodels/"+id.Hex(), nil)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the document was deleted from the database
	err = mongo.GetMongoCol("test_collection").Find(bson.M{"_id": id}, nil).One(&TestModel{})
	assert.Error(t, err)
}
