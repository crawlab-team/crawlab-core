package manager

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/grpc/server"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/task"
	"github.com/crawlab-team/crawlab-db/redis"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
)

type Service struct {
	// dependencies
	interfaces.TaskBaseService
	nodeCfgSvc interfaces.NodeConfigService
	svr        interfaces.GrpcServer
}

func (svc *Service) Enqueue(t interfaces.Task) (err error) {
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
	if err := svc.SaveTask(t, constants.TaskStatusPending); err != nil {
		return err
	}

	return nil
}

func (svc *Service) Cancel(taskId primitive.ObjectID) (err error) {
	// TODO: implement
	return nil
}

func NewTaskManagerService(opts ...Option) (svc2 interfaces.TaskManagerService, err error) {
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
	if err := c.Provide(server.ProvideServer(svc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(nodeCfgSvc interfaces.NodeConfigService, svr interfaces.GrpcServer) {
		svc.nodeCfgSvc = nodeCfgSvc
		svc.svr = svr
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}

func ProvideTaskManagerService(path string, opts ...Option) func() (svc interfaces.TaskManagerService, err error) {
	opts = append(opts, WithConfigPath(path))
	return func() (svc interfaces.TaskManagerService, err error) {
		return NewTaskManagerService(opts...)
	}
}
