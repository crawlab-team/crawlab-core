package middlewares

import (
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testPaginationMiddlewareHandlerFunc(c *gin.Context) {
	p, err := controllers.GetPagination(c)
	if err != nil {
		controllers.HandleErrorBadRequest(c, err)
		return
	}
	controllers.HandleSuccessData(c, p)
}

func TestPaginationMiddleware(t *testing.T) {
	app := gin.New()

	app.Use(PaginationMiddleware())

	app.GET("/test", testPaginationMiddlewareHandlerFunc)

	server := httptest.NewServer(app)
	defer server.Close()

	p := entity.Pagination{
		Page: 2,
		Size: 20,
	}

	e := httpexpect.New(t, server.URL)
	data := e.GET("/test").
		WithQueryObject(p).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Path("$.data").Object()
	data.Path("$.page").Number().Equal(2)
	data.Path("$.size").Number().Equal(20)
}
