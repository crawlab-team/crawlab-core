package controllers

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleSuccessWithDataV2[T any](c *gin.Context, data T) {
	c.AbortWithStatusJSON(http.StatusOK, entity.Response{
		Status:  constants.HttpResponseStatusOk,
		Message: constants.HttpResponseMessageSuccess,
		Data:    data,
	})
}

func HandleSuccessWithListDataV2[T any](c *gin.Context, data []T, total int) {
	c.AbortWithStatusJSON(http.StatusOK, entity.ListResponse{
		Status:  constants.HttpResponseStatusOk,
		Message: constants.HttpResponseMessageSuccess,
		Data:    data,
		Total:   total,
	})
}
