package interfaces

type ModelArtifact interface {
	BaseModelInterface
	GetTags() (res []Tag, err error)
	UpdateTags(tagNames []string) (err error)
	GetSys() (sys ModelArtifactSys)
}
