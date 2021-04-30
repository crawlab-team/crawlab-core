package interfaces

type Store interface {
	Set(string, interface{}) error
	MustSet(string, interface{})
	Get(string) (interface{}, error)
	MustGet(string) interface{}
	GetDefault() (interface{}, error)
	MustGetDefault() interface{}
}
