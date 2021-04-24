package models

type Binder interface {
	Bind() (res interface{}, err error)
	MustBind() (res interface{})
	process(d interface{}, fieldIds ...ModelId) (res interface{}, err error)
}
