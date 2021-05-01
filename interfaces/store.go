package interfaces

import "sync"

type Store interface {
	GetMap() sync.Map
	Set(interface{}, interface{}) error
	MustSet(interface{}, interface{})
	Get(interface{}) (interface{}, error)
	MustGet(interface{}) interface{}
	GetDefault() (interface{}, error)
	MustGetDefault() interface{}
}
