package handler

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/client"
	"github.com/crawlab-team/crawlab-core/task"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	interfaces.TaskBaseService
	nodeCfgSvc     interfaces.NodeConfigService
	modelSvc       interfaces.GrpcClientModelService
	modelNodeSvc   interfaces.GrpcClientModelNodeService
	modelSpiderSvc interfaces.GrpcClientModelSpiderService
	modelTaskSvc   interfaces.GrpcClientModelTaskService

	// settings
	maxRunners        int
	exitWatchDuration time.Duration
	reportInterval    time.Duration

	// internals variables
	stopped      bool
	mu           sync.Mutex
	runnersCount int      // number of task runners
	runners      sync.Map // pool of task runners started
}

func (svc *Service) Run(taskId primitive.ObjectID) (err error) {
	// validate if there are available runners
	if svc.runnersCount >= svc.maxRunners {
		return errors.ErrorTaskNoAvailableRunners
	}

	// attempt to get runner from pool
	_, ok := svc.runners.Load(taskId)
	if ok {
		return errors.ErrorTaskAlreadyExists
	}

	// create a new task runner
	r, err := NewTaskRunner(taskId, svc)
	if err != nil {
		return err
	}

	// add runner to pool
	svc.addRunner(taskId, r)

	// create a goroutine to run task
	go func() {
		// run task process (blocking)
		// error or finish after task runner ends
		if err := r.Run(); err != nil {
			switch err {
			case constants.ErrTaskError:
				log.Errorf("task[%s] finished with error: %v", r.GetTaskId().Hex(), err)
			case constants.ErrTaskCancelled:
				log.Errorf("task[%s] cancelled", r.GetTaskId().Hex())
			default:
				log.Errorf("task[%s] finished with unknown error: %v", r.GetTaskId().Hex(), err)
			}

			// delete runner from pool
			svc.deleteRunner(r.GetTaskId())
		}
		log.Infof("task[%s] finished", r.GetTaskId().Hex())

		// delete runner from pool
		svc.deleteRunner(r.GetTaskId())
	}()

	return nil
}

func (svc *Service) Cancel(taskId primitive.ObjectID) (err error) {
	r, err := svc.getRunner(taskId)
	if err != nil {
		return err
	}
	if err := r.Cancel(); err != nil {
		return err
	}
	return nil
}

func (svc *Service) ReportHandlerStatus() {
	for {
		if svc.stopped {
			return
		}

		// report handler status
		if err := svc.reportHandlerStatus(); err != nil {
			trace.PrintError(err)
		}

		// wait
		time.Sleep(svc.reportInterval)
	}
}

func (svc *Service) GetMaxRunners() (maxRunners int) {
	return svc.maxRunners
}

func (svc *Service) SetMaxRunners(maxRunners int) {
	svc.maxRunners = maxRunners
}

func (svc *Service) GetExitWatchDuration() (duration time.Duration) {
	return svc.exitWatchDuration
}

func (svc *Service) SetExitWatchDuration(duration time.Duration) {
	svc.exitWatchDuration = duration
}

func (svc *Service) GetReportInterval() (interval time.Duration) {
	return svc.reportInterval
}

func (svc *Service) SetReportInterval(interval time.Duration) {
	svc.reportInterval = interval
}

func (svc *Service) GetModelService() (modelSvc interfaces.GrpcClientModelService) {
	return svc.modelSvc
}

func (svc *Service) GetModelSpiderService() (modelSpiderSvc interfaces.GrpcClientModelSpiderService) {
	return svc.modelSpiderSvc
}

func (svc *Service) GetModelTaskService() (modelTaskSvc interfaces.GrpcClientModelTaskService) {
	return svc.modelTaskSvc
}

func (svc *Service) getRunner(taskId primitive.ObjectID) (r interfaces.TaskRunner, err error) {
	v, ok := svc.runners.Load(taskId)
	if !ok {
		return nil, errors.ErrorTaskNotExists
	}
	switch v.(type) {
	case interfaces.TaskRunner:
		r = v.(interfaces.TaskRunner)
	default:
		return nil, errors.ErrorModelInvalidType
	}
	return r, nil
}

func (svc *Service) addRunner(taskId primitive.ObjectID, r interfaces.TaskRunner) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.runners.Store(taskId, r)
	svc.runnersCount++
}

func (svc *Service) deleteRunner(taskId primitive.ObjectID) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.runners.Delete(taskId)
	svc.runnersCount--
}

func (svc *Service) saveTask(t interfaces.Task, status string) (err error) {
	// normalize status
	if status == "" {
		status = constants.TaskStatusPending
	}

	// set task status
	t.SetStatus(status)

	// attempt to get task from database
	_, err = svc.modelTaskSvc.GetTaskById(t.GetId())
	if err != nil {
		// if task does not exist, add to database
		if err == mongo.ErrNoDocuments {
			if err := client.NewModelDelegate(t).Add(); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	} else {
		// otherwise, update
		if err := client.NewModelDelegate(t).Save(); err != nil {
			return err
		}
		return nil
	}
}

func (svc *Service) reportHandlerStatus() (err error) {
	nodeKey := svc.nodeCfgSvc.GetNodeKey()
	n, err := svc.modelNodeSvc.GetNodeByKey(nodeKey)
	if err != nil {
		return err
	}
	ar := svc.maxRunners - svc.runnersCount
	n.SetAvailableRunners(ar)
	n.SetMaxRunners(svc.maxRunners)
	if err := client.NewModelNodeDelegate(n).Save(); err != nil {
		return err
	}
	return nil
}

func NewTaskHandlerService(opts ...Option) (svc2 interfaces.TaskHandlerService, err error) {
	// base service
	baseSvc, err := task.NewBaseService()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// service
	// Service
	svc := &Service{
		TaskBaseService:   baseSvc,
		maxRunners:        8,
		exitWatchDuration: 60 * time.Second,
		reportInterval:    60 * time.Second,
		mu:                sync.Mutex{},
		runnersCount:      0,
		runners:           sync.Map{},
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(client.ProvideServiceDelegate(svc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.NewNodeServiceDelegate); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.NewSpiderServiceDelegate); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.NewTaskServiceDelegate); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(
		modelSvc interfaces.GrpcClientModelService,
		modelNodeSvc interfaces.GrpcClientModelNodeService,
		modelSpiderSvc interfaces.GrpcClientModelSpiderService,
		modelTaskSvc interfaces.GrpcClientModelTaskService,
	) {
		svc.modelSvc = modelSvc
		svc.modelNodeSvc = modelNodeSvc
		svc.modelSpiderSvc = modelSpiderSvc
		svc.modelTaskSvc = modelTaskSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}

func ProvideTaskHandlerService(path string, opts ...Option) func() (svc interfaces.TaskHandlerService, err error) {
	opts = append(opts, WithConfigPath(path))
	return func() (svc interfaces.TaskHandlerService, err error) {
		return NewTaskHandlerService(opts...)
	}
}
