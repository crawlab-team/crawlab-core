package interfaces

type NodeService interface {
	Init() error
	Reload() error
	GetNodeInfo() Entity
	GetNodeKey() string
	IsMaster() bool
}
