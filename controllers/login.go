package controllers

import (
	"github.com/gin-gonic/gin"
)

var LoginController = loginController{}

type loginController struct {
}

func (ctr *loginController) Post(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "login":
		ctr.login(c)
	case "logout":
		ctr.logout(c)
	}
}

func (ctr *loginController) login(c *gin.Context) {
}

func (ctr *loginController) logout(c *gin.Context) {
}
