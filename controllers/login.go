package controllers

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	panic(errors.ErrorControllerNotImplemented)
}

func Logout(c *gin.Context) {
	panic(errors.ErrorControllerNotImplemented)
}

var LoginActions = []Action{
	{Method: http.MethodPost, Path: "/login", HandlerFunc: Login},
	{Method: http.MethodPost, Path: "/logout", HandlerFunc: Logout},
}

var LoginController ActionController
