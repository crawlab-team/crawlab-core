package middlewares

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/gin-gonic/gin"
)

func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var p entity.Pagination
		if err := c.ShouldBindQuery(&p); err != nil {
			c.Next()
			return
		}
		c.Set(constants.PAGINATION, &p)
	}
}
