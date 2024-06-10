package controllers_test

import (
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/middlewares"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostUserChangePassword_Success(t *testing.T) {
	modelSvc := service.NewModelServiceV2[models.UserV2]()
	u := models.UserV2{
		Id: primitive.NewObjectID(),
	}
	_, err := modelSvc.InsertOne(u)
	assert.Nil(t, err)

	userSvc, err := user.GetUserServiceV2()
	require.Nil(t, err)
	token, err := userSvc.MakeToken(&u)
	require.Nil(t, err)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.POST("/users/:id/change-password", controllers.PostUserChangePassword)

	id := u.Id
	password := "newPassword"
	reqBody := strings.NewReader(`{"password":"` + password + `"}`)
	req, _ := http.NewRequest(http.MethodPost, "/users/"+id.Hex()+"/change-password", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserMe_Success(t *testing.T) {
	modelSvc := service.NewModelServiceV2[models.UserV2]()
	u := models.UserV2{
		Id: primitive.NewObjectID(),
	}
	_, err := modelSvc.InsertOne(u)
	assert.Nil(t, err)

	userSvc, err := user.GetUserServiceV2()
	require.Nil(t, err)
	token, err := userSvc.MakeToken(&u)
	require.Nil(t, err)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.GET("/users/me", controllers.GetUserMe)

	req, _ := http.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutUserById_Success(t *testing.T) {
	modelSvc := service.NewModelServiceV2[models.UserV2]()
	u := models.UserV2{
		Id: primitive.NewObjectID(),
	}
	_, err := modelSvc.InsertOne(u)
	assert.Nil(t, err)

	userSvc, err := user.GetUserServiceV2()
	require.Nil(t, err)
	token, err := userSvc.MakeToken(&u)
	require.Nil(t, err)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.PUT("/users/me", controllers.PutUserById)

	id := primitive.NewObjectID()
	reqBody := strings.NewReader(`{"id":"` + id.Hex() + `","username":"newUsername","email":"newEmail@test.com"}`)
	req, _ := http.NewRequest(http.MethodPut, "/users/me", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
