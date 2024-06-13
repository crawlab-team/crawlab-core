package controllers_test

import (
	"bytes"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/middlewares"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateSpider(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.POST("/spiders", controllers.PostSpider)

	payload := models.SpiderV2{
		Name:    "Test Spider",
		ColName: "test_spiders",
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/spiders", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", TestToken)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response controllers.Response[models.SpiderV2]
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.Nil(t, err)
	assert.False(t, response.Data.Id.IsZero())
	assert.Equal(t, payload.Name, response.Data.Name)
}

func TestGetSpiderById(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.GET("/spiders/:id", controllers.GetSpiderById)

	id := primitive.NewObjectID()
	model := models.SpiderV2{
		Name:    "Test Spider",
		ColName: "test_spiders",
	}
	model.SetId(id)
	jsonValue, _ := json.Marshal(model)
	_, err := http.NewRequest("POST", "/spiders", bytes.NewBuffer(jsonValue))
	require.Nil(t, err)
	time.Sleep(100 * time.Millisecond)

	req, _ := http.NewRequest("GET", "/spiders/"+id.Hex(), nil)
	req.Header.Set("Authorization", TestToken)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response controllers.Response[models.SpiderV2]
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	require.Nil(t, err)
	assert.Equal(t, model.Name, response.Data.Name)
}

func TestUpdateSpiderById(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.PUT("/spiders/:id", controllers.PutSpiderById)

	id := primitive.NewObjectID()
	model := models.SpiderV2{
		Name:    "Test Spider",
		ColName: "test_spiders",
	}
	model.SetId(id)
	jsonValue, _ := json.Marshal(model)
	_, err := http.NewRequest("POST", "/spiders", bytes.NewBuffer(jsonValue))
	require.Nil(t, err)

	spiderId := id.Hex()
	payload := models.SpiderV2{
		Name: "Updated Spider",
	}
	jsonValue, _ = json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/spiders/"+spiderId, bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", TestToken)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response controllers.Response[models.SpiderV2]
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	require.Nil(t, err)
	assert.Equal(t, payload.Name, response.Data.Name)

	svc := service.NewModelServiceV2[models.SpiderV2]()
	resModel, err := svc.GetById(id)
	require.Nil(t, err)
	assert.Equal(t, payload.Name, resModel.Name)
}

func TestDeleteSpiderById(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.DELETE("/spiders/:id", controllers.DeleteSpiderById)

	id := primitive.NewObjectID()
	model := models.SpiderV2{
		Name:    "Test Spider",
		ColName: "test_spiders",
	}
	model.SetId(id)
	jsonValue, _ := json.Marshal(model)
	_, err := http.NewRequest("POST", "/spiders", bytes.NewBuffer(jsonValue))
	require.Nil(t, err)

	req, _ := http.NewRequest("DELETE", "/spiders/"+id.Hex(), nil)
	req.Header.Set("Authorization", TestToken)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	svc := service.NewModelServiceV2[models.SpiderV2]()
	_, err = svc.GetById(id)
	require.NotNil(t, err)
}

func TestDeleteSpiderList(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.DELETE("/spiders", controllers.DeleteSpiderList)

	svc := service.NewModelServiceV2[models.SpiderV2]()
	modelList := []models.SpiderV2{
		{
			Name:    "Test Name 1",
			ColName: "test_spiders",
		}, {
			Name:    "Test Name 2",
			ColName: "test_spiders",
		},
	}
	var ids []primitive.ObjectID
	for _, model := range modelList {
		id := primitive.NewObjectID()
		model.SetId(id)
		jsonValue, _ := json.Marshal(model)
		_, err := http.NewRequest("POST", "/spiders", bytes.NewBuffer(jsonValue))
		require.Nil(t, err)
		ids = append(ids, id)
	}

	payload := struct {
		Ids []string `json:"ids"`
	}{}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("DELETE", "/spiders", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", TestToken)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	total, err := svc.Count(bson.M{"_id": bson.M{"$in": ids}})
	require.Nil(t, err)
	require.Equal(t, 0, total)
}
