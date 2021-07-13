package middlewares

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func FilerAuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// auth key
		authKey := c.GetHeader("Authorization")

		// server auth key
		svrAuthKey := viper.GetString("fs.filer.authKey")
		if svrAuthKey == "" {
			svrAuthKey = constants.DefaultFilerAuthKey
		}

		// validate
		if authKey != svrAuthKey {
			// validation failed, return error response
			controllers.HandleErrorUnauthorized(c, errors.ErrorHttpUnauthorized)
			return
		}

		// validation success
		c.Next()
	}
}
