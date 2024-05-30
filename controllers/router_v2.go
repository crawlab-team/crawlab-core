package controllers

import (
	"github.com/crawlab-team/crawlab-core/middlewares"
	"github.com/gin-gonic/gin"
)

type RouterGroups struct {
	AuthGroup      *gin.RouterGroup
	AnonymousGroup *gin.RouterGroup
}

func NewRouterGroups(app *gin.Engine) (groups *RouterGroups) {
	return &RouterGroups{
		AuthGroup:      app.Group("/", middlewares.AuthorizationMiddleware()),
		AnonymousGroup: app.Group("/"),
	}
}

func RegisterController[T any](group *gin.RouterGroup, basePath string, ctr *BaseControllerV2[T]) {
	group.GET(basePath, ctr.GetList)
	group.GET(basePath+"/:id", ctr.GetById)
	group.POST(basePath, ctr.Post)
	group.PUT(basePath+"/:id", ctr.PutById)
	group.PATCH(basePath, ctr.PatchList)
	group.DELETE(basePath+"/:id", ctr.DeleteById)
}

func InitRoutes(app *gin.Engine) {
	// routes groups
	groups := NewRouterGroups(app)

	RegisterController(groups.AuthGroup, "/projects", ProjectV2Controller)
}
