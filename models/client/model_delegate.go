package client

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/entity"
	errors2 "github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/grpc/client"
	"github.com/crawlab-team/crawlab-core/interfaces"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/go-trace"
)

func NewModelDelegate(doc interfaces.Model, opts ...ModelDelegateOption) interfaces.GrpcClientModelDelegate {
	switch doc.(type) {
	case *models2.Artifact:
		return newModelDelegate(interfaces.ModelIdArtifact, doc, opts...)
	case *models2.Tag:
		return newModelDelegate(interfaces.ModelIdTag, doc, opts...)
	case *models2.Node:
		return newModelDelegate(interfaces.ModelIdNode, doc, opts...)
	case *models2.Project:
		return newModelDelegate(interfaces.ModelIdProject, doc, opts...)
	case *models2.Spider:
		return newModelDelegate(interfaces.ModelIdSpider, doc, opts...)
	case *models2.Task:
		return newModelDelegate(interfaces.ModelIdTask, doc, opts...)
	case *models2.Job:
		return newModelDelegate(interfaces.ModelIdJob, doc, opts...)
	case *models2.Schedule:
		return newModelDelegate(interfaces.ModelIdSchedule, doc, opts...)
	case *models2.User:
		return newModelDelegate(interfaces.ModelIdUser, doc, opts...)
	case *models2.Setting:
		return newModelDelegate(interfaces.ModelIdSetting, doc, opts...)
	case *models2.Token:
		return newModelDelegate(interfaces.ModelIdToken, doc, opts...)
	case *models2.Variable:
		return newModelDelegate(interfaces.ModelIdVariable, doc, opts...)
	default:
		_ = trace.TraceError(errors2.ErrorModelInvalidType)
		return nil
	}
}

func newModelDelegate(id interfaces.ModelId, doc interfaces.Model, opts ...ModelDelegateOption) interfaces.GrpcClientModelDelegate {
	var err error

	// model delegate
	d := &ModelDelegate{
		id:      id,
		doc:     doc,
		cfgPath: config.DefaultConfigPath,
	}

	// apply options
	for _, opt := range opts {
		opt(d)
	}

	// grpc client
	d.c, err = client.GetClient(d.cfgPath)
	if err != nil {
		trace.PrintError(errors2.ErrorModelInvalidType)
		return nil
	}
	if !d.c.IsStarted() {
		if err := d.c.Start(); err != nil {
			trace.PrintError(err)
			return nil
		}
	} else if d.c.IsClosed() {
		if err := d.c.Restart(); err != nil {
			trace.PrintError(err)
			return nil
		}
	}

	return d
}

type ModelDelegate struct {
	// settings
	cfgPath string

	// internals
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

func (d *ModelDelegate) GetConfigPath() (path string) {
	return d.cfgPath
}

func (d *ModelDelegate) SetConfigPath(path string) {
	d.cfgPath = path
}

func (d *ModelDelegate) Close() (err error) {
	return d.c.Stop()
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
	_, err = d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.GrpcDelegateMessage{
		ModelId: d.id,
		Method:  interfaces.ModelDelegateMethodAdd,
		Data:    d.mustGetData(),
	}))
	return err
}

func (d *ModelDelegate) save() (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	_, err = d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.GrpcDelegateMessage{
		ModelId: d.id,
		Method:  interfaces.ModelDelegateMethodSave,
		Data:    d.mustGetData(),
	}))
	return err
}

func (d *ModelDelegate) delete() (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	_, err = d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.GrpcDelegateMessage{
		ModelId: d.id,
		Method:  interfaces.ModelDelegateMethodDelete,
		Data:    d.mustGetData(),
	}))
	return err
}

func (d *ModelDelegate) getArtifact() (res2 interfaces.ModelArtifact, err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	res, err := d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.GrpcDelegateMessage{
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
