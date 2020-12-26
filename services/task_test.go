package services

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/model"
	db "github.com/crawlab-team/crawlab-db"
	cfs "github.com/crawlab-team/crawlab-fs"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

type TaskTestObject struct {
	nodes    []*model.Node
	spiders  []*model.Spider
	tasks    []*model.Task
	fs       *FileSystemService
	fsPath   string
	repoPath string
}

func (to *TaskTestObject) GetFsPath(s model.Spider) (fsPath string) {
	return fmt.Sprintf("%s/%s", "/spiders", s.Id.Hex())
}

func (to *TaskTestObject) GetRepoPath(s model.Spider) (fsPath string) {
	return fmt.Sprintf("./tmp/repo/%s", s.Id.Hex())
}

func (to *TaskTestObject) CreateNode(name string, isMaster bool) (n *model.Node, err error) {
	n = &model.Node{
		Id:       bson.NewObjectId(),
		Name:     name,
		IsMaster: isMaster,
	}
	if err := n.Add(); err != nil {
		return n, err
	}
	to.nodes = append(to.nodes, n)
	return n, nil
}

func (to *TaskTestObject) CreateSpider(name string) (s *model.Spider, err error) {
	s = &model.Spider{
		Id:   bson.NewObjectId(),
		Name: name,
		Type: constants.Customized,
		Cmd:  fmt.Sprintf("python %s.py", name),
		Envs: []model.Env{
			{Name: "Env1", Value: "Value1"},
			{Name: "Env2", Value: "Value2"},
		},
		FileId:    bson.ObjectIdHex(constants.ObjectIdNull),
		ProjectId: bson.ObjectIdHex(constants.ObjectIdNull),
		UserId:    bson.ObjectIdHex(constants.ObjectIdNull),
	}
	if err := s.Add(); err != nil {
		return s, err
	}
	to.spiders = append(to.spiders, s)
	return s, nil
}

func (to *TaskTestObject) CreateTask(s *model.Spider) (t *model.Task, err error) {
	t = &model.Task{
		Id:         uuid.New().String(),
		SpiderId:   s.Id,
		Type:       constants.TaskTypeSpider,
		NodeId:     bson.ObjectIdHex(constants.ObjectIdNull),
		ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
		UserId:     bson.ObjectIdHex(constants.ObjectIdNull),
	}
	if err := model.AddTask(*t); err != nil {
		return t, err
	}
	to.tasks = append(to.tasks, t)
	return t, nil
}

func setupTask() (to *TaskTestObject, err error) {
	// test object
	to = &TaskTestObject{}

	// set debug
	viper.Set("debug", true)

	// mongo
	viper.Set("mongo.host", "localhost")
	viper.Set("mongo.port", "27017")
	viper.Set("mongo.db", "test")
	if err := db.InitMongo(); err != nil {
		return to, err
	}

	// redis
	viper.Set("redis.address", "localhost")
	viper.Set("redis.port", "6379")
	viper.Set("redis.database", "0")
	if err := db.InitRedis(); err != nil {
		return to, err
	}

	// set paths
	viper.Set("log.path", "/logs")
	viper.Set("spider.path", "/spiders")
	viper.Set("spider.workspace", "./tmp/workspace")

	// cleanup
	cleanupTask(to)

	// fs
	to.fs, err = NewFileSystemService(&FileSystemServiceOptions{
		IsMaster: true,
		FsPath:   to.fsPath,
		RepoPath: to.repoPath,
	})
	if err != nil {
		return to, err
	}

	// nodes
	nn := 2
	for i := 0; i < nn; i++ {
		name := fmt.Sprintf("node%d", i+1)
		isMaster := false
		if i == 0 {
			isMaster = true
		}
		n, err := to.CreateNode(name, isMaster)
		if err != nil {
			return to, err
		}
		to.nodes = append(to.nodes, n)
	}

	// spiders
	ns := 3
	for i := 0; i < ns; i++ {
		name := fmt.Sprintf("s%d", i+1)
		_, err := to.CreateSpider(name)
		if err != nil {
			return to, err
		}
	}

	// add scripts
	py1 := `print('it works')`
	if err := to.fs.Save("s1.py", []byte(py1)); err != nil {
		return to, err
	}
	py2 := `
import time
import sys
for i in range(3):
    print('line: ' + str(i))
    sys.stdout.flush()
`
	if err := to.fs.Save("s2.py", []byte(py2)); err != nil {
		return to, err
	}
	py3 := `print('it works')`
	if err := to.fs.Save("s3.py", []byte(py3)); err != nil {
		return to, err
	}

	return to, nil
}

