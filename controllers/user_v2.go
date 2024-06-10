package controllers

import (
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostUserChangePassword(c *gin.Context) {
	// get id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// get payload
	var payload struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// get user
	u := GetUserFromContextV2(c)
	modelSvc := service.NewModelServiceV2[models.UserV2]()

	// update password
	user, err := modelSvc.GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	user.SetUpdated(u.Id)
	user.Password = utils.EncryptMd5(payload.Password)
	if err := modelSvc.ReplaceById(user.Id, *user); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// handle success
	HandleSuccess(c)
}

func GetUserMe(c *gin.Context) {
	u := GetUserFromContextV2(c)
	HandleSuccessWithData(c, u)
}

func PutUserById(c *gin.Context) {
	// get payload
	var user models.UserV2
	if err := c.ShouldBindJSON(&user); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// get user
	u := GetUserFromContextV2(c)

	modelSvc := service.NewModelServiceV2[models.UserV2]()

	// update user
	userDb, err := modelSvc.GetById(u.Id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	user.Password = userDb.Password
	user.SetUpdated(u.Id)
	if user.Id.IsZero() {
		user.Id = u.Id
	}
	if err := modelSvc.ReplaceById(u.Id, user); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// handle success
	HandleSuccess(c)
}
