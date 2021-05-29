package test

import (
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/routes"
	"github.com/crawlab-team/go-trace"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	if err := controllers.InitControllers(); err != nil {
		panic(err)
	}
	var err error
	T, err = NewTest()
	if err != nil {
		panic(err)
	}
}

type Test struct {
	// dependencies
	modelSvc service.ModelService

	// internals
	app *gin.Engine
	svr *httptest.Server
}

func (t *Test) Setup(t2 *testing.T) {
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	_ = t.modelSvc.DropAll()
	time.Sleep(200 * time.Millisecond)
}

func (t *Test) NewExpect(t2 *testing.T) (e *httpexpect.Expect) {
	return httpexpect.New(t2, t.svr.URL)
}

var T *Test

func NewTest() (res *Test, err error) {
	// test
	t := &Test{}

	// gin app
	t.app = gin.New()

	// http test server
	t.svr = httptest.NewServer(t.app)

	// init routes
	if err := routes.InitRoutes(t.app); err != nil {
		return nil, err
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(modelSvc service.ModelService) {
		t.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return t, nil
}
