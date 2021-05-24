package scheduler

import (
	"encoding/json"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/task"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"go.uber.org/dig"
)

type Service struct {
	// dependencies
	interfaces.TaskBaseService
	nodeCfgSvc interfaces.NodeConfigService
	modelSvc   service.ModelService

	// internals
	ch chan interfaces.Task
}

func (svc *Service) Start() {
	go svc.Fetch()
	go svc.Assign()
	svc.Wait()
	svc.Stop()
}

func (svc *Service) Fetch() {
	for {
		// return if quit is true
		if svc.IsStopped() {
			return
		}

		// fetch task with retry
		if err := backoff.RetryNotify(func() error {
			// fetch
			t, err := svc.fetch()
			if err != nil {
				return trace.TraceError(err)
			}

			// notify task channel
			svc.ch <- t

			return nil
		}, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("task scheduler fetch task")); err != nil {
			trace.PrintError(err)
		}
	}
}

func (svc *Service) Assign() {
	for {
		// return if quit is true
		if svc.IsStopped() {
			return
		}

		// receive task from channel
		t := <-svc.ch

		// assign task
		if err := backoff.RetryNotify(func() error {
			return svc.assign(t)
		}, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("task scheduler assign task")); err != nil {
			trace.PrintError(err)
		}
	}
}

func (svc *Service) fetch() (t interfaces.Task, err error) {
	// message
	var msg string

	// TODO: implement fetch priority queue message

	// deserialization
	tMsg := entity.TaskMessage{}
	if err := json.Unmarshal([]byte(msg), &tMsg); err != nil {
		return t, err
	}

	// fetch task
	t, err = svc.modelSvc.GetTaskById(tMsg.Id)
	if err != nil {
		return t, err
	}

	return t, nil
}

func (svc *Service) assign(t interfaces.Task) (err error) {
	// TODO: implement
	if err := svc.SaveTask(t, constants.TaskStatusAssigned); err != nil {
		return err
	}
	log.Infof("task scheduler assigned task[%s] successfully", t.GetId())
	return nil
}

func NewTaskSchedulerService(opts ...Option) (svc2 interfaces.TaskSchedulerService, err error) {
	// base service
	baseSvc, err := task.NewBaseService()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// service
	svc := &Service{TaskBaseService: baseSvc}

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
	if err := c.Invoke(func(nodeCfgSvc interfaces.NodeConfigService, modelSvc service.ModelService) {
		svc.nodeCfgSvc = nodeCfgSvc
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// task channel
	svc.ch = make(chan interfaces.Task)

	return svc, nil
}

func ProvideTaskSchedulerService(path string, opts ...Option) func() (svc interfaces.TaskSchedulerService, err error) {
	opts = append(opts, WithConfigPath(path))
	return func() (svc interfaces.TaskSchedulerService, err error) {
		return NewTaskSchedulerService(opts...)
	}
}
