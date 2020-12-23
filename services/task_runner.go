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
	"path/filepath"
	"strings"
	"time"
)

type TaskRunnerInterface interface {
	Run() (err error)
	Cancel() (err error)
	Dispose() (err error)
}

type TaskRunnerOptions struct {
	Task          *model.Task               // Task to run
	Channel       chan constants.TaskSignal // Channel to communicate between TaskService and TaskRunner
	LogDriverType string                    // log driver type
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
	workspacePath := fmt.Sprintf("%s/%s", os.TempDir(), uuid.New().String())
	fs, err := NewFileSystemService(&FileSystemServiceOptions{
		IsMaster:      false,
		FsPath:        fsPath,
		WorkspacePath: workspacePath,
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
	cmd    *exec.Cmd                 // process command instance
	pid    int                       // process id
	killed bool                      // whether process is killed
	fs     *FileSystemService        // file system service
	l      clog.Driver               // log service
	t      *model.Task               // task
	s      *model.Spider             // spider
	ch     chan constants.TaskSignal // channel to communicate between TaskService and TaskRunner
	envs   []model.Env               // environment variables
	opts   *TaskRunnerOptions        // options

	// log internals
	readerStdout     *bufio.Reader
	readerStderr     *bufio.Reader
	bufferLines      []string
	isStdoutFinished bool
	isStderrFinished bool
}

func (r *TaskRunner) Run() (err error) {
	// get cmd
	r.cmd, err = r.getCmd()
	if err != nil {
		return err
	}

	// working directory
	r.cmd.Dir = filepath.Join(
		viper.GetString("spider.path"),
		r.s.Name,
	)

	// TODO: configure environment variables

	// configure pgid to allow killing sub processes
	sys_exec.Setpgid(r.cmd)

	// configure task logging
	if err := r.setLogConfig(); err != nil {
		return err
	}

	// start logging
	go r.startLogging()

	// start process
	if err := r.cmd.Start(); err != nil {
		return err
	}

	// process id
	if r.cmd.Process == nil {
		return constants.ErrNotExists
	}
	r.pid = r.cmd.Process.Pid

	// wait for process to finish
	if err := r.cmd.Wait(); err != nil {
		// if killed is flagged as true, return ErrTaskCancelled
		if r.killed {
			return constants.ErrTaskCancelled
		}

		// standard error
		return err
	}

	return nil
}

func (r *TaskRunner) Cancel() (err error) {
	// flag killed as true
	r.killed = true

	// kill process
	if err := sys_exec.KillProcess(r.cmd); err != nil {
		return err
	}

	// make sure the process does not exist
	cancelWaitSeconds := viper.GetInt("task.cancelWaitSeconds")
	for i := 0; i < cancelWaitSeconds; i++ {
		exists, _ := process.PidExists(int32(r.pid))
		if !exists {
			return nil
		}
	}

	return constants.ErrUnableToCancelTask
}

func (r *TaskRunner) Dispose() (err error) {
	// TODO: implement
	return nil
}

func (r *TaskRunner) getCmd() (cmd *exec.Cmd, err error) {
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
	cmd = sys_exec.BuildCmd(cmdStr)

	return cmd, nil
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

func (r *TaskRunner) setLogConfig() (err error) {
	// get stdout reader
	stdout, err := r.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	r.readerStdout = bufio.NewReader(stdout)

	// get stderr reader
	stderr, err := r.cmd.StderrPipe()
	if err != nil {
		return err
	}
	r.readerStderr = bufio.NewReader(stderr)

	return nil
}

func (r *TaskRunner) startLogging() {
	for {
		// write log lines
		_ = r.l.WriteLines(r.bufferLines)

		// reset buffer log lines
		r.bufferLines = []string{}

		// end if finished flags are set to true
		if r.isStdoutFinished && r.isStderrFinished {
			return
		}

		// wait for a period
		time.Sleep(5 * time.Second)
	}
}

func (r *TaskRunner) setEnv() *exec.Cmd {
	// TODO: refactor
	envs := r.s.Envs
	if r.s.Type == constants.Configurable {
		// 数据库配置
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_HOST", Value: viper.GetString("mongo.host")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_PORT", Value: viper.GetString("mongo.port")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_DB", Value: viper.GetString("mongo.db")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_USERNAME", Value: viper.GetString("mongo.username")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_PASSWORD", Value: viper.GetString("mongo.password")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_AUTHSOURCE", Value: viper.GetString("mongo.authSource")})

		// 设置配置
		for envName, envValue := range r.s.Config.Settings {
			envs = append(envs, model.Env{Name: "CRAWLAB_SETTING_" + envName, Value: envValue})
		}
	}

	// 默认把Node.js的全局node_modules加入环境变量
	envPath := os.Getenv("PATH")
	nodePath := "/usr/lib/node_modules"
	if !strings.Contains(envPath, nodePath) {
		_ = os.Setenv("PATH", nodePath+":"+envPath)
	}
	_ = os.Setenv("NODE_PATH", nodePath)

	// default results collection
	col := utils.GetSpiderCol(r.s.Col, r.s.Name)

	// 默认环境变量
	r.cmd.Env = append(os.Environ(), "CRAWLAB_TASK_ID="+r.t.Id)
	r.cmd.Env = append(r.cmd.Env, "CRAWLAB_COLLECTION="+col)
	r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_HOST="+viper.GetString("mongo.host"))
	r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_PORT="+viper.GetString("mongo.port"))
	if viper.GetString("mongo.db") != "" {
		r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_DB="+viper.GetString("mongo.db"))
	}
	if viper.GetString("mongo.username") != "" {
		r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_USERNAME="+viper.GetString("mongo.username"))
	}
	if viper.GetString("mongo.password") != "" {
		r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_PASSWORD="+viper.GetString("mongo.password"))
	}
	if viper.GetString("mongo.authSource") != "" {
		r.cmd.Env = append(r.cmd.Env, "CRAWLAB_MONGO_AUTHSOURCE="+viper.GetString("mongo.authSource"))
	}
	r.cmd.Env = append(r.cmd.Env, "PYTHONUNBUFFERED=0")
	r.cmd.Env = append(r.cmd.Env, "PYTHONIOENCODING=utf-8")
	r.cmd.Env = append(r.cmd.Env, "TZ=Asia/Shanghai")
	r.cmd.Env = append(r.cmd.Env, "CRAWLAB_DEDUP_FIELD="+r.s.DedupField)
	r.cmd.Env = append(r.cmd.Env, "CRAWLAB_DEDUP_METHOD="+r.s.DedupMethod)
	if r.s.IsDedup {
		r.cmd.Env = append(r.cmd.Env, "CRAWLAB_IS_DEDUP=1")
	} else {
		r.cmd.Env = append(r.cmd.Env, "CRAWLAB_IS_DEDUP=0")
	}

	//任务环境变量
	for _, env := range r.envs {
		r.cmd.Env = append(r.cmd.Env, env.Name+"="+env.Value)
	}

	// 全局环境变量
	variables := model.GetVariableList()
	for _, variable := range variables {
		r.cmd.Env = append(r.cmd.Env, variable.Key+"="+variable.Value)
	}
	return r.cmd
}
