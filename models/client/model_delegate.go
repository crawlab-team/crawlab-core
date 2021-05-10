package client

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/entity"
	errors2 "github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/grpc/client"
	"github.com/crawlab-team/crawlab-core/interfaces"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/go-trace"
	"go.uber.org/dig"
)

func NewModelDelegate(doc interfaces.Model, opts ...ModelDelegateOption) interfaces.GrpcClientModelDelegate {
	switch doc.(type) {
	case *models2.Artifact:
		return newModelDelegate(interfaces.ModelIdArtifact, doc)
	case *models2.Tag:
		return newModelDelegate(interfaces.ModelIdTag, doc)
	case *models2.Node:
		return newModelDelegate(interfaces.ModelIdNode, doc)
	case *models2.Project:
		return newModelDelegate(interfaces.ModelIdProject, doc)
	case *models2.Spider:
		return newModelDelegate(interfaces.ModelIdSpider, doc)
	case *models2.Task:
		return newModelDelegate(interfaces.ModelIdTask, doc)
	case *models2.Job:
		return newModelDelegate(interfaces.ModelIdJob, doc)
	case *models2.Schedule:
		return newModelDelegate(interfaces.ModelIdSchedule, doc)
	case *models2.User:
		return newModelDelegate(interfaces.ModelIdUser, doc)
	case *models2.Setting:
		return newModelDelegate(interfaces.ModelIdSetting, doc)
	case *models2.Token:
		return newModelDelegate(interfaces.ModelIdToken, doc)
	case *models2.Variable:
		return newModelDelegate(interfaces.ModelIdVariable, doc)
	default:
		_ = trace.TraceError(errors2.ErrorModelInvalidType)
		return nil
	}
}

func newModelDelegate(id interfaces.ModelId, doc interfaces.Model, opts ...ModelDelegateOption) interfaces.GrpcClientModelDelegate {
	// model delegate
	d := &ModelDelegate{
		id:  id,
		doc: doc,
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(client.NewClient); err != nil {
		_ = trace.TraceError(errors2.ErrorModelInvalidType)
		return nil
	}
	if err := c.Invoke(func(c interfaces.GrpcClient) {
		d.c = c
	}); err != nil {
		_ = trace.TraceError(errors2.ErrorModelInvalidType)
		return nil
	}

	// apply options
	for _, opt := range opts {
		opt(d)
	}

	return d
}

type ModelDelegate struct {
	id  interfaces.ModelId
	c   interfaces.GrpcClient
	doc interfaces.Model
}

func (d *ModelDelegate) Add() (err error) {
	_, err = d.do(interfaces.ModelDelegateMethodAdd)
	return trace.TraceError(err)
}

func (d *ModelDelegate) Save() (err error) {
	_, err = d.do(interfaces.ModelDelegateMethodSave)
	return trace.TraceError(err)
}

func (d *ModelDelegate) Delete() (err error) {
	_, err = d.do(interfaces.ModelDelegateMethodDelete)
	return trace.TraceError(err)
}

func (d *ModelDelegate) GetArtifact() (res interfaces.ModelArtifact, err error) {
	return d.do(interfaces.ModelDelegateMethodGetArtifact)
}

func (d *ModelDelegate) GetModel() (res interfaces.Model) {
	return d.doc
}

func (d *ModelDelegate) Refresh() (err error) {
	return d.refresh()
}

func (d *ModelDelegate) do(method interfaces.ModelDelegateMethod) (a interfaces.ModelArtifact, err error) {
	switch method {
	case interfaces.ModelDelegateMethodAdd:
		return nil, d.add()
	case interfaces.ModelDelegateMethodSave:
		return nil, d.save()
	case interfaces.ModelDelegateMethodDelete:
		return nil, d.delete()
	case interfaces.ModelDelegateMethodGetArtifact:
		return d.getArtifact()
	default:
		return nil, trace.TraceError(errors2.ErrorModelInvalidType)
	}
}

func (d *ModelDelegate) add() (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	_, err = d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.DelegateMessage{
		ModelId: d.id,
		Method:  interfaces.ModelDelegateMethodAdd,
		Data:    d.mustGetData(),
	}))
	return err
}

func (d *ModelDelegate) save() (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	_, err = d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.DelegateMessage{
		ModelId: d.id,
		Method:  interfaces.ModelDelegateMethodSave,
		Data:    d.mustGetData(),
	}))
	return err
}

func (d *ModelDelegate) delete() (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	_, err = d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.DelegateMessage{
		ModelId: d.id,
		Method:  interfaces.ModelDelegateMethodDelete,
		Data:    d.mustGetData(),
	}))
	return err
}

func (d *ModelDelegate) getArtifact() (res2 interfaces.ModelArtifact, err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	res, err := d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.DelegateMessage{
		ModelId: d.id,
		Method:  interfaces.ModelDelegateMethodGetArtifact,
		Data:    d.mustGetData(),
	}))
	if err != nil {
		return nil, err
	}
	var a models2.Artifact
	if err := json.Unmarshal(res.Data, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (d *ModelDelegate) refresh() (err error) {
	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors2.ErrorModelMissingId)
	}
	// TODO: implement
	return nil
}

func (d *ModelDelegate) mustGetData() (data []byte) {
	data, err := d.getData()
	if err != nil {
		panic(err)
	}
	return data
}

func (d *ModelDelegate) getData() (data []byte, err error) {
	return json.Marshal(d.doc)
}
