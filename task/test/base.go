package test

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/test"
	"github.com/crawlab-team/crawlab-core/task/handler"
	"github.com/crawlab-team/crawlab-core/task/manager"
	"github.com/crawlab-team/crawlab-core/task/scheduler"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	managerSvc         interfaces.TaskManagerService
	schedulerSvc       interfaces.TaskSchedulerService
	handlerSvc         interfaces.TaskHandlerService
	modelSvc           service.ModelService
	TestNode           interfaces.Node
	TestSpider         interfaces.Spider
	TestTask           interfaces.Task
	TestTaskWithNodeId interfaces.Task
	TestTaskMessage    entity.TaskMessage
}

func (t *Test) Setup(t2 *testing.T) {
	// add test node
	if err := delegate.NewModelDelegate(t.TestNode).Add(); err != nil {
		panic(err)
	}
	// add test spider
	if err := delegate.NewModelDelegate(t.TestSpider).Add(); err != nil {
		panic(err)
	}

	// add test task
	if err := delegate.NewModelDelegate(t.TestTask).Add(); err != nil {
		panic(err)
	}

	// add test task with node id
	if err := delegate.NewModelDelegate(t.TestTaskWithNodeId).Add(); err != nil {
		panic(err)
	}

	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	_ = t.modelSvc.DropAll()
}

func NewTest() (res *Test, err error) {
	t := &Test{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(manager.ProvideTaskManagerService(test.T.MasterSvc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(scheduler.ProvideTaskSchedulerService(test.T.MasterSvc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(handler.ProvideTaskHandlerService(test.T.MasterSvc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(service.NewService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(managerSvc interfaces.TaskManagerService, schedulerSvc interfaces.TaskSchedulerService, handlerSvc interfaces.TaskHandlerService, modelSvc service.ModelService) {
		t.managerSvc = managerSvc
		t.schedulerSvc = schedulerSvc
		t.handlerSvc = handlerSvc
		t.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// test node
	t.TestNode = &models.Node{
		Id: primitive.NewObjectID(),
	}

	// test spider
	t.TestSpider = &models.Spider{
		Id: primitive.NewObjectID(),
	}

	// test task
	t.TestTask = &models.Task{
		Id:       primitive.NewObjectID(),
		SpiderId: primitive.NewObjectID(),
	}

	// test task with node id
	t.TestTaskWithNodeId = &models.Task{
		Id:       primitive.NewObjectID(),
		SpiderId: primitive.NewObjectID(),
		NodeId:   t.TestNode.GetId(),
	}

	return t, nil
}
