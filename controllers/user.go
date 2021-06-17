package controllers

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"net/http"
)

var UserController *userController

var UserActions = []Action{
	{
		Method:      http.MethodPost,
		Path:        "/login",
		HandlerFunc: userCtx.login,
	},
	{
		Method:      http.MethodPost,
		Path:        "/logout",
		HandlerFunc: userCtx.logout,
	},
	{
		Method:      http.MethodGet,
		Path:        "/me",
		HandlerFunc: userCtx.me,
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

func (ctx *userContext) login(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	token, loggedInUser, err := ctx.userSvc.Login(&interfaces.UserLoginOptions{
		Username: u.Username,
		Password: u.Password,
	})
	if err != nil {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}
	c.Set(constants.UserContextKey, loggedInUser)
	HandleSuccessWithData(c, token)
}

func (ctx *userContext) logout(c *gin.Context) {
	c.Set(constants.UserContextKey, nil)
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
