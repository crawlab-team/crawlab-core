package controllers

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

var UserControllerV2 = NewControllerV2[models.UserV2](
	Action{
		Method:      http.MethodPost,
		Path:        "/:id/change-password",
		HandlerFunc: PostUserChangePassword,
	},
	Action{
		Method:      http.MethodGet,
		Path:        "/me",
		HandlerFunc: GetUserMe,
	},
	Action{
		Method:      http.MethodPut,
		Path:        "/me",
		HandlerFunc: PutUserById,
	},
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
	res, ok := c.Get(constants.UserContextKey)
	if !ok {
		HandleErrorUnauthorized(c, errors.ErrorUserNotExists)
		return
	}
	u, ok := res.(models.UserV2)
	if !ok {
		HandleErrorUnauthorized(c, errors.ErrorUserNotExists)
		return
	}
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
	res, ok := c.Get(constants.UserContextKey)
	if !ok {
		HandleErrorUnauthorized(c, errors.ErrorUserNotExists)
		return
	}
	u, ok := res.(models.UserV2)
	if !ok {
		HandleErrorUnauthorized(c, errors.ErrorUserNotExists)
		return
	}
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
	res, ok := c.Get(constants.UserContextKey)
	if !ok {
		HandleErrorUnauthorized(c, errors.ErrorUserNotExists)
		return
	}
	u, ok := res.(models.UserV2)
	if !ok {
		HandleErrorUnauthorized(c, errors.ErrorUserNotExists)
		return
	}
	modelSvc := service.NewModelServiceV2[models.UserV2]()

	// update user
	userDb, err := modelSvc.GetById(u.Id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	user.Password = userDb.Password
	user.SetUpdated(u.Id)
	if err := modelSvc.ReplaceById(u.Id, user); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// handle success
	HandleSuccess(c)
}
