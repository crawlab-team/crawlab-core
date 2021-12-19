package interfaces

type Plugin interface {
	Model
	GetName() (name string)
	SetName(name string)
	GetFullName() (fullName string)
	SetFullName(fullName string)
	GetInstallUrl() (url string)
	SetInstallUrl(url string)
	GetInstallType() (t string)
	SetInstallType(t string)
}
