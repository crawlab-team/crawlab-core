package controllers

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func NewBasicControllerDelegate(id ControllerId, svc models.PublicServiceInterface) (d *BasicControllerDelegate) {
	return &BasicControllerDelegate{
		id:  id,
		svc: svc,
	}
}

type BasicControllerDelegate struct {
	id  ControllerId
	svc models.PublicServiceInterface
}

func (d *BasicControllerDelegate) Get(c *gin.Context) {
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

func (d *BasicControllerDelegate) Post(c *gin.Context) {
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

func (d *BasicControllerDelegate) Put(c *gin.Context) {
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

func (d *BasicControllerDelegate) Delete(c *gin.Context) {
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
