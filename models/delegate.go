package models

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/entity"
	errors2 "github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node"
	"github.com/crawlab-team/crawlab-core/store"
	"github.com/crawlab-team/crawlab-db/errors"
	"github.com/crawlab-team/crawlab-db/mongo"
	grpc2 "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

func NewDelegate(id interfaces.ModelId, obj interface{}) Delegate {
	colName := getModelColName(id)
	var doc BaseModel
	data, err := json.Marshal(obj)
	if err != nil {
		return Delegate{
			colName: colName,
		}
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return Delegate{
			colName: colName,
		}
	}
	a := Artifact{
		Col: colName,
	}
	return Delegate{
		id:      id,
		colName: colName,
		obj:     obj,
		doc:     &doc,
		a:       &a,
	}
}

type Delegate struct {
	id      interfaces.ModelId
	colName string
	obj     interface{}
	doc     *BaseModel
	a       *Artifact
}

func (d *Delegate) do(method interfaces.ModelDelegateMethod) (a interfaces.ModelArtifact, err error) {
	if store.NodeService.IsMaster() {
		return d.doLocal(method)
	} else {
		return d.doRemote(method)
	}
}

func (d *Delegate) doLocal(method interfaces.ModelDelegateMethod) (a interfaces.ModelArtifact, err error) {
	switch method {
	case interfaces.ModelDelegateMethodAdd:
		return a, d.add()
	case interfaces.ModelDelegateMethodSave:
		return a, d.save()
	case interfaces.ModelDelegateMethodDelete:
		return a, d.delete()
	case interfaces.ModelDelegateMethodGetArtifact:
		return d.getArtifact()
	default:
		return a, trace.TraceError(errors2.ErrorModelInvalidType)
	}
}

func (d *Delegate) doRemote(method interfaces.ModelDelegateMethod) (res interfaces.ModelArtifact, err error) {
	var a Artifact
	data, err := json.Marshal(d.obj)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	msg := entity.DelegateMessage{
		ModelId: d.id,
		Method:  method,
		Data:    data,
	}
	client, err := store.GrpcService.GetDefaultClient()
	if err != nil {
		return nil, trace.TraceError(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // TODO: configure timeout
	defer cancel()
	nodeKey, err := node.GetNodeKey()
	if err != nil {
		return nil, trace.TraceError(err)
	}
	req := &grpc2.Request{
		NodeKey: nodeKey,
		Data:    msg.ToBytes(),
	}
	_, err = client.GetModelDelegateClient().Do(ctx, req)
	if err != nil {
		return res, trace.TraceError(err)
	}
	return &a, nil
}

func (d *Delegate) Add() (err error) {
	_, err = d.do(interfaces.ModelDelegateMethodAdd)
	return trace.TraceError(err)
}

func (d *Delegate) Save() (err error) {
	_, err = d.do(interfaces.ModelDelegateMethodSave)
	return trace.TraceError(err)
}

func (d *Delegate) Delete() (err error) {
	_, err = d.do(interfaces.ModelDelegateMethodDelete)
	return trace.TraceError(err)
}

func (d *Delegate) GetArtifact() (res interfaces.ModelArtifact, err error) {
	return d.do(interfaces.ModelDelegateMethodGetArtifact)
}

func (d *Delegate) add() (err error) {
	if d.doc == nil || d.doc.Id.IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}
	col := mongo.GetMongoCol(d.colName)
	if _, err = col.Insert(d.obj); err != nil {
		return trace.TraceError(err)
	}
	if err := d.upsertArtifact(); err != nil {
		return trace.TraceError(err)
	}
	if err := d.updateTags(); err != nil {
		return trace.TraceError(err)
	}
	return d.refresh()
}

func (d *Delegate) save() (err error) {
	if d.doc == nil || d.doc.Id.IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}
	col := mongo.GetMongoCol(d.colName)
	if err := col.ReplaceId(d.doc.Id, d.obj); err != nil {
		return trace.TraceError(err)
	}
	if err := d.upsertArtifact(); err != nil {
		return trace.TraceError(err)
	}
	if err := d.updateTags(); err != nil {
		return trace.TraceError(err)
	}
	return d.refresh()
}

func (d *Delegate) delete() (err error) {
	if d.doc.Id.IsZero() {
		return trace.TraceError(errors2.ErrorModelMissingId)
	}
	col := mongo.GetMongoCol(d.colName)
	if err := col.FindId(d.doc.Id).One(d.obj); err != nil {
		return trace.TraceError(err)
	}
	if err := col.DeleteId(d.doc.Id); err != nil {
		return trace.TraceError(err)
	}
	return d.deleteArtifact()
}

func (d *Delegate) getArtifact() (res interfaces.ModelArtifact, err error) {
	var a Artifact
	if d.doc.Id.IsZero() {
		return &a, trace.TraceError(errors2.ErrorModelMissingId)
	}
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	if err := col.FindId(d.doc.Id).One(a); err != nil {
		return &a, err
	}
	return &a, nil
}

func (d *Delegate) refresh() (err error) {
	if d.doc.Id.IsZero() {
		return trace.TraceError(errors2.ErrorModelMissingId)
	}
	col := mongo.GetMongoCol(d.colName)
	fr := col.FindId(d.doc.Id)
	if err := fr.One(d.obj); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (d *Delegate) upsertArtifact() (err error) {
	if d.doc.Id.IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	ctx := col.GetContext()
	// TODO: implement user
	user, ok := ctx.Value(UserContextKey).(*User)
	d.a.Id = d.doc.Id
	if err := col.FindId(d.doc.Id).One(d.a); err != nil {
		if err == mongo2.ErrNoDocuments {
			// new artifact
			d.a.CreateTs = time.Now()
			d.a.UpdateTs = time.Now()
			if ok {
				d.a.CreateUid = user.Id
				d.a.UpdateUid = user.Id
			}
			_, err = col.Insert(d.a)
			if err != nil {
				return trace.TraceError(err)
			}
			return nil
		} else {
			// error
			return trace.TraceError(err)
		}
	}

	// existing artifact
	d.a.UpdateTs = time.Now()
	if ok {
		d.a.UpdateUid = user.Id
	}
	return col.ReplaceId(d.a.Id, d.a)
}

func (d *Delegate) deleteArtifact() (err error) {
	if d.doc.Id.IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	ctx := col.GetContext()
	d.a.Id = d.doc.Id
	d.a.Obj = d.obj
	d.a.Del = true
	d.a.DeleteTs = time.Now()
	// TODO: implement user
	user, ok := ctx.Value(UserContextKey).(*User)
	if ok {
		d.a.DeleteUid = user.Id
	}
	return col.ReplaceId(d.doc.Id, d.a)
}

func (d *Delegate) updateTags() (err error) {
	if d.doc.Id.IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}
	//ctx := col.GetContext()
	if _, err := TagService.UpdateTagsById(d.colName, d.doc.Id, d.doc.Tags); err != nil {
		return trace.TraceError(err)
	}

	return nil
}
