package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getSystemInfo(c *gin.Context) {
	// TODO: implement me
	panic("not implemented")
}

func getSystemInfoActions() []Action {
	return []Action{
		{
			Path:        "",
			Method:      http.MethodGet,
			HandlerFunc: getSystemInfo,
		},
	}
}

var SystemInfoController ActionController
