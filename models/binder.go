package models

import "github.com/crawlab-team/crawlab-core/interfaces"

type Binder interface {
	Bind() (res interface{}, err error)
	MustBind() (res interface{})
	process(d interface{}, fieldIds ...interfaces.ModelId) (res interface{}, err error)
}
