package models

type Binder interface {
	Bind() (res interface{}, err error)
	MustBind() (res interface{})
}
