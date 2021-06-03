package scheduler

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/grpc/server"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/client"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/task"
	"github.com/crawlab-team/crawlab-core/task/handler"
	"github.com/crawlab-team/crawlab-db/mongo"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"github.com/joeshaw/multierror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	interfaces.TaskBaseService
	nodeCfgSvc interfaces.NodeConfigService
	modelSvc   service.ModelService
	svr        interfaces.GrpcServer
	handlerSvc interfaces.TaskHandlerService

	// settings
	interval time.Duration
}

func (svc *Service) Start() {
	go svc.DequeueAndSchedule()
	svc.Wait()
	svc.Stop()
}

func (svc *Service) Enqueue(t interfaces.Task) (err error) {
	// set task status
	t.SetStatus(constants.TaskStatusPending)

	// run mongo transaction
	if err := mongo.RunTransaction(func(sc mongo2.SessionContext) error {
		var err error

		// add task
		if err = delegate.NewModelDelegate(t).Add(); err != nil {
			return err
		}

		// task queue item
		tq := &models.TaskQueueItem{
			Id:       t.GetId(),
			Priority: t.GetPriority(),
		}

		// task stat
		ts := &models.TaskStat{
			Id: t.GetId(),
		}

		// enqueue task
		_, err = mongo.GetMongoCol(interfaces.ModelColNameTaskQueue).Insert(tq)
		if err != nil {
			return trace.TraceError(err)
		}

		// add task stat
		_, err = mongo.GetMongoCol(interfaces.ModelColNameTaskStat).Insert(ts)
		if err != nil {
			return trace.TraceError(err)
		}

		return nil
	}); err != nil {
		return trace.TraceError(err)
	}

	// success
	return nil
}

func (svc *Service) DequeueAndSchedule() {
	for {
		if svc.IsStopped() {
			return
		}

		// wait
		time.Sleep(svc.interval)

		if err := mongo.RunTransaction(func(sc mongo2.SessionContext) error {
			// dequeue tasks
			tasks, err := svc.Dequeue()
			if err != nil {
				return trace.TraceError(err)
			}

			// schedule tasks
			if err := svc.Schedule(tasks); err != nil {
				return trace.TraceError(err)
			}

			return nil
		}); err != nil {
			trace.PrintError(err)
		}
	}
}

