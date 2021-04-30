package grpc

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/utils/binders"
)

func NewDelegate(id interfaces.ModelId, obj interface{}) Delegate {
	colName := binders.NewColNameBinder(id).MustBindString()
	return Delegate{
		id:      id,
		colName: colName,
		obj:     obj,
	}
}

type Delegate struct {
	id      interfaces.ModelId
	colName string
	obj     interface{}
}

func (d *Delegate) Add() error {
	panic("implement me")
}

func (d *Delegate) Save() error {
	panic("implement me")
}

func (d *Delegate) Delete() error {
	panic("implement me")
}
