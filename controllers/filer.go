package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

var FilerController ActionController

func getFilerActions() []Action {
	filerCtx := newFilerContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "",
			HandlerFunc: filerCtx.get,
		},
		//{
		//	Method:      http.MethodPost,
		//	Path:        "",
		//	HandlerFunc: filerCtx.post,
		//},
		//{
		//	Method:      http.MethodPut,
		//	Path:        "",
		//	HandlerFunc: filerCtx.put,
		//},
		//{
		//	Method:      http.MethodDelete,
		//	Path:        "",
		//	HandlerFunc: filerCtx.del,
		//},
	}
}

type filerContext struct {
	endpoint string
}

func (ctx *filerContext) get(c *gin.Context) {
	requestPath := strings.Replace(c.Request.URL.Path, "/filer", "/", 1)
	requestUrl := fmt.Sprintf("%s%s", ctx.endpoint, requestPath)
	if c.Request.URL.RawQuery != "" {
		requestUrl += "?" + c.Request.URL.RawQuery
	}
	res, err := req.Get(requestUrl)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	c.Request.Response = res.Response()
	c.AbortWithStatus(http.StatusOK)
}

var _filerCtx *filerContext

func newFilerContext() *filerContext {
	if _filerCtx != nil {
		return _filerCtx
	}

	ctx := &filerContext{
		endpoint: "http://localhost:8888",
	}

	if viper.GetString("fs.filer.proxy") != "" {
		ctx.endpoint = viper.GetString("fs.filer.proxy")
	}

	_filerCtx = ctx

	return ctx
}
