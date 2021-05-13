package task

import (
	"bufio"
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/fs"
	"github.com/crawlab-team/crawlab-core/models"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/services"
	"github.com/crawlab-team/crawlab-core/services/sys_exec"
	"github.com/crawlab-team/crawlab-core/utils"
	clog "github.com/crawlab-team/crawlab-log"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"os/exec"
	"path"
	"time"
)

type TaskRunnerInterface interface {
	Run() (err error)
	Cancel() (err error)
	Dispose() (err error)
}

type TaskRunnerOptions struct {
	TaskService   *taskService       // TaskService on which the TaskRunner is executed
	TaskId        primitive.ObjectID // id of task (model.Task) to run
	LogDriverType string             // log driver type
}

func NewTaskRunner(options *TaskRunnerOptions) (r *TaskRunner, err error) {
	// validate options
	if options == nil {
		return r, constants.ErrInvalidOptions
	}
	if options.TaskId.IsZero() {
		return r, constants.ErrInvalidOptions
	}

	// normalize LogDriverType
	if options.LogDriverType == "" {
		options.LogDriverType = clog.DriverTypeFs
	}

	// task service
	svc := options.TaskService

	// task runner
	r = &TaskRunner{
		svc:  svc,
		tid:  options.TaskId,
		ch:   make(chan constants.TaskSignal),
		opts: options,
	}

	// initialize task runner
	if err := r.Init(); err != nil {
		return r, err
	}

	return r, nil
}

type TaskRunner struct {
	cmd  *exec.Cmd                   // process command instance
	pid  int                         // process id
	svc  *taskService                // taskService
	fs   *services.fileSystemService // file system service fileSystemService
	l    clog.Driver                 // log service log.Driver
	tid  primitive.ObjectID          // id of t (model.Task)
	t    *models2.Task               // task model.Task
	s    *models2.Spider             // spider model.Spider
	ch   chan constants.TaskSignal   // channel to communicate between taskService and TaskRunner
	envs []models2.Env               // environment variables
	opts *TaskRunnerOptions          // options
	cwd  string                      // working directory

	// log internals
	scannerStdout *bufio.Scanner
	scannerStderr *bufio.Scanner
}

func (r *TaskRunner) Init() (err error) {
	// update task
	if err := r.updateTask(""); err != nil {
		return err
	}

	// spider
	s, err := models.MustGetRootService().GetSpiderById(r.t.SpiderId)
	if err != nil {
		return err
	}
	r.s = s

	// worker file system service using a temp directory
	fsPath := fmt.Sprintf("%s/%s", viper.GetString("spider.path"), r.s.Id.Hex())
	//cwd := fmt.Sprintf("%s/%s", viper.GetString("spider.workspace"), uuid.New().String())
	r.cwd = path.Join(os.TempDir(), uuid.New().String())
	r.fs, err = fs.NewFsService(&fs.FileSystemServiceOptions{
		IsMaster:      false,
		FsPath:        fsPath,
		WorkspacePath: r.cwd,
	})
	if err != nil {
		return err
	}

	// sync files to workspace
	if err := r.fs.SyncToWorkspace(); err != nil {
		return err
	}

	// log driver
	// TODO: configure TTL
	r.l, err = r.getLogDriver()
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRunner) Run() (err error) {
	// update task status (processing)
	if err := r.updateTask(constants.StatusRunning); err != nil {
		return err
	}

	// configure cmd
	if err := r.configureCmd(); err != nil {
		return err
	}

	// configure environment variables
	if err := r.configureEnv(); err != nil {
		return err
	}

	// configure logging
	if err := r.configureLogging(); err != nil {
		return err
	}

	// start process
	if err := r.cmd.Start(); err != nil {
		return err
	}

	// start logging
	go r.startLogging()

	// process id
	if r.cmd.Process == nil {
		return constants.ErrNotExists
	}
	r.pid = r.cmd.Process.Pid

	// wait for process to finish
	go r.wait()

	// start health check
	go r.startHealthCheck()

	// declare task status
	status := ""

	// wait for signal
	signal := <-r.ch
	switch signal {
	case constants.TaskSignalFinish:
		err = nil
		status = constants.StatusFinished
	case constants.TaskSignalCancel:
		err = constants.ErrTaskCancelled
		status = constants.StatusCancelled
	case constants.TaskSignalError:
		err = constants.ErrTaskError
		status = constants.StatusError
	case constants.TaskSignalLost:
		err = constants.ErrTaskLost
		status = constants.StatusError
	default:
		err = constants.ErrInvalidSignal
		status = constants.StatusError
	}

	// validate task status
	if status == "" {
		return constants.ErrInvalidType
	}

	// update task status
	if err := r.updateTask(status); err != nil {
		return err
	}

	// flush log
	_ = r.l.Flush()

	// dispose
	_ = r.Dispose()

	return err
}

