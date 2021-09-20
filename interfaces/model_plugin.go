package interfaces

type Plugin interface {
	Model
	GetName() (name string)
	SetName(name string)
	GetInstallUrl() (url string)
	SetInstallUrl(url string)
}
