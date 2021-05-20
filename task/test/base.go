package test

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/test"
	"github.com/crawlab-team/crawlab-core/task/handler"
	"github.com/crawlab-team/crawlab-core/task/scheduler"
	"github.com/crawlab-team/go-trace"
	"go.uber.org/dig"
	"testing"
)

func init() {
	var err error
	T, err = NewTest()
	if err != nil {
		panic(err)
	}
}

var T *Test

type Test struct {
	schedulerSvc interfaces.TaskSchedulerService
	handlerSvc   interfaces.TaskHandlerService
}

func (t *Test) Setup(t2 *testing.T) {
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
}

func NewTest() (res *Test, err error) {
	t := &Test{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(scheduler.ProvideTaskSchedulerService(test.T.MasterSvc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(handler.ProvideTaskHandlerService(test.T.MasterSvc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(schedulerSvc interfaces.TaskSchedulerService, handlerSvc interfaces.TaskHandlerService) {
		t.schedulerSvc = schedulerSvc
		t.handlerSvc = handlerSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return t, nil
}
