package services

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/model"
	db "github.com/crawlab-team/crawlab-db"
	cfs "github.com/crawlab-team/crawlab-fs"
	clog "github.com/crawlab-team/crawlab-log"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

var spiderId = bson.NewObjectId()
var spider model.Spider
var taskId = uuid.New().String()
var task model.Task
var fs *FileSystemService
var fsPath = fmt.Sprintf("%s/%s", "/spiders", spiderId.Hex())
var repoPath = fmt.Sprintf("./tmp/repo/%s", spiderId.Hex())

func setupTaskRunner() (err error) {
	// set debug
	viper.Set("debug", true)

	// mongo
	viper.Set("mongo.host", "localhost")
	viper.Set("mongo.port", "27017")
	viper.Set("mongo.db", "test")
	if err := db.InitMongo(); err != nil {
		return err
	}

	// set paths
	viper.Set("log.path", "/logs")
	viper.Set("spider.path", "/spiders")
	viper.Set("spider.workspace", "./tmp/workspace")

	// cleanup
	cleanupTaskRunner()

	// fs
	fs, err = NewFileSystemService(&FileSystemServiceOptions{
		IsMaster: true,
		FsPath:   fsPath,
		RepoPath: repoPath,
	})
	if err != nil {
		return err
	}

	// spider
	spider = model.Spider{
		Id:   spiderId,
		Name: "test_spider",
		Type: constants.Customized,
		Cmd:  "python main.py",
		Envs: []model.Env{
			{Name: "Env1", Value: "Value1"},
			{Name: "Env2", Value: "Value2"},
		},
		FileId:    bson.ObjectIdHex(constants.ObjectIdNull),
		ProjectId: bson.ObjectIdHex(constants.ObjectIdNull),
		UserId:    bson.ObjectIdHex(constants.ObjectIdNull),
	}
	if err := spider.Add(); err != nil {
		return err
	}

	// task
	task = model.Task{
		Id:       taskId,
		SpiderId: spiderId,
		Type:     constants.TaskTypeSpider,
	}

	// add python script
	pythonScript := `print('it works')`
	if err := fs.Save("main.py", []byte(pythonScript)); err != nil {
		return err
	}

	// commit
	if err := fs.Commit("initial commit"); err != nil {
		return err
	}

	return nil
}

func cleanupTaskRunner() {
	if m, err := cfs.NewSeaweedFSManager(); err == nil {
		_ = m.DeleteDir("/logs")
		_ = m.DeleteDir("/spiders")
	}
	_ = model.RemoveSpider(spiderId)
	_ = model.RemoveTask(taskId)
	_ = os.RemoveAll("./tmp/repo")
}

func TestNewTaskRunner(t *testing.T) {
	err := setupTaskRunner()
	require.Nil(t, err)

	// create task runner
	runner, err := NewTaskRunner(&TaskRunnerOptions{
		Task:          &task,
		LogDriverType: clog.DriverTypeFs,
	})
	require.Nil(t, err)
	require.NotNil(t, runner.fs)
	require.NotNil(t, runner.t)
	require.NotNil(t, runner.s)
	require.NotNil(t, runner.l)
	require.NotNil(t, runner.ch)

	cleanupTaskRunner()
}

func TestTaskRunner_Run(t *testing.T) {
	err := setupTaskRunner()
	require.Nil(t, err)

	// create task runner
	runner, err := NewTaskRunner(&TaskRunnerOptions{
		Task:          &task,
		LogDriverType: clog.DriverTypeFs,
	})
	require.Nil(t, err)

	// run
	err = runner.Run()
	require.Nil(t, err)

	// test logs
	lines, err := runner.l.Find("", 0, 100)
	require.Nil(t, err)
	require.Len(t, lines, 1)
	require.Equal(t, "it works", lines[0])

	cleanupTaskRunner()
}

func TestTaskRunner_RunWithError(t *testing.T) {
	err := setupTaskRunner()
	require.Nil(t, err)

	// add error python script
	pythonScript := `
raise Exception('an error')
`
	err = fs.Save("main.py", []byte(pythonScript))
	require.Nil(t, err)

	// create task runner
	runner, err := NewTaskRunner(&TaskRunnerOptions{
		Task:          &task,
		LogDriverType: clog.DriverTypeFs,
	})
	require.Nil(t, err)

	// run
	err = runner.Run()
	require.Equal(t, constants.ErrTaskError, err)

	// test logs
	lines, err := runner.l.Find("", 0, 100)
	require.Nil(t, err)
	require.Greater(t, len(lines), 0)
	hasExceptionLog := false
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "exception") {
			hasExceptionLog = true
			break
		}
	}
	require.True(t, hasExceptionLog)

	cleanupTaskRunner()
}

func TestTaskRunner_RunLong(t *testing.T) {
	err := setupTaskRunner()
	require.Nil(t, err)

	// add a long task python script
	n := 5
	pythonScript := fmt.Sprintf(`
import time
for i in range(%d):
   print('line: ' + str(i))
   time.sleep(1)
`, n)
	err = fs.Save("main.py", []byte(pythonScript))
	require.Nil(t, err)

	// create task runner
	runner, err := NewTaskRunner(&TaskRunnerOptions{
		Task:          &task,
		LogDriverType: clog.DriverTypeFs,
	})
	require.Nil(t, err)

	// run
	err = runner.Run()
	require.Nil(t, err)

	// test logs
	lines, err := runner.l.Find("", 0, 100)
	require.Nil(t, err)
	require.Equal(t, n, len(lines))
	for i, line := range lines {
		require.Equal(t, fmt.Sprintf("line: %d", i), line)
	}

	cleanupTaskRunner()
}

func TestTaskRunner_Cancel(t *testing.T) {
	err := setupTaskRunner()
	require.Nil(t, err)

	// add a long task python script
	n := 10
	pythonScript := fmt.Sprintf(`
import time
for i in range(%d):
   print('line: ' + str(i))
   time.sleep(1)
`, n)
	err = fs.Save("main.py", []byte(pythonScript))
	require.Nil(t, err)

	// create task runner
	runner, err := NewTaskRunner(&TaskRunnerOptions{
		Task:          &task,
		LogDriverType: clog.DriverTypeFs,
	})
	require.Nil(t, err)

	// cancel
	go func() {
		time.Sleep(5 * time.Second)
		err = runner.Cancel()
		require.Nil(t, err)
	}()

	// run
	err = runner.Run()
	require.Equal(t, constants.ErrTaskCancelled, err)

	// test logs
	lines, err := runner.l.Find("", 0, 100)
	require.Nil(t, err)
	require.Greater(t, len(lines), 0)

	cleanupTaskRunner()
}
