package controllers

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/gin-gonic/gin"
)

func GetDefaultPagination() (p *entity.Pagination) {
	return &entity.Pagination{
		Page: constants.PaginationDeafultPage,
		Size: constants.PaginationDefaultSize,
	}
}

func GetPagination(c *gin.Context) (p *entity.Pagination, err error) {
	if err := c.ShouldBindQuery(&p); err != nil {
		return GetDefaultPagination(), err
	}
	return p, nil
}

func MustGetPagination(c *gin.Context) (p *entity.Pagination) {
	if err := c.ShouldBindQuery(&p); err != nil {
		return GetDefaultPagination()
	}
	return p
}
