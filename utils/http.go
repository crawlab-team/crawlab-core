package utils

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func HandleError(statusCode int, c *gin.Context, err error) {
	c.AbortWithStatusJSON(statusCode, entity.Response{
		Status:  "error",
		Message: "failure",
		Error:   err.Error(),
	})
}

func HandleErrorF(statusCode int, c *gin.Context, err string) {
	debug.PrintStack()
	c.AbortWithStatusJSON(statusCode, entity.Response{
		Status:  "ok",
		Message: "error",
		Error:   err,
	})
}

func HandleSuccess(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, entity.Response{
		Status:  "ok",
		Message: "success",
	})
}

func HandleSuccessData(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, entity.Response{
		Status:  "ok",
		Message: "success",
		Data:    data,
	})
}