func cleanupTask(to *TaskTestObject) {
	if m, err := cfs.NewSeaweedFSManager(); err == nil {
		_ = m.DeleteDir("/logs")
		_ = m.DeleteDir("/spiders")
	}
	for _, s := range to.spiders {
		_ = model.RemoveSpider(s.Id)
	}
	to.spiders = []*model.Spider{}
	for _, t := range to.tasks {
		_ = model.RemoveTask(t.Id)
	}
	to.tasks = []*model.Task{}
	_ = os.RemoveAll("./tmp/repo")
	_ = db.RedisClient.Del("tasks:public")
}

func TestNewTaskService(t *testing.T) {
	to, err := setupTask()
	require.Nil(t, err)

	// create master TaskService
	s, err := NewTaskService(&TaskServiceOptions{
		IsMaster: true,
	})
	require.Nil(t, err)
	require.Equal(t, 0, s.runnersCount)
	require.Equal(t, true, s.opts.IsMaster)
	require.Equal(t, 8, s.opts.MaxRunners)
	require.Equal(t, 5, s.opts.PollWaitSeconds)

	cleanupTask(to)
}

func TestTaskService_Init(t *testing.T) {
	to, err := setupTask()
	require.Nil(t, err)

	// create master TaskService
	s, err := NewTaskService(&TaskServiceOptions{
		IsMaster: true,
	})
	require.Nil(t, err)

	// test init
	go s.Init()

	// TODO: test

	cleanupTask(to)
}

func TestTaskService_Assign(t *testing.T) {
	to, err := setupTask()
	require.Nil(t, err)

	// create master TaskService
	s, err := NewTaskService(&TaskServiceOptions{
		IsMaster: true,
	})
	require.Nil(t, err)

	// test assign (without init)
	task, err := to.CreateTask(to.spiders[0])
	require.Nil(t, err)
	err = s.Assign(*task)
	require.Nil(t, err)
	count, err := db.RedisClient.LLen("tasks:public")
	require.Nil(t, err)
	require.Equal(t, 1, count)
	result, err := db.RedisClient.LPop("tasks:public")
	require.Nil(t, err)
	require.NotEmpty(t, result)
	count, err = db.RedisClient.LLen("tasks:public")
	require.Nil(t, err)
	require.Equal(t, 0, count)

	cleanupTask(to)
}

func TestTaskService_Fetch(t *testing.T) {
	to, err := setupTask()
	require.Nil(t, err)

	// create master TaskService
	s, err := NewTaskService(&TaskServiceOptions{
		IsMaster: true,
	})
	require.Nil(t, err)

	// test fetch (without init)
	task, err := to.CreateTask(to.spiders[0])
	require.Nil(t, err)
	err = s.Assign(*task)
	require.Nil(t, err)
	count, err := db.RedisClient.LLen("tasks:public")
	require.Nil(t, err)
	require.Equal(t, 1, count)
	task2, err := s.Fetch()
	require.Nil(t, err)
	require.Equal(t, task.Id, task2.Id)
	require.Nil(t, err)
	count, err = db.RedisClient.LLen("tasks:public")
	require.Nil(t, err)
	require.Equal(t, 0, count)

	cleanupTask(to)
}

func TestTaskService_Run(t *testing.T) {
	to, err := setupTask()
	require.Nil(t, err)

	// create master TaskService
	s, err := NewTaskService(&TaskServiceOptions{
		IsMaster: true,
	})
	require.Nil(t, err)

	// test run (full process: assign -> fetch -> run)
	task, err := to.CreateTask(to.spiders[1])
	require.Nil(t, err)
	err = s.Assign(*task)
	require.Nil(t, err)
	*task, err = s.Fetch()
	require.Nil(t, err)
	err = s.Run(task.Id)
	require.Nil(t, err)
	err = s.Run(task.Id)
	require.Equal(t, constants.ErrAlreadyExists, err)
	require.Equal(t, 1, s.runnersCount)
	*task, err = model.GetTask(task.Id)
	require.Nil(t, err)
	require.Equal(t, constants.StatusRunning, task.Status)
	time.Sleep(3 * time.Second)
	*task, err = model.GetTask(task.Id)
	require.Nil(t, err)
	require.Equal(t, constants.StatusFinished, task.Status)

	cleanupTask(to)
}
