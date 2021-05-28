package handler

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"time"
)

type Option func(svc interfaces.TaskHandlerService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.TaskHandlerService) {
		svc.SetConfigPath(path)
	}
}

func WithMaxRunners(maxRunners int) Option {
	return func(svc interfaces.TaskHandlerService) {
		svc.SetMaxRunners(maxRunners)
	}
}

func WithExitWatchDuration(duration time.Duration) Option {
	return func(svc interfaces.TaskHandlerService) {
		svc.SetExitWatchDuration(duration)
	}
}

func WithReportInterval(interval time.Duration) Option {
	return func(svc interfaces.TaskHandlerService) {
		svc.SetReportInterval(interval)
	}
}

type RunnerOption func(runner interfaces.TaskRunner)

func WithLogDriverType(driverType string) RunnerOption {
	return func(runner interfaces.TaskRunner) {
		runner.SetLogDriverType(driverType)
	}
}
