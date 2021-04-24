package controllers

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

func NewListControllerDelegate(id ControllerId, svc models.PublicServiceInterface) (d *ListControllerDelegate) {
	if svc == nil {
		panic(errors.ErrorControllerNoModelService)
	}

	return &ListControllerDelegate{
		id:  id,
		svc: svc,
		bc:  NewBasicControllerDelegate(id, svc),
	}
}

type ListControllerDelegate struct {
	id  ControllerId
	svc models.PublicServiceInterface
	bc  BasicController
}

func (d *ListControllerDelegate) Get(c *gin.Context) {
	d.bc.Get(c)
}

func (d *ListControllerDelegate) Post(c *gin.Context) {
	d.bc.Post(c)
}

func (d *ListControllerDelegate) Put(c *gin.Context) {
	d.bc.Put(c)
}

func (d *ListControllerDelegate) Delete(c *gin.Context) {
	d.bc.Delete(c)
}

func (d *ListControllerDelegate) GetList(c *gin.Context) {
	// get all if query field "all" is set true
	all := MustGetFilterAll(c)
	if all {
		d.getAll(c)
		return
	}

	// params
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)

	// get list
	list, err := d.svc.GetList(query, &mongo.FindOptions{
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			HandleErrorNotFound(c, err)
		} else {
			HandleErrorInternalServerError(c, err)
		}
		return
	}
	data := list.Values()

	// total count
	total, err := d.svc.Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// response
	HandleSuccessListData(c, data, total)
}

func (d *ListControllerDelegate) getAll(c *gin.Context) {
	// get list
	list, err := d.svc.GetList(nil, nil)
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			HandleErrorNotFound(c, err)
		} else {
			HandleErrorInternalServerError(c, err)
		}
		return
	}
	data := list.Values()

	// total count
	total, err := d.svc.Count(nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// response
	HandleSuccessListData(c, data, total)
}

func (d *ListControllerDelegate) PostList(c *gin.Context) {
	payload, doc, err := NewJsonBinder(d.id).BindBatchRequestPayloadWithStringData(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	query := bson.M{
		"_id": bson.M{
			"$in": payload.Ids,
		},
	}
	if err := d.svc.Update(query, doc, payload.Fields); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (d *ListControllerDelegate) PutList(c *gin.Context) {
	// bind
	docs, err := NewJsonBinder(d.id).BindList(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// success ids
	var ids []primitive.ObjectID

	// reflect
	switch reflect.TypeOf(docs).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(docs)
		for i := 0; i < s.Len(); i++ {
			item := s.Index(i)
			if !item.CanAddr() {
				HandleErrorInternalServerError(c, errors.ErrorModelInvalidType)
				return
			}
			ptr := item.Addr()
			doc, ok := ptr.Interface().(models.BaseModelInterface)
			if !ok {
				HandleErrorInternalServerError(c, errors.ErrorModelInvalidType)
				return
			}
			if err := doc.Add(); err != nil {
				_ = trace.TraceError(err)
				continue
			}
			ids = append(ids, doc.GetId())
		}
	}

	// check
	items, err := utils.GetArrayItems(docs)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if len(ids) < len(items) {
		HandleErrorInternalServerError(c, errors.ErrorControllerAddError)
		return
	}

	// success
	HandleSuccessData(c, docs)
}

func (d *ListControllerDelegate) DeleteList(c *gin.Context) {
	payload, err := NewJsonBinder(d.id).BindBatchRequestPayload(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := d.svc.DeleteList(bson.M{
		"_id": bson.M{
			"$in": payload.Ids,
		},
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}
