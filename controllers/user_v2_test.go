package controllers_test

import (
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostUserChangePassword_Success(t *testing.T) {
	modelSvc := service.NewModelServiceV2[models.UserV2]()
	user := models.UserV2{
		Id: primitive.NewObjectID(),
	}
	_, err := modelSvc.InsertOne(user)
	assert.Nil(t, err)

	router := gin.Default()
	router.POST("/users/:id/change-password", controllers.PostUserChangePassword)

	id := user.Id
	password := "newPassword"
	reqBody := strings.NewReader(`{"password":"` + password + `"}`)
	req, _ := http.NewRequest(http.MethodPost, "/users/"+id.Hex()+"/change-password", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserMe_Success(t *testing.T) {
	router := gin.Default()
	router.GET("/users/me", controllers.GetUserMe)

	req, _ := http.NewRequest(http.MethodGet, "/users/me", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutUserById_Success(t *testing.T) {
	router := gin.Default()
	router.PUT("/users/me", controllers.PutUserById)

	id := primitive.NewObjectID()
	reqBody := strings.NewReader(`{"id":"` + id.Hex() + `","username":"newUsername","email":"newEmail@test.com"}`)
	req, _ := http.NewRequest(http.MethodPut, "/users/me", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
