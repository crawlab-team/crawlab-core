package interfaces

type ModelDelegateMethod int

type ModelDelegate interface {
	Add() error
	Save() error
	Delete() error
	GetArtifact() (ModelArtifact, error)
	GetModel() Model
	Refresh() error
}

const (
	ModelDelegateMethodAdd = iota
	ModelDelegateMethodSave
	ModelDelegateMethodDelete
	ModelDelegateMethodGetArtifact
	ModelDelegateMethodRefresh
)
