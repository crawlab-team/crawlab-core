package services

import (
	"bufio"
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/crawlab-team/crawlab-core/services/sys_exec"
	"github.com/crawlab-team/crawlab-core/utils"
	clog "github.com/crawlab-team/crawlab-log"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"sync"
	"time"
)

type TaskRunnerInterface interface {
	Run() (err error)
	Cancel() (err error)
	Dispose() (err error)
}

type TaskRunnerOptions struct {
	Task          *model.Task // Task to run
	LogDriverType string      // log driver type
}

func NewTaskRunner(options *TaskRunnerOptions) (r *TaskRunner, err error) {
	// validate options
	if options.Task == nil {
		return r, constants.ErrInvalidOptions
	}

	// task
	t := options.Task

	// spider
	spider, err := model.GetSpider(t.SpiderId)
	if err != nil {
		return r, err
	}
	s := &spider

	// worker file system service using a temp directory
	fsPath := fmt.Sprintf("%s/%s", viper.GetString("spider.path"), s.Id.Hex())
	//cwd := fmt.Sprintf("%s/%s", viper.GetString("spider.workspace"), uuid.New().String())
	cwd := fmt.Sprintf("%s/%s", os.TempDir(), uuid.New().String())
	fs, err := NewFileSystemService(&FileSystemServiceOptions{
		IsMaster:      false,
		FsPath:        fsPath,
		WorkspacePath: cwd,
	})
	if err != nil {
		return r, err
	}

	// sync files to workspace
	if err := fs.SyncToWorkspace(); err != nil {
		return r, err
	}

	// task runner
	r = &TaskRunner{
		fs:   fs,
		t:    t,
		s:    s,
		ch:   make(chan constants.TaskSignal),
		cwd:  cwd,
		opts: options,
	}

	// log driver
	// TODO: configure TTL
	r.l, err = r.getLogDriver()
	if err != nil {
		return r, err
	}

	return r, nil
}

type TaskRunner struct {
	cmd  *exec.Cmd                 // process command instance
	pid  int                       // process id
	fs   *FileSystemService        // file system service
	l    clog.Driver               // log service
	t    *model.Task               // task
	s    *model.Spider             // spider
	ch   chan constants.TaskSignal // channel to communicate between TaskService and TaskRunner
	envs []model.Env               // environment variables
	opts *TaskRunnerOptions        // options
	cwd  string                    // working directory

	// log internals
	writeLock     *sync.Mutex
	scannerStdout *bufio.Scanner
	scannerStderr *bufio.Scanner
	readerStdout  *bufio.Reader
	readerStderr  *bufio.Reader
}

func (r *TaskRunner) Run() (err error) {
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

	// wait for signal
	signal := <-r.ch
	switch signal {
	case constants.TaskSignalFinish:
		err = nil
	case constants.TaskSignalCancel:
		err = constants.ErrTaskCancelled
	case constants.TaskSignalError:
		err = constants.ErrTaskError
	case constants.TaskSignalLost:
		err = constants.ErrTaskLost
	default:
		return constants.ErrInvalidSignal
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
	if r.cwd == "" {
		return constants.ErrUnableToDispose
	}
	if _, err := os.Stat(r.cwd); err != nil {
		return constants.ErrAlreadyDisposed
	}
	if err := os.RemoveAll(r.cwd); err != nil {
		return err
	}
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
			Prefix:  r.t.Id,
		}
	}
	return options
}

func (r *TaskRunner) configureLogging() (err error) {
	// set write lock
	r.writeLock = &sync.Mutex{}

	// buffer
	//buf := new([]byte)
	//buf := make([]byte, 1024)

	// set stdout reader
	stdout, err := r.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	r.scannerStdout = bufio.NewScanner(stdout)
	//r.scannerStdout.Buffer(buf, 4096)
	//r.cmd.Stdout = r.l

	// set stderr reader
	stderr, err := r.cmd.StderrPipe()
	if err != nil {
		return err
	}
	r.scannerStderr = bufio.NewScanner(stderr)
	//r.scannerStderr.Buffer(buf, 4096)
	//r.cmd.Stdout = r.l

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
	r.cmd.Env = append(os.Environ(), "CRAWLAB_TASK_ID="+r.t.Id)
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
	variables := model.GetVariableList()
	for _, variable := range variables {
		r.cmd.Env = append(r.cmd.Env, variable.Key+"="+variable.Value)
	}
	return nil
}

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
