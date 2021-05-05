package interfaces

import "github.com/emirpasic/gods/lists/arraylist"

type ModelListBinder interface {
	Bind() (res interface{}, err error)
	MustBind() (res interface{})
	BindList() (res arraylist.List, err error)
	MustBindList() (res arraylist.List)
	MustBindListAsPtr() (res arraylist.List)
	BindListAsPtr() (res arraylist.List, err error)
	Process(d interface{}, fieldIds ...ModelId) (res interface{}, err error)
}
