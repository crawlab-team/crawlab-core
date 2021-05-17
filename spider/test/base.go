package test

import (
	fstest "github.com/crawlab-team/crawlab-core/fs/test"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	ntest "github.com/crawlab-team/crawlab-core/node/test"
	"github.com/crawlab-team/crawlab-core/spider/admin"
	"github.com/crawlab-team/crawlab-core/spider/sync"
	"github.com/crawlab-team/go-trace"
	"go.uber.org/dig"
	"os"
	"testing"
	"time"
)

func init() {
	// remove tmp directory
	if _, err := os.Stat("./tmp"); err == nil {
		if err := os.RemoveAll("./tmp"); err != nil {
			panic(err)
		}
	}
	if err := os.MkdirAll("./tmp", os.ModePerm); err != nil {
		panic(err)
	}

	var err error
	T, err = NewTest()
	if err != nil {
		panic(err)
	}
}

var T *Test

type Test struct {
	s        *models.Spider
	script   string
	adminSvc interfaces.SpiderAdminService
	syncSvc  interfaces.SpiderSyncService
	modelSvc service.ModelService
}

// Setup spider fs service test setup
func (t *Test) Setup(t2 *testing.T) {
	t2.Cleanup(t.Cleanup)
}

// Cleanup spider fs service test cleanup
func (t *Test) Cleanup() {
	// fs service cleanup
	fstest.T.Cleanup()

	// wait to avoid caching
	time.Sleep(500 * time.Millisecond)
}

func NewTest() (res *Test, err error) {
	// test
	t := &Test{
		s: &models.Spider{
			Name: "test_spider",
			Cmd:  "go run main.go",
		},
		script: `package main
import "fmt"
func main() {
  fmt.Println("it works")
}`,
	}

	// add spider to db
	if err := delegate.NewModelDelegate(t.s).Add(); err != nil {
		return nil, err
	}

	// spider service
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(admin.ProvideSpiderAdminService(ntest.T.MasterSvc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(sync.ProvideSpiderSyncService(ntest.T.WorkerSvc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(modelSvc service.ModelService, adminSvc interfaces.SpiderAdminService, sync interfaces.SpiderSyncService) {
		t.adminSvc = adminSvc
		t.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return t, nil
}
