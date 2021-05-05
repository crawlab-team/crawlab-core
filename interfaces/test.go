package interfaces

type Test interface {
	Injectable
	Setup()
	Cleanup()
}

type NodeServiceTest Test
