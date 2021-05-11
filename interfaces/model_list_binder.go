package interfaces

import "github.com/emirpasic/gods/lists/arraylist"

type modelListBinder interface {
	Bind() (list arraylist.List, err error)
	Process(d interface{}) (list arraylist.List, err error)
}

type ModelListBinder interface {
	modelListBinder
	MustBindListAsPtr() (res arraylist.List)
}
