package interfaces

type EventFn func(data ...interface{}) (err error)

type EventService interface {
	Register(key string, ch chan EventData)
	Unregister(key string)
	SendEvent(eventName string, data ...interface{})
}
