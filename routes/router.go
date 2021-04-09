package routes

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RouterServiceInterface interface {
	RegisterControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.ListController)
	RegisterHandlerToGroup(group *gin.RouterGroup, path string, method string, handler gin.HandlerFunc)
}

type RouterService struct {
	app *gin.Engine
}

func NewRouterService(app *gin.Engine) (svc *RouterService) {
	return &RouterService{
		app: app,
	}
}

func (svc *RouterService) RegisterControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.BasicController) {
	group.GET(basePath, ctr.Get)
	group.PUT(basePath, ctr.Put)
	group.POST(basePath, ctr.Post)
	group.DELETE(basePath, ctr.Delete)
}

func (svc *RouterService) RegisterListControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.ListController) {
	group.GET(basePath+"/:id", ctr.Get)
	group.GET(basePath, ctr.GetList)
	group.PUT(basePath, ctr.Put)
	group.PUT(basePath+"/batch", ctr.PutList)
	group.POST(basePath+"/:id", ctr.Post)
	group.POST(basePath, ctr.PostList)
	group.DELETE(basePath+"/:id", ctr.Delete)
	group.DELETE(basePath, ctr.DeleteList)
}

func (svc *RouterService) RegisterPostActionControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.PostActionController) {
	group.POST(basePath+"/:action", func(c *gin.Context) {
		param := c.Param("action")
		for _, action := range ctr.Actions() {
			if action.Name == param {
				action.HandlerFunc(c)
				return
			}
		}
		controllers.HandleErrorNotFound(c, errors.ErrorHttpNotFound)
	})
}

func (svc *RouterService) RegisterListPostActionControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.ListPostActionController) {
	svc.RegisterControllerToGroup(group, basePath, ctr)
	svc.RegisterPostActionControllerToGroup(group, basePath, ctr)
}

func (svc *RouterService) RegisterHandlerToGroup(group *gin.RouterGroup, path string, method string, handler gin.HandlerFunc) {
	switch method {
	case http.MethodGet:
		group.GET(path, handler)
	case http.MethodPut:
		group.PUT(path, handler)
	case http.MethodPost:
		group.POST(path, handler)
	case http.MethodDelete:
		group.DELETE(path, handler)
	default:
		log.Warn(fmt.Sprintf("%s is not a valid http method", method))
	}
}

func InitRoutes(app *gin.Engine) (err error) {
	// routes groups
	groups := NewRouterGroups(app)

	// router service
	svc := NewRouterService(app)

	// login/logout
	svc.RegisterPostActionControllerToGroup(groups.AnonymousGroup, "/auth", controllers.AuthController)

	// project
	svc.RegisterListControllerToGroup(groups.AuthGroup, "/projects", controllers.ProjectController)

	// user
	svc.RegisterListControllerToGroup(groups.AuthGroup, "/users", controllers.UserController)

	// spider
	svc.RegisterListControllerToGroup(groups.AuthGroup, "/spiders", controllers.SpiderController)

	return nil
}
