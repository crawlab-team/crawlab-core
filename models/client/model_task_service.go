package client

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
)

type TaskServiceDelegate struct {
	interfaces.GrpcClientModelBaseService
}

func (svc *TaskServiceDelegate) GetTaskById(id primitive.ObjectID) (t interfaces.Task, err error) {
	res, err := svc.GetById(id)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Task)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *TaskServiceDelegate) GetTask(query bson.M, opts *mongo.FindOptions) (t interfaces.Task, err error) {
	res, err := svc.Get(query, opts)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Task)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *TaskServiceDelegate) GetTaskList(query bson.M, opts *mongo.FindOptions) (res []interfaces.Task, err error) {
	list, err := svc.GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, item := range list.Values() {
		s, ok := item.(interfaces.Task)
		if !ok {
			return nil, errors.ErrorModelInvalidType
		}
		res = append(res, s)
	}
	return res, nil
}

func NewTaskServiceDelegate() (svc2 interfaces.GrpcClientModelTaskService, err error) {
	svc := &TaskServiceDelegate{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(ProvideBaseServiceDelegate(interfaces.ModelIdTask)); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(baseSvc interfaces.GrpcClientModelBaseService) {
		svc.GrpcClientModelBaseService = baseSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}
