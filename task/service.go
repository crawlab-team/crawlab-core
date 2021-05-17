package task

import (
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-db/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	nodeCfgSvc interfaces.NodeConfigService
	modelSvc   service.ModelService

	// settings variables
	maxRunners       int
	pollWaitDuration time.Duration

	// internals variables
	runnersCount int      // number of task runners
	runners      sync.Map // pool of task runners started
	active       bool     // whether the task service is active
}

func (svc *Service) Init() (err error) {
	// set Service.active to true
	svc.active = true

	for {
		// stop if Service.active is set to false
		if !svc.active {
			return constants.ErrStopped
		}

		// wait for a period
		time.Sleep(svc.pollWaitDuration)

		// stop if Service.active is set to false
		if !svc.active {
			return constants.ErrStopped
		}

		// skip if exceeding max runners
		if svc.runnersCount >= svc.maxRunners {
			continue
		}

		// fetch task
		t, err := svc.Fetch()
		if err != nil {
			if err != constants.ErrNoTasksAvailable {
				log.Error("fetch task error: " + err.Error())
			}
			continue
		}
		if t.GetId().IsZero() {
			return constants.ErrEmptyValue
		}

		// run task (async)
		if err := svc.Run(t.GetId()); err != nil {
			log.Error("run task error: " + err.Error())
		}
	}
}

func (svc *Service) Close() {
	svc.active = false
}

func (svc *Service) Assign(t interfaces.Task) (err error) {
	// TODO: implement task assign via grpc
	// validate options
	if !svc.nodeCfgSvc.IsMaster() {
		return constants.ErrForbidden
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

func (svc *Service) Fetch() (t interfaces.Task, err error) {
	// message
	var msg string

	// fetch task from node queue
	// TODO: implement priority queue
	n, err := svc.modelSvc.GetNodeByKey(svc.nodeCfgSvc.GetNodeKey(), nil)
	if err != nil {
		return t, err
	}
	if n != nil {
		queueCur := "tasks:node:" + n.Id.Hex()
		msg, err = redis.RedisClient.LPop(queueCur)
	}

	// fetch task from public queue if first fetch is not successful
	if msg == "" {
		err = nil
		queuePub := "tasks:public"
		msg, err = redis.RedisClient.LPop(queuePub)
		if err != nil {
			return t, err
		}
	}

	// no task fetched
	if msg == "" {
		return t, constants.ErrNoTasksAvailable
	}

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

func (svc *Service) Run(taskId primitive.ObjectID) (err error) {
	_, ok := svc.runners.Load(taskId)
	if ok {
		return constants.ErrAlreadyExists
	}

	// create a new task runner
	r, err := NewTaskRunner(&TaskRunnerOptions{
		TaskId: taskId,
	})
	if err != nil {
		return err
	}

	// save runner to pool
	svc.runners.Store(taskId, r)
	svc.runnersCount++

	// create a goroutine to run task
	go func() {
		// run task process (blocking)
		// error or finish after task runner ends
		if err := r.Run(); err != nil {
			switch err {
			case constants.ErrTaskError:
				log.Error(fmt.Sprintf("task (_id=%svc) finished with error: %svc", r.tid, err.Error()))
			case constants.ErrTaskCancelled:
				log.Error(fmt.Sprintf("task (_id=%svc) was cancelled", r.tid))
			default:
				log.Error(fmt.Sprintf("task (_id=%svc) finished with unknown error: %svc", r.tid, err.Error()))
			}
			return
		}
		log.Info(fmt.Sprintf("task (_id=%svc) finished", r.tid))
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

func (svc *Service) FindLogs(taskId primitive.ObjectID, pattern string, skip, size int) (lines []string, err error) {
	r, err := svc.getTaskRunner(taskId)
	if err != nil {
		return lines, err
	}
	lines, err = r.l.Find(pattern, skip, size)
	if err != nil {
		return lines, err
	}
	return lines, nil
}

func (svc *Service) SetMaxRunners(maxRunners int) {
	svc.maxRunners = maxRunners
}

func (svc *Service) SetPollWaitDuration(duration time.Duration) {
	svc.pollWaitDuration = duration
}

func (svc *Service) getTaskRunner(taskId primitive.ObjectID) (r *TaskRunner, err error) {
	v, ok := svc.runners.Load(taskId)
	if !ok {
		return r, constants.ErrNotExists
	}
	switch v.(type) {
	case *TaskRunner:
		r = v.(*TaskRunner)
	default:
		return r, constants.ErrInvalidType
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
	_, err = svc.modelSvc.GetTaskById(t.GetId())
	if err != nil {
		// if task does not exist, add to database
		if err == mongo.ErrNoDocuments {
			if err := t.Add(); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	} else {
		// otherwise, update
		if err := t.Save(); err != nil {
			return err
		}
		return nil
	}
}

func NewTaskService(opts ...Option) (svc2 interfaces.TaskService, err error) {
	// construct Service
	svc := &Service{
		maxRunners:       8,
		pollWaitDuration: 5,
		runnersCount:     0,
		runners:          sync.Map{},
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	return svc, nil
}
