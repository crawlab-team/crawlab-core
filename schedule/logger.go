package schedule

import (
	"github.com/apex/log"
	"github.com/crawlab-team/go-trace"
	"github.com/robfig/cron/v3"
)

type Logger struct {
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	log.Infof(msg, keysAndValues...)
}

func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	log.Errorf(msg, keysAndValues...)
	trace.PrintError(err)
}

func NewLogger() cron.Logger {
	return &Logger{}
}
