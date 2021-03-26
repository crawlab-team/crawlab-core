package routes

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RouterServiceInterface interface {
	RegisterControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.Controller)
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

func (svc *RouterService) RegisterControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.Controller) {
	group.GET(basePath, ctr.Get)
	group.GET(basePath, ctr.GetList)
	group.PUT(basePath, ctr.Put)
	group.PUT(basePath, ctr.PutList)
	group.POST(basePath, ctr.Post)
	group.POST(basePath, ctr.PostList)
	group.DELETE(basePath, ctr.Delete)
	group.DELETE(basePath, ctr.DeleteList)
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

	// login
	svc.RegisterHandlerToGroup(groups.AnonymousGroup, "/login", http.MethodPost, controllers.PostLogin)

	// project
	svc.RegisterControllerToGroup(groups.AuthGroup, "/projects", &controllers.ProjectController)

	return nil
}
