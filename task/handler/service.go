package handler

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/client"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/task"
	"github.com/crawlab-team/go-trace"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	interfaces.TaskBaseService
	cfgSvc                 interfaces.NodeConfigService
	modelSvc               service.ModelService
	clientModelSvc         interfaces.GrpcClientModelService
	clientModelNodeSvc     interfaces.GrpcClientModelNodeService
	clientModelSpiderSvc   interfaces.GrpcClientModelSpiderService
	clientModelTaskSvc     interfaces.GrpcClientModelTaskService
	clientModelTaskStatSvc interfaces.GrpcClientModelTaskStatService

	// settings
	maxRunners        int
	exitWatchDuration time.Duration
	reportInterval    time.Duration

	// internals variables
	stopped      bool
	mu           sync.Mutex
	runnersCount int      // number of task runners
	runners      sync.Map // pool of task runners started
	syncLocks    sync.Map // files sync locks map of task runners
}

func (svc *Service) Start() {
	go svc.ReportStatus()
}

func (svc *Service) Run(taskId primitive.ObjectID) (err error) {
	// validate if there are available runners
	if svc.runnersCount >= svc.maxRunners {
		return trace.TraceError(errors.ErrorTaskNoAvailableRunners)
	}

	// attempt to get runner from pool
	_, ok := svc.runners.Load(taskId)
	if ok {
		return trace.TraceError(errors.ErrorTaskAlreadyExists)
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

func (svc *Service) Reset() {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.runners = sync.Map{}
	svc.runnersCount = 0
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

func (svc *Service) ReportStatus() {
	for {
		if svc.stopped {
			return
		}

		// report handler status
		if err := svc.reportStatus(); err != nil {
			trace.PrintError(err)
		}

		// wait
		time.Sleep(svc.reportInterval)
	}
}

func (svc *Service) IsSyncLocked(spiderId primitive.ObjectID) (ok bool) {
	_, ok = svc.syncLocks.Load(spiderId)
	return ok
}

func (svc *Service) LockSync(spiderId primitive.ObjectID) {
	svc.syncLocks.Store(spiderId, true)
}

func (svc *Service) UnlockSync(spiderId primitive.ObjectID) {
	svc.syncLocks.Delete(spiderId)
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
	return svc.clientModelSvc
}

func (svc *Service) GetModelSpiderService() (modelSpiderSvc interfaces.GrpcClientModelSpiderService) {
	return svc.clientModelSpiderSvc
}

func (svc *Service) GetModelTaskService() (modelTaskSvc interfaces.GrpcClientModelTaskService) {
	return svc.clientModelTaskSvc
}

func (svc *Service) GetModelTaskStatService() (modelTaskSvc interfaces.GrpcClientModelTaskStatService) {
	return svc.clientModelTaskStatSvc
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
	if svc.runnersCount < svc.maxRunners {
		svc.runnersCount++
	}
}

func (svc *Service) deleteRunner(taskId primitive.ObjectID) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.runners.Delete(taskId)
	if svc.runnersCount > 0 {
		svc.runnersCount--
	}
}

func (svc *Service) saveTask(t interfaces.Task, status string) (err error) {
	// normalize status
	if status == "" {
		status = constants.TaskStatusPending
	}

	// set task status
	t.SetStatus(status)

	// attempt to get task from database
	_, err = svc.clientModelTaskSvc.GetTaskById(t.GetId())
	if err != nil {
		// if task does not exist, add to database
		if err == mongo.ErrNoDocuments {
			if err := client.NewModelDelegate(t, client.WithDelegateConfigPath(svc.GetConfigPath())).Add(); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	} else {
		// otherwise, update
		if err := client.NewModelDelegate(t, client.WithDelegateConfigPath(svc.GetConfigPath())).Save(); err != nil {
			return err
		}
		return nil
	}
}

func (svc *Service) reportStatus() (err error) {
	// node key
	nodeKey := svc.cfgSvc.GetNodeKey()

	// current node
	var n interfaces.Node
	if svc.cfgSvc.IsMaster() {
		n, err = svc.modelSvc.GetNodeByKey(nodeKey, nil)
	} else {
		n, err = svc.clientModelNodeSvc.GetNodeByKey(nodeKey)
	}
	if err != nil {
		return err
	}

	// update node
	ar := svc.maxRunners - svc.runnersCount
	n.SetAvailableRunners(ar)
	n.SetMaxRunners(svc.maxRunners)

	// save node
	if svc.cfgSvc.IsMaster() {
		err = client.NewModelDelegate(n, client.WithDelegateConfigPath(svc.GetConfigPath())).Save()
	} else {
		err = delegate.NewModelDelegate(n).Save()
	}
	if err != nil {
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
	svc := &Service{
		TaskBaseService:   baseSvc,
		maxRunners:        8,
		exitWatchDuration: 60 * time.Second,
		reportInterval:    60 * time.Second,
		mu:                sync.Mutex{},
		runnersCount:      0,
		runners:           sync.Map{},
		syncLocks:         sync.Map{},
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
	if err := c.Provide(client.ProvideServiceDelegate(svc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.ProvideNodeServiceDelegate(svc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.ProvideSpiderServiceDelegate(svc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.ProvideTaskServiceDelegate(svc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.ProvideTaskStatServiceDelegate(svc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(
		cfgSvc interfaces.NodeConfigService,
		modelSvc service.ModelService,
		clientModelSvc interfaces.GrpcClientModelService,
		clientModelNodeSvc interfaces.GrpcClientModelNodeService,
		clientModelSpiderSvc interfaces.GrpcClientModelSpiderService,
		clientModelTaskSvc interfaces.GrpcClientModelTaskService,
		clientModelTaskStatSvc interfaces.GrpcClientModelTaskStatService,
	) {
		svc.cfgSvc = cfgSvc
		svc.modelSvc = modelSvc
		svc.clientModelSvc = clientModelSvc
		svc.clientModelNodeSvc = clientModelNodeSvc
		svc.clientModelSpiderSvc = clientModelSpiderSvc
		svc.clientModelTaskSvc = clientModelTaskSvc
		svc.clientModelTaskStatSvc = clientModelTaskStatSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}

func ProvideTaskHandlerService(path string, opts ...Option) func() (svc interfaces.TaskHandlerService, err error) {
	// config path
	opts = append(opts, WithConfigPath(path))

	// max runners
	maxRunners := viper.GetInt("task.handler.maxRunners")
	if maxRunners > 0 {
		opts = append(opts, WithMaxRunners(maxRunners))
	}

	// report interval
	reportIntervalSeconds := viper.GetInt("task.handler.reportInterval")
	if reportIntervalSeconds > 0 {
		opts = append(opts, WithReportInterval(time.Duration(reportIntervalSeconds)*time.Second))
	}

	return func() (svc interfaces.TaskHandlerService, err error) {
		return NewTaskHandlerService(opts...)
	}
}
