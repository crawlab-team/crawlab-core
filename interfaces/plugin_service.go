package interfaces

type PluginService interface {
	Module
	SetDirPath(path string)
	SetFsPathBase(path string)
}
