package interfaces

import (
	"os/exec"
	"time"
)

type ProcessDaemon interface {
	GetMaxErrors() (maxErrors int)
	SetMaxErrors(maxErrors int)
	GetExitTimeout() (timeout time.Duration)
	SetExitTimeout(timeout time.Duration)
	GetCmd() (cmd *exec.Cmd)
	GetCh() (ch chan int)
	Start() (err error)
	Stop()
}
