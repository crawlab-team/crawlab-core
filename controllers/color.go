package controllers

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetColorList(c *gin.Context) {
	panic(errors.ErrorControllerNotImplemented)
}

var ColorActions = []Action{
	{Method: http.MethodGet, Path: "", HandlerFunc: GetColorList},
}

var ColorController ActionController
