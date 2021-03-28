package controllers

import (
	"github.com/gin-gonic/gin"
)

var LoginController = loginController{}

type loginController struct {
}

func (ctr *loginController) Login(c *gin.Context) {
}

func (ctr *loginController) Logout(c *gin.Context) {
}
