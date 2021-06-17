package controllers

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"net/http"
)

var UserController *userController

var UserActions = []Action{
	{
		Method:      http.MethodPost,
		Path:        "/:id/change-password",
		HandlerFunc: userCtx.changePassword,
	},
}

type userController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *userContext
}

func (ctr *userController) Put(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := ctr.ctx.userSvc.Create(&interfaces.UserCreateOptions{
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		Role:     u.Role,
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

var userCtx = newUserContext()

type userContext struct {
	modelSvc service.ModelService
	userSvc  interfaces.UserService
}

func (ctx *userContext) changePassword(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	var payload models.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if len(payload.Password) < 6 {
		HandleErrorBadRequest(c, errors.ErrorUserInvalidPassword)
		return
	}
	if err := ctx.userSvc.ChangePassword(id, payload.Password); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *userContext) me(c *gin.Context) {
	res, ok := c.Get(constants.UserContextKey)
	if !ok {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}
	u, ok := res.(interfaces.User)
	if !ok {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}
	HandleSuccessWithData(c, u)
}

func newUserContext() *userContext {
	// context
	ctx := &userContext{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		panic(err)
	}
	if err := c.Provide(user.ProvideGetUserService()); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		userSvc interfaces.UserService,
	) {
		ctx.modelSvc = modelSvc
		ctx.userSvc = userSvc
	}); err != nil {
		panic(err)
	}

	return ctx
}

func newUserController() *userController {
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdUser, modelSvc.GetBaseService(interfaces.ModelIdUser), UserActions)
	d := NewListPostActionControllerDelegate(ControllerIdUser, modelSvc.GetBaseService(interfaces.ModelIdUser), UserActions)
	ctx := newUserContext()

	return &userController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
