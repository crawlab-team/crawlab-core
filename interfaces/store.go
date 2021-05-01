package interfaces

type Store interface {
	Set(interface{}, interface{}) error
	MustSet(interface{}, interface{})
	Get(interface{}) (interface{}, error)
	MustGet(interface{}) interface{}
	GetDefault() (interface{}, error)
	MustGetDefault() interface{}
}
