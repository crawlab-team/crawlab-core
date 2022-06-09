package controllers

import (
	"github.com/crawlab-team/crawlab-core/sys_exec"
	"github.com/crawlab-team/go-trace"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getDemoActions() []Action {
	ctx := newDemoContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "/import",
			HandlerFunc: ctx.import_,
		},
		{
			Method:      http.MethodGet,
			Path:        "/reimport",
			HandlerFunc: ctx.reimport,
		},
	}
}

type demoContext struct {
}

func (ctx *demoContext) import_(c *gin.Context) {
	cmdStr := "python -m crawlab-demo import"
	cmd := sys_exec.BuildCmd(cmdStr)
	if err := cmd.Run(); err != nil {
		trace.PrintError(err)
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *demoContext) reimport(c *gin.Context) {
	cmdStr := "python -m crawlab-demo reimport"
	cmd := sys_exec.BuildCmd(cmdStr)
	if err := cmd.Run(); err != nil {
		trace.PrintError(err)
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

var _demoCtx *demoContext

func newDemoContext() *demoContext {
	if _demoCtx != nil {
		return _demoCtx
	}

	_demoCtx = &demoContext{}

	return _demoCtx
}

var DemoController ActionController
