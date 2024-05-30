package controllers_test

import (
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterGroups_AuthGroup(t *testing.T) {
	router := gin.Default()
	groups := controllers.NewRouterGroups(router)

	assert.NotNil(t, groups.AuthGroup)
}

func TestRouterGroups_AnonymousGroup(t *testing.T) {
	router := gin.Default()
	groups := controllers.NewRouterGroups(router)

	assert.NotNil(t, groups.AnonymousGroup)
}

func TestRegisterController_Routes(t *testing.T) {
	router := gin.Default()
	groups := controllers.NewRouterGroups(router)
	ctr := controllers.NewControllerV2[TestModel]()
	basePath := "/testmodels"

	controllers.RegisterController(groups.AuthGroup, basePath, ctr)

	// Check if all routes are registered
	routes := router.Routes()

	assert.Equal(t, 6, len(routes))
	assert.Contains(t, routes, gin.RouteInfo{Method: "GET", Path: basePath})
	assert.Contains(t, routes, gin.RouteInfo{Method: "GET", Path: basePath + "/:id"})
	assert.Contains(t, routes, gin.RouteInfo{Method: "POST", Path: basePath})
	assert.Contains(t, routes, gin.RouteInfo{Method: "PUT", Path: basePath + "/:id"})
	assert.Contains(t, routes, gin.RouteInfo{Method: "PATCH", Path: basePath})
	assert.Contains(t, routes, gin.RouteInfo{Method: "DELETE", Path: basePath + "/:id"})
}

func TestInitRoutes_ProjectsRoute(t *testing.T) {
	router := gin.Default()

	controllers.InitRoutes(router)

	// Check if the projects route is registered
	routes := router.Routes()

	assert.Contains(t, routes, gin.RouteInfo{Method: "GET", Path: "/projects"})
	assert.Contains(t, routes, gin.RouteInfo{Method: "GET", Path: "/projects/:id"})
	assert.Contains(t, routes, gin.RouteInfo{Method: "POST", Path: "/projects"})
	assert.Contains(t, routes, gin.RouteInfo{Method: "PUT", Path: "/projects/:id"})
	assert.Contains(t, routes, gin.RouteInfo{Method: "PATCH", Path: "/projects"})
	assert.Contains(t, routes, gin.RouteInfo{Method: "DELETE", Path: "/projects/:id"})
}

func TestInitRoutes_UnauthorizedAccess(t *testing.T) {
	router := gin.Default()

	controllers.InitRoutes(router)

	// Create a test request
	req, _ := http.NewRequest("GET", "/projects", nil)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
