package interfaces

type ModelArtifact interface {
	GetTags() (res interface{}, err error)
	UpdateTags(tagNames []string) (err error)
}
