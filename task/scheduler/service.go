package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/task"
	db "github.com/crawlab-team/crawlab-db"
	"github.com/crawlab-team/crawlab-db/redis"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	interfaces.TaskBaseService
	nodeCfgSvc interfaces.NodeConfigService
	modelSvc   service.ModelService
	redis      db.RedisClient

	// settings
	fetchInterval time.Duration

	// internals
	ch chan []interfaces.Task
}

func (svc *Service) Start() {
	//go svc.Fetch()
	//go svc.Assign()
	svc.Wait()
	svc.Stop()
}

func (svc *Service) Fetch(nodeKey string) (t interfaces.Task, err error) {
	return svc.fetch(nodeKey)
	//for {
	//	// return if quit is true
	//	if svc.IsStopped() {
	//		return
	//	}
	//
	//	// fetch task with retry
	//	if err := backoff.RetryNotify(func() error {
	//		// fetch
	//		tasks, err := svc.fetch()
	//		if err != nil {
	//			return trace.TraceError(err)
	//		}
	//
	//		// skip if no task fetched
	//		if tasks == nil {
	//			return nil
	//		}
	//
	//		// notify tasks channel
	//		svc.ch <- tasks
	//
	//		return nil
	//	}, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("task scheduler fetch task")); err != nil {
	//		trace.PrintError(err)
	//	}
	//
	//	// wait
	//	time.Sleep(svc.fetchInterval)
	//}
}

func (svc *Service) Assign() {
	for {
		// return if quit is true
		if svc.IsStopped() {
			return
		}

		// receive task from channel
		tasks := <-svc.ch

		// assign task
		if err := svc.assign(tasks); err != nil {
			trace.PrintError(err)
		}
	}
}

func (svc *Service) SetFetchInterval(interval time.Duration) {
	svc.fetchInterval = interval
}

func (svc *Service) GetTaskChannel() (ch chan []interfaces.Task) {
	return svc.ch
}

func (svc *Service) fetch(nodeKey string) (t interfaces.Task, err error) {
	// node
	n, err := svc.modelSvc.GetNodeByKey(nodeKey, nil)
	if err != nil {
		return nil, err
	}

	// validate node
	if !n.Enabled {
		return nil, errors.ErrorTaskForbidden
	}

	// attempt to fetch task from the queue for dedicated node
	t, err = svc._fetch(n.GetId())
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	// fetch task from the public queue (random)
	t, err = svc._fetch(primitive.NilObjectID)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	// no task fetched
	return nil, nil
}

func (svc *Service) _fetch(nodeId primitive.ObjectID) (t interfaces.Task, err error) {
	// dequeue record with max score
	value, err := svc.redis.ZPopMaxOne(svc.GetQueue(nodeId))
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// skip if empty result returned
	if value == "" {
		return nil, nil
	}

	// deserialize task message
	data := []byte(value)
	var tMsg entity.TaskMessage
	if err := json.Unmarshal(data, &tMsg); err != nil {
		return nil, trace.TraceError(err)
	}

	// task in db
	t, err = svc.modelSvc.GetTaskById(tMsg.Id)
	if err != nil {
		return nil, trace.TraceError(err)
	}

	return t, nil
}

func (svc *Service) assign(tasks []interfaces.Task) (err error) {
	wg := sync.WaitGroup{}
	wg.Add(len(tasks))
	for _, t := range tasks {
		go func(t interfaces.Task) {
			if err := svc._assign(t); err != nil {
				trace.PrintError(err)
				t.SetStatus(constants.TaskStatusError)
				t.SetError(err.Error())
				if err := delegate.NewModelDelegate(t).Save(); err != nil {
					trace.PrintError(err)
				}
			}
			wg.Done()
		}(t)
	}
	wg.Wait()
	return nil
}

func (svc *Service) _assign(t interfaces.Task) (err error) {
	// node
	n, err := svc.modelSvc.GetNodeById(t.GetNodeId())
	if err != nil {
		return err
	}

	// save task
	if err := svc.SaveTask(t, constants.TaskStatusAssigned); err != nil {
		return err
	}

	// log
	log.Infof("task scheduler assigned task[%s] successfully", t.GetId())

	return nil
}

func (svc *Service) getZScanPattern(nodeId primitive.ObjectID) (pattern string) {
	if nodeId.IsZero() {
		// public
		pattern = fmt.Sprintf("%s%s%s", constants.TaskKeyAnchor, constants.TaskListQueuePrefixPublic, constants.TaskKeyAnchor)
	} else {
		pattern = fmt.Sprintf("%s%s:%s%s", constants.TaskKeyAnchor, constants.TaskListQueuePrefixNodes, nodeId.Hex(), constants.TaskKeyAnchor)
	}
	return "*" + pattern + "*"
}

func (svc *Service) getAvailableNodes() (nodes []models.Node, err error) {
	query := bson.M{
		"enabled":   true,
		"active":    true,
		"available": true,
	}
	return svc.modelSvc.GetNodeList(query, nil)
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
		fetchInterval:   15 * time.Second,
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
	if err := c.Provide(redis.GetRedisClient); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(
		nodeCfgSvc interfaces.NodeConfigService,
		modelSvc service.ModelService,
		redis db.RedisClient,
	) {
		svc.nodeCfgSvc = nodeCfgSvc
		svc.modelSvc = modelSvc
		svc.svr = svr
		svc.redis = redis
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// task channel
	svc.ch = make(chan []interfaces.Task)

	return svc, nil
}

func ProvideTaskSchedulerService(path string, opts ...Option) func() (svc interfaces.TaskSchedulerService, err error) {
	opts = append(opts, WithConfigPath(path))
	return func() (svc interfaces.TaskSchedulerService, err error) {
		return NewTaskSchedulerService(opts...)
	}
}
