package interfaces

type NodeConfigService interface {
	Init() error
	Reload() error
	GetBasicNodeInfo() Entity
	GetNodeKey() string
	IsMaster() bool
	SetConfigPath(string)
}
