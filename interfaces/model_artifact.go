package interfaces

type ModelArtifact interface {
	GetTags() (res []Tag, err error)
	UpdateTags(tagNames []string) (err error)
}
