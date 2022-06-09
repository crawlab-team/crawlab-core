package utils

import (
	"github.com/crawlab-team/crawlab-core/sys_exec"
	"github.com/crawlab-team/go-trace"
)

func ImportDemo() (err error) {
	cmdStr := "crawlab-demo import && python -m crawlab-demo import"
	cmd := sys_exec.BuildCmd(cmdStr)
	if err := cmd.Run(); err != nil {
		return trace.TraceError(err)
	}
	return nil
}
