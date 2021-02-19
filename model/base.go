package model

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-db/errors"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type BaseModelInterface interface {
	Add() (err error)
	Save() (err error)
	Delete() (err error)
	GetArtifact() (a Artifact, err error)
}

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

type BaseModel struct {
	Id primitive.ObjectID `bson:"_id" json:"_id"`
}

type Delegate struct {
	colName string
	obj     interface{}
	doc     *BaseModel
	a       *Artifact
}

const ArtifactColName = "artifacts"

type Artifact struct {
	Id  primitive.ObjectID `bson:"_id" json:"_id"`
	Col string             `bson:"_col" json:"_col"`
	Del bool               `bson:"_del" json:"_del"`
	Sys `bson:"_sys" json:"_sys"`
	Obj interface{} `bson:"_obj" json:"_obj"`
}

type Sys struct {
	CreateTs  time.Time          `json:"create_ts" bson:"create_ts"`
	CreateUid primitive.ObjectID `json:"create_uid" bson:"create_uid"`
	UpdateTs  time.Time          `json:"update_ts" bson:"update_ts"`
	UpdateUid primitive.ObjectID `json:"update_uid" bson:"update_uid"`
	DeleteTs  time.Time          `json:"delete_ts" bson:"delete_ts"`
	DeleteUid primitive.ObjectID `json:"delete_uid" bson:"delete_uid"`
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
		return constants.ErrMissingId
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
	col := mongo.GetMongoCol(ArtifactColName)
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
	col := mongo.GetMongoCol(ArtifactColName)
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
	col := mongo.GetMongoCol(ArtifactColName)
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
