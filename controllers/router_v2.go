package controllers

import (
	"github.com/crawlab-team/crawlab-core/middlewares"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/gin-gonic/gin"
	"net/http"
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

	RegisterController(groups.AuthGroup, "/data/collections", NewControllerV2[models.DataCollectionV2]())
	RegisterController(groups.AuthGroup, "/data-sources", NewControllerV2[models.DataSourceV2]())
	RegisterController(groups.AuthGroup, "/environments", NewControllerV2[models.EnvironmentV2]())
	RegisterController(groups.AuthGroup, "/gits", NewControllerV2[models.GitV2]())
	RegisterController(groups.AuthGroup, "/nodes", NewControllerV2[models.NodeV2]())
	RegisterController(groups.AuthGroup, "/notifications/settings", NewControllerV2[models.SettingV2]())
	RegisterController(groups.AuthGroup, "/permissions", NewControllerV2[models.PermissionV2]())
	RegisterController(groups.AuthGroup, "/projects", NewControllerV2[models.ProjectV2]())
	RegisterController(groups.AuthGroup, "/roles", NewControllerV2[models.RoleV2]())
	RegisterController(groups.AuthGroup, "/schedules", NewControllerV2[models.ScheduleV2](
	// TODO: implement actions
	))
	RegisterController(groups.AuthGroup, "/settings", NewControllerV2[models.SettingV2]())
	RegisterController(groups.AuthGroup, "/spiders", NewControllerV2[models.SpiderV2](
	// TODO: implement actions
	))
	RegisterController(groups.AuthGroup, "/tasks", NewControllerV2[models.TaskV2](
	// TODO: implement actions
	))
	RegisterController(groups.AuthGroup, "/tokens", NewControllerV2[models.TokenV2](
	// TODO: implement actions
	))
	RegisterController(groups.AuthGroup, "/users", NewControllerV2[models.UserV2](
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/change-password",
			HandlerFunc: PostUserChangePassword,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "/me",
			HandlerFunc: GetUserMe,
		},
		Action{
			Method:      http.MethodPut,
			Path:        "/me",
			HandlerFunc: PutUserById,
		},
	))
}
