package controllers

import (
	"github.com/gin-gonic/gin"
)

var UserController = userController{}

type userController struct {
}

func (ctr *userController) Get(c *gin.Context) {
	panic("implement me")
}

func (ctr *userController) GetList(c *gin.Context) {
	panic("implement me")
}

func (ctr *userController) Post(c *gin.Context) {
	panic("implement me")
}

func (ctr *userController) PostList(c *gin.Context) {
	panic("implement me")
}

func (ctr *userController) Put(c *gin.Context) {
	panic("implement me")
}

func (ctr *userController) PutList(c *gin.Context) {
	panic("implement me")
}

func (ctr *userController) Delete(c *gin.Context) {
	panic("implement me")
}

func (ctr *userController) DeleteList(c *gin.Context) {
	panic("implement me")
}
