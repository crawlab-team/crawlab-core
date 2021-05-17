package scheduler

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/crawlab-db/redis"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
)

type Service struct {
	// dependencies
	nodeCfgSvc interfaces.NodeConfigService
	modelSvc   service.ModelService

	// settings variables
	cfgPath string
}

func (svc *Service) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *Service) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *Service) Init() error {
	panic("implement me")
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

func (svc *Service) Assign(t interfaces.Task) (err error) {
	// TODO: implement assign priority queue message
	// TODO: implement task assign via grpc
	// validate node type
	if !svc.nodeCfgSvc.IsMaster() {
		return errors.ErrorTaskForbidden
	}

	// task message
	msg := entity.TaskMessage{
		Id: t.GetId(),
	}

	// serialization
	msgStr, err := msg.ToString()
	if err != nil {
		return err
	}

	// queue name
	var queue string
	if t.GetNodeId().IsZero() {
		queue = "tasks:public"
	} else {
		queue = "tasks:node:" + t.GetNodeId().Hex()
	}

	// enqueue
	if err := redis.RedisClient.RPush(queue, msgStr); err != nil {
		return err
	}

	// set task status as "pending" and save to database
	if err := svc.saveTask(t, constants.StatusPending); err != nil {
		return err
	}

	return nil
}

func (svc *Service) Fetch() (err error) {
	panic("implement me")
}

func (svc *Service) Run(taskId primitive.ObjectID) (err error) {
	// TODO: implement
	return nil
}

func (svc *Service) Cancel(taskId primitive.ObjectID) (err error) {
	// TODO: implement
	return nil
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

func (svc *Service) saveTask(t interfaces.Task, status string) (err error) {
	// normalize status
	if status == "" {
		status = constants.StatusPending
	}

	// set task status
	t.SetStatus(status)

	// attempt to get task from database
	_, err = svc.modelSvc.GetTaskById(t.GetId())
	if err != nil {
		// if task does not exist, add to database
		if err == mongo.ErrNoDocuments {
			if err := delegate.NewModelDelegate(t).Add(); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	} else {
		// otherwise, update
		if err := delegate.NewModelDelegate(t).Save(); err != nil {
			return err
		}
		return nil
	}
}

func NewTaskSchedulerService(opts ...Option) (svc2 interfaces.TaskSchedulerService, err error) {
	// service
	svc := &Service{}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(config.NewNodeConfigService); err != nil {
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

	return svc, nil
}

func ProvideTaskSchedulerService(path string, opts ...Option) func() (svc interfaces.TaskSchedulerService, err error) {
	return func() (svc interfaces.TaskSchedulerService, err error) {
		return NewTaskSchedulerService(opts...)
	}
}
