package services

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/crawlab-team/crawlab-core/services/sys_exec"
	"github.com/spf13/viper"
	"os/exec"
	"path/filepath"
)

type TaskRunnerInterface interface {
	Run() (err error)
	Cancel() (err error)
	Dispose() (err error)
}

type TaskRunnerOptions struct {
	Task    *model.Task               // Task to run
	Channel chan constants.TaskSignal // Channel to communicate between TaskService and TaskRunner
}

func NewTaskRunner(options *TaskRunnerOptions) (r *TaskRunner, err error) {
	if options.Task == nil {
		return r, constants.ErrInvalidOptions
	}
	spider, err := model.GetSpider(options.Task.SpiderId)
	if err != nil {
		return r, err
	}
	r = &TaskRunner{
		t:    options.Task,
		s:    &spider,
		opts: options,
	}
	return r, nil
}

type TaskRunner struct {
	cmd    *exec.Cmd                 // process command instance
	pid    int                       // process id
	killed bool                      // whether process is killed
	t      *model.Task               // task
	s      *model.Spider             // spider
	ch     chan constants.TaskSignal // channel to communicate between TaskService and TaskRunner
	opts   *TaskRunnerOptions        // options
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

	// TODO: configure log settings

	// TODO: configure environment variables

	// configure to allow kill sub processes
	sys_exec.Setpgid(r.cmd)

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
	// TODO: make sure the process does not exist
	return nil
}

func (r *TaskRunner) Dispose() (err error) {
	panic("implement me")
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
