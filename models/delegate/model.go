package delegate

import (
	errors2 "github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/errors"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

func NewModelDelegate(doc interfaces.Model) interfaces.ModelDelegate {
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

func newModelDelegate(id interfaces.ModelId, doc interfaces.Model) interfaces.ModelDelegate {
	// collection name
	colName := models2.GetModelColName(id)

	// model delegate
	d := &ModelDelegate{
		id:  id,
		doc: doc,
		a: &models2.Artifact{
			Col: colName,
		},
		colName: colName,
	}

	return d
}

type ModelDelegate struct {
	id      interfaces.ModelId
	colName string
	doc     interfaces.Model
	a       interfaces.ModelArtifact
}

// Add model
func (d *ModelDelegate) Add() (err error) {
	return d.do(interfaces.ModelDelegateMethodAdd)
}

// Save model
func (d *ModelDelegate) Save() (err error) {
	return d.do(interfaces.ModelDelegateMethodSave)
}

// Delete model
func (d *ModelDelegate) Delete() (err error) {
	return d.do(interfaces.ModelDelegateMethodDelete)
}

// GetArtifact refresh artifact and return it
func (d *ModelDelegate) GetArtifact() (res interfaces.ModelArtifact, err error) {
	if err := d.do(interfaces.ModelDelegateMethodGetArtifact); err != nil {
		return nil, err
	}
	return d.a, nil
}

// Refresh model
func (d *ModelDelegate) Refresh() (err error) {
	return d.refresh()
}

// GetModel return model
func (d *ModelDelegate) GetModel() (res interfaces.Model) {
	return d.doc
}

// do action given the model delegate method
func (d *ModelDelegate) do(method interfaces.ModelDelegateMethod) (err error) {
	switch method {
	case interfaces.ModelDelegateMethodAdd:
		return d.add()
	case interfaces.ModelDelegateMethodSave:
		return d.save()
	case interfaces.ModelDelegateMethodDelete:
		return d.delete()
	case interfaces.ModelDelegateMethodGetArtifact, interfaces.ModelDelegateMethodRefresh:
		return d.refresh()
	default:
		return trace.TraceError(errors2.ErrorModelInvalidType)
	}
}

// add model
func (d *ModelDelegate) add() (err error) {
	if d.doc == nil {
		return trace.TraceError(errors.ErrMissingValue)
	}
	if d.doc.GetId().IsZero() {
		d.doc.SetId(primitive.NewObjectID())
	}
	col := mongo.GetMongoCol(d.colName)
	if _, err = col.Insert(d.doc); err != nil {
		return trace.TraceError(err)
	}
	if err := d.upsertArtifact(); err != nil {
		return trace.TraceError(err)
	}
	// TODO: implement with alternative
	if err := d.updateTags(); err != nil {
		return trace.TraceError(err)
	}
	return d.refresh()
}

// save model
func (d *ModelDelegate) save() (err error) {
	if d.doc == nil || d.doc.GetId().IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}
	col := mongo.GetMongoCol(d.colName)
	if err := col.ReplaceId(d.doc.GetId(), d.doc); err != nil {
		return trace.TraceError(err)
	}
	if err := d.upsertArtifact(); err != nil {
		return trace.TraceError(err)
	}
	// TODO: implement with alternative
	//if err := d.updateTags(); err != nil {
	//	return trace.TraceError(err)
	//}
	return d.refresh()
}

// delete model
func (d *ModelDelegate) delete() (err error) {
	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors2.ErrorModelMissingId)
	}
	col := mongo.GetMongoCol(d.colName)
	if err := col.FindId(d.doc.GetId()).One(d.doc); err != nil {
		return trace.TraceError(err)
	}
	if err := col.DeleteId(d.doc.GetId()); err != nil {
		return trace.TraceError(err)
	}
	return d.deleteArtifact()
}

// refresh model and artifact
func (d *ModelDelegate) refresh() (err error) {
	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors2.ErrorModelMissingId)
	}
	col := mongo.GetMongoCol(d.colName)
	fr := col.FindId(d.doc.GetId())
	if err := fr.One(d.doc); err != nil {
		return trace.TraceError(err)
	}
	return d.refreshArtifact()
}

// refresh artifact
func (d *ModelDelegate) refreshArtifact() (err error) {
	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors2.ErrorModelMissingId)
	}
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	if err := col.FindId(d.doc.GetId()).One(d.a); err != nil {
		return err
	}
	return nil
}

// upsertArtifact
func (d *ModelDelegate) upsertArtifact() (err error) {
	// skip artifact
	if d.id == interfaces.ModelIdArtifact {
		return nil
	}

	// validate id
	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}

	// mongo collection
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)

	// context
	// TODO: implement user
	ctx := col.GetContext()
	user, ok := ctx.Value(models2.UserContextKey).(*models2.User)

	// assign id to artifact
	d.a.SetId(d.doc.GetId())

	// attempt to find artifact
	if err := col.FindId(d.doc.GetId()).One(d.a); err != nil {
		if err == mongo2.ErrNoDocuments {
			// new artifact
			d.a.GetSys().SetCreateTs(time.Now())
			d.a.GetSys().SetUpdateTs(time.Now())
			if ok {
				d.a.GetSys().SetCreateUid(user.Id)
				d.a.GetSys().SetUpdateUid(user.Id)
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
	d.a.GetSys().SetUpdateTs(time.Now())
	if ok {
		d.a.GetSys().SetUpdateUid(user.Id)
	}

	// save new artifact
	return col.ReplaceId(d.a.GetId(), d.a)
}

// deleteArtifact
func (d *ModelDelegate) deleteArtifact() (err error) {
	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	ctx := col.GetContext()
	d.a.SetId(d.doc.GetId())
	d.a.SetObj(d.doc)
	d.a.SetDel(true)
	d.a.GetSys().SetDeleteTs(time.Now())
	// TODO: implement user
	user, ok := ctx.Value(models2.UserContextKey).(*models2.User)
	if ok {
		d.a.GetSys().SetDeleteUid(user.Id)
	}
	return col.ReplaceId(d.doc.GetId(), d.a)
}

// updateTags
func (d *ModelDelegate) updateTags() (err error) {
	// validate id
	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}
	//ctx := col.GetContext()

	// convert to model with tags
	doc, ok := d.doc.(interfaces.ModelWithTags)
	if !ok {
		return nil
	}

	// skip if not tags
	if doc.GetTags() == nil || len(doc.GetTags()) == 0 {
		return nil
	}

	// upsert tags and add to tag ids
	var tagIds []primitive.ObjectID
	for _, tag := range doc.GetTags() {
		if tag.GetId().IsZero() {
			tag.SetCol(d.colName)
			if err := NewModelDelegate(tag).Add(); err != nil {
				return err
			}
		}
		tagIds = append(tagIds, tag.GetId())
	}

	// assign tag ids to artifact
	d.a.SetTagIds(tagIds)

	// update tag ids
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	return col.ReplaceId(d.a.GetId(), d.a)
}
