package test

import (
	"context"
	"github.com/crawlab-team/crawlab-core/entity"
	gtest "github.com/crawlab-team/crawlab-core/grpc/test"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	ntest "github.com/crawlab-team/crawlab-core/node/test"
	stest "github.com/crawlab-team/crawlab-core/spider/test"
	"github.com/crawlab-team/crawlab-core/task/handler"
	"github.com/crawlab-team/crawlab-core/task/scheduler"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"io/ioutil"
	"os"
	"path"
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
	schedulerSvc  interfaces.TaskSchedulerService
	handlerSvc    interfaces.TaskHandlerService
	modelSvc      service.ModelService
	masterFsSvc   interfaces.SpiderFsService
	workerFsSvc   interfaces.SpiderFsService
	masterSyncSvc interfaces.SpiderSyncService
	client        interfaces.GrpcClient
	server        interfaces.GrpcServer
	sub           grpc.NodeService_SubscribeClient

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
	//if err := delegate.NewModelDelegate(t.TestSpider).Add(); err != nil {
	//	panic(err)
	//}

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
	ntest.T.StartMasterWorker()
}

func (t *Test) StopMasterWorker() {
	ntest.T.StopMasterWorker()
}

func NewTest() (res *Test, err error) {
	t := &Test{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(scheduler.ProvideTaskSchedulerService(ntest.T.MasterSvc.GetConfigPath(), scheduler.WithInterval(5*time.Second))); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(handler.ProvideTaskHandlerService(ntest.T.WorkerSvc.GetConfigPath())); err != nil {
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
	t.masterFsSvc = stest.T.GetMasterFsSvc()
	t.workerFsSvc = stest.T.GetWorkerFsSvc()
	t.masterSyncSvc = stest.T.GetMasterSyncSvc()

	// test node
	t.TestNode = &models.Node{
		Id:               primitive.NewObjectID(),
		Key:              "test_node_key",
		Enabled:          true,
		Active:           true,
		AvailableRunners: 20,
	}

	// test spider
	t.TestSpider = stest.T.TestSpider

	// test task
	t.TestTask = t.NewTask()

	// test task with node id
	t.TestTaskWithNodeId = t.NewTask()
	t.TestTaskWithNodeId.SetNodeId(t.TestNode.GetId())

	// add file to spider fs
	filePath := path.Join(t.masterFsSvc.GetWorkspacePath(), stest.T.ScriptName)
	if err := ioutil.WriteFile(filePath, []byte(stest.T.Script), os.ModePerm); err != nil {
		panic(err)
	}
	if err := t.masterFsSvc.GetFsService().Commit("initial commit"); err != nil {
		return nil, err
	}
	if err := t.masterSyncSvc.SyncToFs(t.TestSpider.GetId()); err != nil {
		panic(err)
	}

	// grpc server/client
	grpcT, _ := gtest.NewTest()
	t.server = grpcT.Server
	t.client = grpcT.Client
	if err := t.client.Start(); err != nil {
		return nil, err
	}
	req := &grpc.Request{NodeKey: t.TestNode.GetKey()}
	t.sub, err = t.client.GetNodeClient().Subscribe(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return t, nil
}