func (r *TaskRunner) Cancel() (err error) {
	// kill process
	if err := sys_exec.KillProcess(r.cmd); err != nil {
		return err
	}

	// make sure the process does not exist
	cancelWaitSeconds := viper.GetInt("task.cancelWaitSeconds")
	if cancelWaitSeconds == 0 {
		cancelWaitSeconds = 30
	}
	for i := 0; i < cancelWaitSeconds; i++ {
		exists, _ := process.PidExists(int32(r.pid))
		if !exists {
			// successfully cancelled
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	// unable to cancel
	return constants.ErrUnableToCancel
}

func (r *TaskRunner) Dispose() (err error) {
	// validate whether it is disposable
	if r.cwd == "" {
		return constants.ErrUnableToDispose
	}
	if _, err := os.Stat(r.cwd); err != nil {
		return constants.ErrAlreadyDisposed
	}

	// remove working directory
	if err := os.RemoveAll(r.cwd); err != nil {
		return err
	}

	// remove in taskService
	if r.svc != nil {
		if err := r.removeTaskRunner(r.t.Id); err != nil {
			return err
		}
	}

	// remove
	return nil
}

func (r *TaskRunner) configureCmd() (err error) {
	var cmdStr string
	if r.t.Type == constants.TaskTypeSpider {
		// spider task
		if r.s.Type == constants.Configurable {
			// configurable spider
			cmdStr = "scrapy crawl config_spider"
		} else {
			// customized spider
			cmdStr = r.s.Cmd
		}
	} else if r.t.Type == constants.TaskTypeSystem {
		// system task
		cmdStr = r.t.Cmd
	}

	// parameters
	if r.t.Param != "" {
		cmdStr += " " + r.t.Param
	}

	// get cmd instance
	r.cmd = sys_exec.BuildCmd(cmdStr)

	// set working directory
	r.cmd.Dir = r.cwd

	// configure pgid to allow killing sub processes
	sys_exec.Setpgid(r.cmd)

	return nil
}

func (r *TaskRunner) getLogDriver() (driver clog.Driver, err error) {
	options := r.getLogDriverOptions()
	driver, err = clog.NewLogDriver(r.opts.LogDriverType, options)
	if err != nil {
		return driver, err
	}
	return driver, nil
}

func (r *TaskRunner) getLogDriverOptions() (options interface{}) {
	switch r.opts.LogDriverType {
	case clog.DriverTypeFs:
		options = &clog.SeaweedFSLogDriverOptions{
			BaseDir: viper.GetString("log.path"),
			Prefix:  r.t.Id.Hex(),
		}
	}
	return options
}

func (r *TaskRunner) configureLogging() (err error) {
	// set stdout reader
	stdout, _ := r.cmd.StdoutPipe()
	r.scannerStdout = bufio.NewScanner(stdout)

	// set stderr reader
	stderr, _ := r.cmd.StderrPipe()
	r.scannerStderr = bufio.NewScanner(stderr)

	return nil
}

func (r *TaskRunner) startLogging() {
	// start reading stdout
	go r.startLoggingReaderStdout()

	// start reading stderr
	go r.startLoggingReaderStderr()
}

func (r *TaskRunner) startLoggingReaderStdout() {
	utils.LogDebug("begin startLoggingReaderStdout")
	for r.scannerStdout.Scan() {
		line := r.scannerStdout.Text()
		utils.LogDebug(fmt.Sprintf("scannerStdout line: %s", line))
		_ = r.l.WriteLine(line)
	}
	// reach end
	utils.LogDebug("scannerStdout reached end")
}

func (r *TaskRunner) startLoggingReaderStderr() {
	utils.LogDebug("begin startLoggingReaderStderr")
	for r.scannerStderr.Scan() {
		line := r.scannerStderr.Text()
		utils.LogDebug(fmt.Sprintf("scannerStderr line: %s", line))
		_ = r.l.WriteLine(line)
	}
	// reach end
	utils.LogDebug("scannerStderr reached end")
}

func (r *TaskRunner) startHealthCheck() {
	for {
		exists, _ := process.PidExists(int32(r.pid))
		if !exists {
			// process lost
			r.ch <- constants.TaskSignalLost
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (r *TaskRunner) configureEnv() (err error) {
	// TODO: refactor
	//envs := r.s.Envs
	//if r.s.Type == constants.Configurable {
	//	// 数据库配置
	//	envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_HOST", Value: viper.GetString("mongo.host")})
	//	envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_PORT", Value: viper.GetString("mongo.port")})
	//	envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_DB", Value: viper.GetString("mongo.db")})
	//	envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_USERNAME", Value: viper.GetString("mongo.username")})
	//	envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_PASSWORD", Value: viper.GetString("mongo.password")})
	//	envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_AUTHSOURCE", Value: viper.GetString("mongo.authSource")})
	//
	//	// 设置配置
	//	for envName, envValue := range r.s.Config.Settings {
	//		envs = append(envs, model.Env{Name: "CRAWLAB_SETTING_" + envName, Value: envValue})
	//	}
	//}

	// 默认把Node.js的全局node_modules加入环境变量
	//envPath := os.Getenv("PATH")
	//nodePath := "/usr/lib/node_modules"
	//if !strings.Contains(envPath, nodePath) {
	//	_ = os.Setenv("PATH", nodePath+":"+envPath)
	//}
	//_ = os.Setenv("NODE_PATH", nodePath)

	// default results collection
	//col := utils.GetSpiderCol(r.s.Col, r.s.Name)

	// 默认环境变量
	r.cmd.Env = append(os.Environ(), "CRAWLAB_TASK_ID="+r.t.Id.Hex())
	//r.cmd.Env = append(r.cmd.Env, "CRAWLAB_COLLECTION="+col)
	//r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_HOST="+viper.GetString("mongo.host"))
	//r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_PORT="+viper.GetString("mongo.port"))
	//if viper.GetString("mongo.db") != "" {
	//	r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_DB="+viper.GetString("mongo.db"))
	//}
	//if viper.GetString("mongo.username") != "" {
	//	r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_USERNAME="+viper.GetString("mongo.username"))
	//}
	//if viper.GetString("mongo.password") != "" {
	//	r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_PASSWORD="+viper.GetString("mongo.password"))
	//}
	//if viper.GetString("mongo.authSource") != "" {
	//	r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_AUTHSOURCE="+viper.GetString("mongo.authSource"))
	//}
	//r.cmd.Env = append(r.cmd.Env, "PYTHONUNBUFFERED=0")
	//r.cmd.Env = append(r.cmd.Env, "PYTHONIOENCODING=utf-8")
	//r.cmd.Env = append(r.cmd.Env, "TZ=Asia/Shanghai")
	//r.cmd.Env = append(r.cmd.Env, "CRAWLAB_DEDUP_FIELD="+r.s.DedupField)
	//r.cmd.Env = append(r.cmd.Env, "CRAWLAB_DEDUP_METHOD="+r.s.DedupMethod)
	//if r.s.IsDedup {
	//	r.cmd.Env = append(r.cmd.Env, "CRAWLAB_IS_DEDUP=1")
	//} else {
	//	r.cmd.Env = append(r.cmd.Env, "CRAWLAB_IS_DEDUP=0")
	//}

	// task environment variables
	for _, env := range r.s.Envs {
		r.cmd.Env = append(r.cmd.Env, env.Name+"="+env.Value)
	}

	// global environment variables
	variables, err := models.MustGetRootService().GetVariableList(nil, nil)
	if err != nil {
		return err
	}
	for _, variable := range variables {
		r.cmd.Env = append(r.cmd.Env, variable.Key+"="+variable.Value)
	}
	return nil
}

// wait for process to finish and send task signal (constants.TaskSignal)
// to task runner's channel (TaskRunner.ch) according to exit code
func (r *TaskRunner) wait() {
	// wait for process to finish
	if err := r.cmd.Wait(); err != nil {
		exitError, ok := err.(*exec.ExitError)
		if !ok {
			r.ch <- constants.TaskSignalError
			return
		}
		exitCode := exitError.ExitCode()
		if exitCode == -1 {
			// cancel error
			r.ch <- constants.TaskSignalCancel
			return
		}

		// standard error
		r.ch <- constants.TaskSignalError
		return
	}

	// success
	r.ch <- constants.TaskSignalFinish
}

// update and get updated info of task (TaskRunner.t)
func (r *TaskRunner) updateTask(status string) (err error) {
	// update task status
	if r.t != nil && status != "" {
		r.t.Status = status
		if err := r.t.Save(); err != nil {
			return err
		}
	}

	// get task
	t, err := models.MustGetRootService().GetTaskById(r.tid)
	if err != nil {
		return err
	}

	// set task
	r.t = t

	return nil
}

func (r *TaskRunner) removeTaskRunner(taskId primitive.ObjectID) (err error) {
	_, ok := r.svc.runners.Load(taskId)
	if !ok {
		return constants.ErrNotExists
	}
	r.svc.runners.Delete(taskId)
	return nil
}
