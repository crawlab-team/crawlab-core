package interfaces

type ModelDelegateMethod int

type ModelDelegate interface {
	Add() error
	Save() error
	Delete() error
	GetArtifact() (ModelArtifact, error)
}

type ModelDelegateMessage interface {
	GetModelId() ModelId
	GetMethod() ModelDelegateMethod
	GetData() []byte
	ToBytes() (data []byte)
}

const (
	ModelDelegateMethodAdd = iota
	ModelDelegateMethodSave
	ModelDelegateMethodDelete
	ModelDelegateMethodGetArtifact
)
