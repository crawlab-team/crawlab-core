package stats

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/task"
	"github.com/crawlab-team/go-trace"
	"sync"
)

type Service struct {
	// dependencies
	interfaces.TaskBaseService
	nodeCfgSvc interfaces.NodeConfigService
}

func (svc *Service) InsertData(t interfaces.Task, records ...interface{}) (err error) {
	panic("implement me")
}

func (svc *Service) InsertLogs(t interfaces.Task, logs ...string) (err error) {
	panic("implement me")
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
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	return svc, nil
}

var store = sync.Map{}

func GetTaskStatsService(path string, opts ...Option) (svr interfaces.TaskStatsService, err error) {
	if path == "" {
		path = config.DefaultConfigPath
	}
	opts = append(opts, WithConfigPath(path))
	res, ok := store.Load(path)
	if !ok {
		return NewTaskStatsService(opts...)
	}
	svr, ok = res.(interfaces.TaskStatsService)
	if !ok {
		return NewTaskStatsService(opts...)
	}
	return svr, nil
}

func ProvideGetTaskStatsService(path string, opts ...Option) func() (svr interfaces.TaskStatsService, err error) {
	return func() (svr interfaces.TaskStatsService, err error) {
		return GetTaskStatsService(path, opts...)
	}
}
