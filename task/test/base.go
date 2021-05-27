package test

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/test"
	"github.com/crawlab-team/crawlab-core/task/handler"
	"github.com/crawlab-team/crawlab-core/task/scheduler"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"testing"
	"time"
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
	// dependencies
	schedulerSvc interfaces.TaskSchedulerService
	handlerSvc   interfaces.TaskHandlerService
	modelSvc     service.ModelService

	// data
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

	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	_ = t.modelSvc.DropAll()
}

func (t *Test) NewTask() (t2 interfaces.Task) {
	return &models.Task{
		SpiderId: t.TestSpider.GetId(),
	}
}

func (t *Test) StartMasterWorker() {
	test.T.StartMasterWorker()
}

func (t *Test) StopMasterWorker() {
	test.T.StopMasterWorker()
}

func NewTest() (res *Test, err error) {
	t := &Test{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(scheduler.ProvideTaskSchedulerService(test.T.MasterSvc.GetConfigPath(), scheduler.WithInterval(5*time.Second))); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(handler.ProvideTaskHandlerService(test.T.MasterSvc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(service.NewService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(
		schedulerSvc interfaces.TaskSchedulerService,
		handlerSvc interfaces.TaskHandlerService,
		modelSvc service.ModelService,
	) {
		t.schedulerSvc = schedulerSvc
		t.handlerSvc = handlerSvc
		t.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// test node
	t.TestNode = &models.Node{
		Id:               primitive.NewObjectID(),
		Key:              "test_node_key",
		Enabled:          true,
		Active:           true,
		AvailableRunners: 20,
	}

	// test spider
	t.TestSpider = &models.Spider{
		Id: primitive.NewObjectID(),
	}

	// test task
	t.TestTask = t.NewTask()

	// test task with node id
	t.TestTaskWithNodeId = t.NewTask()
	t.TestTaskWithNodeId.SetNodeId(t.TestNode.GetId())

	return t, nil
}
