package controllers

import (
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func NewBasicControllerDelegate(svc model.PublicServiceInterface) (d *BasicControllerDelegate) {
	return &BasicControllerDelegate{svc: svc}
}

type BasicControllerDelegate struct {
	svc model.PublicServiceInterface
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
	panic("implement me")
}

func (d *BasicControllerDelegate) Put(c *gin.Context) {
	panic("implement me")
}

func (d *BasicControllerDelegate) Delete(c *gin.Context) {
	panic("implement me")
}
