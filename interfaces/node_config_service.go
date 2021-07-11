package interfaces

type NodeConfigService interface {
	WithConfigPath
	Init() error
	Reload() error
	GetBasicNodeInfo() Entity
	GetNodeKey() string
	IsMaster() bool
	GetAuthKey() string
	GetMaxRunners() int
}