func (svc *Service) Dequeue() (tasks []interfaces.Task, err error) {
	// get task queue items
	tqList, err := svc.getTaskQueueItems()
	if err != nil {
		return nil, err
	}
	if tqList == nil {
		return nil, nil
	}

	// match resources
	tasks, nodesMap, err := svc.matchResources(tqList)
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		return nil, nil
	}

	// update resources
	if err := svc.updateResources(nodesMap); err != nil {
		return nil, err
	}

	// dequeue tasks
	if err := svc.dequeueTasks(tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (svc *Service) Schedule(tasks []interfaces.Task) (err error) {
	var e multierror.Errors

	// nodes cache
	nodesCache := sync.Map{}

	// wait group
	wg := sync.WaitGroup{}
	wg.Add(len(tasks))

	// iterate tasks and execute each of them
	for _, t := range tasks {
		go func(t interfaces.Task) {
			var err error

			// node of the task
			var n interfaces.Node
			res, ok := nodesCache.Load(t.GetNodeId())
			if !ok {
				// not exists in cache
				n, err = svc.modelSvc.GetNodeById(t.GetNodeId())
				if err == nil {
					nodesCache.Store(n.GetId(), n)
				}
			} else {
				// exists in cache
				n, ok = res.(interfaces.Node)
				if !ok {
					err = errors.ErrorTaskInvalidType
				}
			}
			if err != nil {
				e = append(e, err)
				svc.handleTaskError(n, t, err)
				wg.Done()
				return
			}

			// schedule task
			if n.GetIsMaster() {
				// execute task on master
				err = svc.handlerSvc.Run(t.GetId())
			} else {
				// send to execute task on worker nodes
				err = svc.svr.SendStreamMessageWithData(n.GetKey(), grpc.StreamMessageCode_RUN_TASK, t)
			}
			if err != nil {
				e = append(e, err)
				svc.handleTaskError(n, t, err)
				wg.Done()
				return
			}

			// success
			wg.Done()
		}(t)
	}

	// wait
	wg.Wait()

	return e.Err()
}

func (svc *Service) SetInterval(interval time.Duration) {
	svc.interval = interval
}

func (svc *Service) getTaskQueueItems() (tqList []models.TaskQueueItem, err error) {
	opts := &mongo.FindOptions{
		Sort: bson.M{
			"p": -1,
		},
	}
	if err := mongo.GetMongoCol(interfaces.ModelColNameTaskQueue).Find(nil, opts).All(&tqList); err != nil {
		if err == mongo2.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return tqList, nil
}

func (svc *Service) getResourcesAndNodesMap() (resources map[string]interfaces.Node, nodesMap map[primitive.ObjectID]interfaces.Node, err error) {
	nodesMap = map[primitive.ObjectID]interfaces.Node{}
	resources = map[string]interfaces.Node{}
	query := bson.M{
		// enabled: true
		"en": true,
		// active: true
		"a": true,
		// available_runners > 0
		"ar": bson.M{
			"$gt": 0,
		},
	}
	nodes, err := svc.modelSvc.GetNodeList(query, nil)
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	for _, n := range nodes {
		nodesMap[n.Id] = &n
		for i := 0; i < n.AvailableRunners; i++ {
			key := fmt.Sprintf("%s:%d", n.Id.Hex(), i)
			resources[key] = &n
		}
	}
	return resources, nodesMap, nil
}

func (svc *Service) matchResources(tqList []models.TaskQueueItem) (tasks []interfaces.Task, nodesMap map[primitive.ObjectID]interfaces.Node, err error) {
	resources, nodesMap, err := svc.getResourcesAndNodesMap()
	if err != nil {
		return nil, nil, err
	}
	if resources == nil || len(resources) == 0 {
		return nil, nil, nil
	}
	for _, tq := range tqList {
		t, err := svc.modelSvc.GetTaskById(tq.GetId())
		if err != nil {
			return nil, nil, err
		}
		for key, r := range resources {
			if t.GetNodeId().IsZero() ||
				t.GetNodeId() == r.GetId() {
				t.NodeId = r.GetId()
				tasks = append(tasks, t)
				delete(resources, key)

				n := nodesMap[r.GetId()]
				n.DecrementAvailableRunners()
				break
			}
		}
	}

	return tasks, nodesMap, nil
}

func (svc *Service) updateResources(nodesMap map[primitive.ObjectID]interfaces.Node) (err error) {
	for _, n := range nodesMap {
		if err := delegate.NewModelNodeDelegate(n).Save(); err != nil {
			return err
		}
	}
	return nil
}

func (svc *Service) dequeueTasks(tasks []interfaces.Task) (err error) {
	for _, t := range tasks {
		// save task with node id
		if err := delegate.NewModelDelegate(t).Save(); err != nil {
			return err
		}

		// remove task queue item
		if err := mongo.GetMongoCol(interfaces.ModelColNameTaskQueue).DeleteId(t.GetId()); err != nil {
			return err
		}
	}
	return nil
}

func (svc *Service) handleTaskError(n interfaces.Node, t interfaces.Task, err error) {
	trace.PrintError(err)
	t.SetStatus(constants.TaskStatusError)
	t.SetError(err.Error())
	if n.GetIsMaster() {
		_ = delegate.NewModelDelegate(t).Save()
	} else {
		_ = client.NewModelDelegate(t).Save()
	}
}

func NewTaskSchedulerService(opts ...Option) (svc2 interfaces.TaskSchedulerService, err error) {
	// base service
	baseSvc, err := task.NewBaseService()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// service
	svc := &Service{
		TaskBaseService: baseSvc,
		interval:        15 * time.Second,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(config.ProvideConfigService(svc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(service.NewService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(server.ProvideGetServer(svc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(handler.ProvideTaskHandlerService(svc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(
		nodeCfgSvc interfaces.NodeConfigService,
		modelSvc service.ModelService,
		svr interfaces.GrpcServer,
		handlerSvc interfaces.TaskHandlerService,
	) {
		svc.nodeCfgSvc = nodeCfgSvc
		svc.modelSvc = modelSvc
		svc.svr = svr
		svc.handlerSvc = handlerSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}

func ProvideTaskSchedulerService(path string, opts ...Option) func() (svc interfaces.TaskSchedulerService, err error) {
	opts = append(opts, WithConfigPath(path))
	return func() (svc interfaces.TaskSchedulerService, err error) {
		return NewTaskSchedulerService(opts...)
	}
}

var store = sync.Map{}

func GetTaskSchedulerService(path string, opts ...Option) (svr interfaces.TaskSchedulerService, err error) {
	if path == "" {
		path = config.DefaultConfigPath
	}
	opts = append(opts, WithConfigPath(path))
	res, ok := store.Load(path)
	if !ok {
		return NewTaskSchedulerService(opts...)
	}
	svr, ok = res.(interfaces.TaskSchedulerService)
	if !ok {
		return NewTaskSchedulerService(opts...)
	}
	return svr, nil
}

func ProvideGetTaskSchedulerService(path string, opts ...Option) func() (svr interfaces.TaskSchedulerService, err error) {
	return func() (svr interfaces.TaskSchedulerService, err error) {
		return GetTaskSchedulerService(path, opts...)
	}
}
