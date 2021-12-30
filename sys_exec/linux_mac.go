// +build !windows

package sys_exec

import (
	"github.com/crawlab-team/go-trace"
	"github.com/shirou/gopsutil/process"
	"os/exec"
	"syscall"
	"time"
)

func BuildCmd(cmdStr string) *exec.Cmd {
	return exec.Command("sh", "-c", cmdStr)
}

func Setpgid(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	} else {
		cmd.SysProcAttr.Setpgid = true
	}
}

func KillProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	return ForceKillProcess(cmd)
}

func KillProcessWithTimeout(cmd *exec.Cmd, timeout time.Duration) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	go func() {
		if err := syscall.Kill(cmd.Process.Pid, syscall.SIGTERM); err != nil {
			trace.PrintError(err)
		}
	}()
	for i := 0; i < int(timeout.Seconds()); i++ {
		ok, err := process.PidExists(int32(cmd.Process.Pid))
		if err == nil && !ok {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return ForceKillProcess(cmd)
}

func ForceKillProcess(cmd *exec.Cmd) error {
	return syscall.Kill(cmd.Process.Pid, syscall.SIGKILL)
}
