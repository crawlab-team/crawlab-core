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
	return &ListControllerDelegate{
		id:  id,
		svc: svc,
	}
}

type ListControllerDelegate struct {
	id  ControllerId
	svc models.PublicServiceInterface
	ListController
}

func (d *ListControllerDelegate) Get(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	doc, err := d.svc.GetById(id)
	if err == mongo2.ErrNoDocuments {
		HandleErrorNotFound(c, err)
		return
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessData(c, doc)
}

func (d *ListControllerDelegate) GetList(c *gin.Context) {
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)
	data, err := d.svc.GetList(query, &mongo.FindOptions{
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
	total, err := d.svc.Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessListData(c, data, total)
}

func (d *ListControllerDelegate) Post(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	doc, err := NewJsonBinder(d.id).Bind(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if doc.GetId() != id {
		HandleErrorBadRequest(c, errors.ErrorHttpBadRequest)
		return
	}
	_, err = d.svc.GetById(id)
	if err != nil {
		HandleErrorNotFound(c, err)
		return
	}
	if err := doc.Save(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessData(c, doc)
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
	if err := d.svc.Update(query, doc); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (d *ListControllerDelegate) Put(c *gin.Context) {
	doc, err := NewJsonBinder(d.id).Bind(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := doc.Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessData(c, doc)
}

func (d *ListControllerDelegate) PutList(c *gin.Context) {
	docs, err := NewJsonBinder(d.id).BindList(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	var ids []primitive.ObjectID

	// reflect
	switch reflect.TypeOf(docs).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(docs)
		for i := 0; i < s.Len(); i++ {
			item := s.Index(i)
			ptr := item.Addr()
			doc, ok := ptr.Interface().(models.BaseModelInterface)
			//doc, ok := item.Interface().(model.BaseModelInterface)
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

	//for _, doc := range docs {
	//	if err := doc.Add(); err != nil {
	//		_ = trace.TraceError(err)
	//		continue
	//	}
	//	ids = append(ids, doc.GetId())
	//}

	items, err := utils.GetArrayItems(docs)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if len(ids) < len(items) {
		HandleErrorInternalServerError(c, errors.ErrorCrudAddError)
		return
	}
	HandleSuccessData(c, docs)
}

func (d *ListControllerDelegate) Delete(c *gin.Context) {
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := d.svc.DeleteById(oid); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
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
