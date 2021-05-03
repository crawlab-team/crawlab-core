package interfaces

type ModelArtifact interface {
	BaseModelInterface
	GetTags() (res []Tag, err error)
	GetSys() (sys ModelArtifactSys)
}
