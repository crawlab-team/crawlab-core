package middlewares

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token string
		tokenStr := c.GetHeader("Authorization")

		// validate token
		user, err := services.CheckToken(tokenStr)

		// validation failed, return error response
		if err != nil {
			controllers.HandleErrorF(http.StatusUnauthorized, c, errors.AuthUnauthorized)
			return
		}

		// set user in context
		c.Set(constants.ContextUser, &user)

		// validation success
		c.Next()
	}
}
