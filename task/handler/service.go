package handler

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/client"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	nodeCfgSvc     interfaces.NodeConfigService
	modelSvc       interfaces.GrpcClientModelService
	modelSpiderSvc interfaces.GrpcClientModelSpiderService
	modelTaskSvc   interfaces.GrpcClientModelTaskService

	// settings variables
	cfgPath           string
	maxRunners        int
	exitWatchDuration time.Duration

	// internals variables
	mu           sync.Mutex
	runnersCount int      // number of task runners
	runners      sync.Map // pool of task runners started
}

func (svc *Service) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *Service) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *Service) Init() (err error) {
	// TODO: implement
	return nil
}

func (svc *Service) Start() {
	panic("implement me")
}

func (svc *Service) Wait() {
	utils.DefaultWait()
}

func (svc *Service) Stop() {
	panic("implement me")
}

func (svc *Service) Run(taskId primitive.ObjectID) (err error) {
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
	svc.AddRunner(taskId, r)

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
			return
		}
		log.Infof("task[%s] finished", r.GetTaskId().Hex())
	}()

	return nil
}

func (svc *Service) Cancel(taskId primitive.ObjectID) (err error) {
	r, err := svc.getTaskRunner(taskId)
	if err != nil {
		return err
	}
	if err := r.Cancel(); err != nil {
		return err
	}
	return nil
}

func (svc *Service) AddRunner(taskId primitive.ObjectID, r interfaces.TaskRunner) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.runners.Store(taskId, r)
	svc.runnersCount++
}

func (svc *Service) DeleteRunner(taskId primitive.ObjectID) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.runners.Delete(taskId)
	svc.runnersCount--
}

func (svc *Service) GetRunner(taskId primitive.ObjectID) (r interfaces.TaskRunner, err error) {
	res, ok := svc.runners.Load(taskId)
	if !ok {
		return nil, trace.TraceError(errors.ErrorTaskNotExists)
	}
	r, ok = res.(interfaces.TaskRunner)
	if !ok {
		return nil, trace.TraceError(errors.ErrorTaskInvalidType)
	}
	return r, nil
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

func (svc *Service) GetModelService() (modelSvc interfaces.GrpcClientModelService) {
	return svc.modelSvc
}

func (svc *Service) GetModelSpiderService() (modelSpiderSvc interfaces.GrpcClientModelSpiderService) {
	return svc.modelSpiderSvc
}

func (svc *Service) GetModelTaskService() (modelTaskSvc interfaces.GrpcClientModelTaskService) {
	return svc.modelTaskSvc
}

func (svc *Service) getTaskRunner(taskId primitive.ObjectID) (r interfaces.TaskRunner, err error) {
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

func (svc *Service) saveTask(t interfaces.Task, status string) (err error) {
	// normalize status
	if status == "" {
		status = constants.StatusPending
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

func NewTaskHandlerService(opts ...Option) (svc2 interfaces.TaskHandlerService, err error) {
	// construct Service
	svc := &Service{
		maxRunners:        8,
		exitWatchDuration: 60 * time.Second,
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
	if err := c.Provide(config.NewNodeConfigService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(nodeCfgSvc interfaces.NodeConfigService) {
		svc.nodeCfgSvc = nodeCfgSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.ProvideServiceDelegate(svc.nodeCfgSvc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.NewSpiderServiceDelegate); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.NewTaskServiceDelegate); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(modelSvc interfaces.GrpcClientModelService, modelSpiderSvc interfaces.GrpcClientModelSpiderService, modelTaskSvc interfaces.GrpcClientModelTaskService) {
		svc.modelSvc = modelSvc
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
