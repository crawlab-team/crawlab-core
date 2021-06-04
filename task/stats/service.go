package stats

import (
	"github.com/apex/log"
	config2 "github.com/crawlab-team/crawlab-core/config"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/client"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/task"
	clog "github.com/crawlab-team/crawlab-log"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"sync"
)

type Service struct {
	// dependencies
	interfaces.TaskBaseService
	nodeCfgSvc         interfaces.NodeConfigService
	clientModelSvc     interfaces.GrpcClientModelService
	clientModelTaskSvc interfaces.GrpcClientModelTaskService

	// internals
	mu      sync.Mutex
	cache   sync.Map
	drivers sync.Map
}

func (svc *Service) InsertData(id primitive.ObjectID, records ...interface{}) (err error) {
	//t, err := svc.getTask(id)
	//if err != nil {
	//	return err
	//}
	panic("implement me")
}

func (svc *Service) InsertLogs(id primitive.ObjectID, logs ...string) (err error) {
	log.Infof(logs[0])
	l, err := svc.getLogDriver(id)
	if err != nil {
		return err
	}
	return l.WriteLines(logs)
}

func (svc *Service) getTask(id primitive.ObjectID) (t interfaces.Task, err error) {
	// attempt to get from cache
	res, ok := svc.cache.Load(id)
	if ok {
		t, ok = res.(interfaces.Task)
		if ok {
			return t, nil
		}
	}

	// task
	t, err = svc.clientModelTaskSvc.GetTaskById(id)
	if err != nil {
		return nil, err
	}
	svc.cache.Store(id, t)

	return t, nil
}

func (svc *Service) getLogDriver(id primitive.ObjectID) (l clog.Driver, err error) {
	// attempt to get from cache
	res, ok := svc.drivers.Load(id)
	if ok {
		l, ok = res.(clog.Driver)
		if ok {
			return l, nil
		}
	}

	// TODO: other types of log drivers
	l, err = clog.NewSeaweedFsLogDriver(&clog.SeaweedFsLogDriverOptions{
		Prefix: id.Hex(),
	})
	if err != nil {
		return nil, err
	}
	svc.drivers.Store(id, l)

	return l, nil
}

func NewTaskStatsService(opts ...Option) (svc2 interfaces.TaskStatsService, err error) {
	// base service
	baseSvc, err := task.NewBaseService()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// service
	svc := &Service{
		TaskBaseService: baseSvc,
		cache:           sync.Map{},
		drivers:         sync.Map{},
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// node config service
	nodeCfgSvc, err := config.NewNodeConfigService()
	if err != nil {
		return nil, trace.TraceError(err)
	}
	svc.nodeCfgSvc = nodeCfgSvc

	// dependency injection
	c := dig.New()
	if err := c.Provide(client.ProvideServiceDelegate(svc.nodeCfgSvc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.ProvideTaskServiceDelegate(svc.nodeCfgSvc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(
		clientModelSvc interfaces.GrpcClientModelService,
		clientModelTaskSvc interfaces.GrpcClientModelTaskService,
	) {
		svc.clientModelSvc = clientModelSvc
		svc.clientModelTaskSvc = clientModelTaskSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}

var store = sync.Map{}

func GetTaskStatsService(path string, opts ...Option) (svr interfaces.TaskStatsService, err error) {
	if path == "" {
		path = config2.DefaultConfigPath
	}
	opts = append(opts, WithConfigPath(path))
	res, ok := store.Load(path)
	if ok {
		svr, ok = res.(interfaces.TaskStatsService)
		if ok {
			return svr, nil
		}
	}
	svr, err = NewTaskStatsService(opts...)
	if err != nil {
		return nil, err
	}
	store.Store(path, svr)
	return svr, nil
}

func ProvideGetTaskStatsService(path string, opts ...Option) func() (svr interfaces.TaskStatsService, err error) {
	return func() (svr interfaces.TaskStatsService, err error) {
		return GetTaskStatsService(path, opts...)
	}
}
