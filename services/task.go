package services

import (
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/models"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/redis"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

type TaskServiceInterface interface {
	Init() (err error)
	Close()
	Assign(t *models2.Task) (err error)
	Fetch() (t *models2.Task, err error)
	Run(taskId primitive.ObjectID) (err error)
	Cancel(taskId primitive.ObjectID) (err error)
	FindLogs(id primitive.ObjectID, pattern string, skip, size int) (lines []string, err error)
}

type TaskServiceOptions struct {
	IsMaster        bool // whether taskService is on master node
	MaxRunners      int  // max TaskRunner count that can run on taskService, default: 8
	PollWaitSeconds int  // number of seconds that taskService.Fetch will be executed, default: 5
}

func NewTaskService(options *TaskServiceOptions) (s *taskService, err error) {
	// normalize options
	if options == nil {
		options = &TaskServiceOptions{
			IsMaster: false,
		}
	}

	// normalize MaxRunners
	if options.MaxRunners == 0 {
		options.MaxRunners = 8
	}

	// normalize PollWaitSeconds
	if options.PollWaitSeconds == 0 {
		options.PollWaitSeconds = 5
	}

	// construct taskService
	s = &taskService{
		runnersCount: 0,
		runners:      sync.Map{},
		opts:         options,
	}

	return s, nil
}

func InitTaskService() (err error) {
	TaskService, err = NewTaskService(&TaskServiceOptions{
		IsMaster:        viper.GetBool("server.master"),
		MaxRunners:      viper.GetInt("task.maxRunners"),      // TODO: implement in db
		PollWaitSeconds: viper.GetInt("task.pollWaitSeconds"), // TODO: implement in db
	})
	if err != nil {
		return err
	}
	go TaskService.Init()
	return nil
}

func CloseTaskService() {
	TaskService.Close()
}

type taskService struct {
	runnersCount int                 // number of task runners
	runners      sync.Map            // pool of task runners started
	active       bool                // whether the task service is active
	opts         *TaskServiceOptions // options
}

func (s *taskService) Init() (err error) {
	// set taskService.active to true
	s.active = true

	for {
		// stop if taskService.active is set to false
		if !s.active {
			return constants.ErrStopped
		}

		// wait for a period
		time.Sleep(time.Duration(s.opts.PollWaitSeconds) * time.Second)

		// stop if taskService.active is set to false
		if !s.active {
			return constants.ErrStopped
		}

		// skip if exceeding max runners
		if s.runnersCount >= s.opts.MaxRunners {
			continue
		}

		// fetch task
		t, err := s.Fetch()
		if err != nil {
			if err != constants.ErrNoTasksAvailable {
				log.Error("fetch task error: " + err.Error())
			}
			continue
		}
		if t.Id.IsZero() {
			return constants.ErrEmptyValue
		}

		// run task (async)
		if err := s.Run(t.Id); err != nil {
			log.Error("run task error: " + err.Error())
		}
	}
}

func (s *taskService) Close() {
	s.active = false
}

func (s *taskService) Assign(t *models2.Task) (err error) {
	// validate options
	if !s.opts.IsMaster {
		return constants.ErrForbidden
	}

	// task message
	msg := entity.TaskMessage{
		Id: t.Id,
	}

	// serialization
	msgStr, err := msg.ToString()
	if err != nil {
		return err
	}

	// queue name
	var queue string
	if t.NodeId.IsZero() {
		queue = "tasks:public"
	} else {
		queue = "tasks:node:" + t.NodeId.Hex()
	}

	// enqueue
	if err := redis.RedisClient.RPush(queue, msgStr); err != nil {
		return err
	}

	// set task status as "pending" and save to database
	if err := s.saveTask(t, constants.StatusPending); err != nil {
		return err
	}

	return nil
}

func (s *taskService) Fetch() (t *models2.Task, err error) {
	// message
	var msg string

	// fetch task from node queue
	n, err := NodeService.GetCurrentNode()
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
	t, err = models.MustGetRootService().GetTaskById(tMsg.Id)
	if err != nil {
		return t, err
	}

	return t, nil
}

func (s *taskService) Run(taskId primitive.ObjectID) (err error) {
	_, ok := s.runners.Load(taskId)
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
	s.runners.Store(taskId, r)
	s.runnersCount++

	// create a goroutine to run task
	go func() {
		// run task process (blocking)
		// error or finish after task runner ends
		if err := r.Run(); err != nil {
			switch err {
			case constants.ErrTaskError:
				log.Error(fmt.Sprintf("task (_id=%s) finished with error: %s", r.tid, err.Error()))
			case constants.ErrTaskCancelled:
				log.Error(fmt.Sprintf("task (_id=%s) was cancelled", r.tid))
			default:
				log.Error(fmt.Sprintf("task (_id=%s) finished with unknown error: %s", r.tid, err.Error()))
			}
			return
		}
		log.Info(fmt.Sprintf("task (_id=%s) finished", r.tid))
	}()

	return nil
}

func (s *taskService) Cancel(taskId primitive.ObjectID) (err error) {
	r, err := s.getTaskRunner(taskId)
	if err != nil {
		return err
	}
	if err := r.Cancel(); err != nil {
		return err
	}
	return nil
}

func (s *taskService) FindLogs(taskId primitive.ObjectID, pattern string, skip, size int) (lines []string, err error) {
	r, err := s.getTaskRunner(taskId)
	if err != nil {
		return lines, err
	}
	lines, err = r.l.Find(pattern, skip, size)
	if err != nil {
		return lines, err
	}
	return lines, nil
}

func (s *taskService) getTaskRunner(taskId primitive.ObjectID) (r *TaskRunner, err error) {
	v, ok := s.runners.Load(taskId)
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

func (s *taskService) saveTask(t *models2.Task, status string) (err error) {
	// normalize status
	if status == "" {
		status = constants.StatusPending
	}

	// set task status
	t.Status = status

	// attempt to get task from database
	_, err = models.MustGetRootService().GetTaskById(t.Id)
	if err != nil {
		// if task does not exist, add to database
		if err == mongo2.ErrNoDocuments {
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

var TaskService *taskService
