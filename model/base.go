package model

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BaseModelInterface interface {
	Add() (err error)
	Save() (err error)
	Delete() (err error)
}

func NewDelegate(id primitive.ObjectID, colName string, obj interface{}, sys *Sys) Delegate {
	return Delegate{
		id:      id,
		colName: colName,
		obj:     obj,
		sys:     sys,
	}
}

type Delegate struct {
	id      primitive.ObjectID
	colName string
	obj     interface{}
	sys     *Sys
}

type Sys struct {
	CreateTs  time.Time          `json:"create_ts" bson:"create_ts"`
	CreateUid primitive.ObjectID `json:"create_uid" bson:"create_uid"`
	UpdateTs  time.Time          `json:"update_ts" bson:"update_ts"`
	UpdateUid primitive.ObjectID `json:"update_uid" bson:"update_uid"`
}

func (d *Delegate) Add() (err error) {
	col := mongo.GetMongoCol(JobColName)
	d.id, err = col.Insert(d)
	d.sys.CreateTs = time.Now()
	d.sys.UpdateTs = time.Now()
	// TODO: CreateUid & UpdateUid
	if err != nil {
		return err
	}
	return d.refresh()
}

func (d *Delegate) Save() (err error) {
	col := mongo.GetMongoCol(JobColName)
	d.sys.UpdateTs = time.Now()
	// TODO: UpdateUid
	if err := col.ReplaceId(d.id, d); err != nil {
		return err
	}
	return d.refresh()
}

func (d *Delegate) Delete() (err error) {
	col := mongo.GetMongoCol(JobColName)
	return col.DeleteId(d.id)
}

func (d *Delegate) refresh() (err error) {
	if d.id.IsZero() {
		return constants.ErrMissingId
	}
	col := mongo.GetMongoCol(JobColName)
	res, err := col.FindId(d.id)
	if err != nil {
		return err
	}
	return res.One(d.obj)
}
