package models

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-db/errors"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

func NewDelegate(colName string, obj interface{}) Delegate {
	var doc BaseModel
	data, err := bson.Marshal(obj)
	if err != nil {
		return Delegate{
			colName: colName,
		}
	}
	if err := bson.Unmarshal(data, &doc); err != nil {
		return Delegate{
			colName: colName,
		}
	}
	a := Artifact{
		Col: colName,
	}
	return Delegate{
		colName: colName,
		obj:     obj,
		doc:     &doc,
		a:       &a,
	}
}

type Delegate struct {
	colName string
	obj     interface{}
	doc     *BaseModel
	a       *Artifact
}

func (d *Delegate) Add() (err error) {
	if d.doc == nil || d.doc.Id.IsZero() {
		return errors.ErrMissingValue
	}
	col := mongo.GetMongoCol(d.colName)
	if _, err = col.Insert(d.obj); err != nil {
		return err
	}
	if err := d.upsertArtifact(); err != nil {
		return err
	}
	return d.refresh()
}

func (d *Delegate) Save() (err error) {
	if d.doc == nil || d.doc.Id.IsZero() {
		return errors.ErrMissingValue
	}
	col := mongo.GetMongoCol(d.colName)
	if err := col.ReplaceId(d.doc.Id, d.obj); err != nil {
		return err
	}
	if err := d.upsertArtifact(); err != nil {
		return err
	}
	return d.refresh()
}

func (d *Delegate) Delete() (err error) {
	if d.doc.Id.IsZero() {
		return trace.TraceError(constants.ErrMissingId)
	}
	col := mongo.GetMongoCol(d.colName)
	if err := col.FindId(d.doc.Id).One(d.obj); err != nil {
		return err
	}
	if err := col.DeleteId(d.doc.Id); err != nil {
		return err
	}
	return d.deleteArtifact()
}

func (d *Delegate) GetArtifact() (a Artifact, err error) {
	if d.doc.Id.IsZero() {
		return a, constants.ErrMissingId
	}
	col := mongo.GetMongoCol(ModelColNameArtifact)
	if err := col.FindId(d.doc.Id).One(&a); err != nil {
		return a, err
	}
	return a, nil
}

func (d *Delegate) refresh() (err error) {
	if d.doc.Id.IsZero() {
		return constants.ErrMissingId
	}
	col := mongo.GetMongoCol(d.colName)
	if err := col.FindId(d.doc.Id).One(d.obj); err != nil {
		return err
	}
	return nil
}

func (d *Delegate) upsertArtifact() (err error) {
	if d.doc.Id.IsZero() {
		return errors.ErrMissingValue
	}
	col := mongo.GetMongoCol(ModelColNameArtifact)
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
				return err
			}
			return nil
		} else {
			// error
			return err
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
		return errors.ErrMissingValue
	}
	col := mongo.GetMongoCol(ModelColNameArtifact)
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
