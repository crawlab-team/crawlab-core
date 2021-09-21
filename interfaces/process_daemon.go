package interfaces

import "os/exec"

type ProcessDaemon interface {
	GetMaxErrors() (maxErrors int)
	SetMaxErrors(maxErrors int)
	GetCmd() (cmd *exec.Cmd)
	GetCh() (ch chan int)
	Start() (err error)
	Stop()
}
