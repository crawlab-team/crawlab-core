package controllers

import (
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func NewListControllerDelegate(svc model.PublicServiceInterface) (d *ListControllerDelegate) {
	return &ListControllerDelegate{svc: svc}
}

type ListControllerDelegate struct {
	svc model.PublicServiceInterface
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
	panic("implement me")
}

func (d *ListControllerDelegate) Post(c *gin.Context) {
	panic("implement me")
}

func (d *ListControllerDelegate) PostList(c *gin.Context) {
	panic("implement me")
}

func (d *ListControllerDelegate) Put(c *gin.Context) {
	var doc model.BaseModelInterface
	if err := c.ShouldBindJSON(&doc); err != nil {
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
	panic("implement me")
}

func (d *ListControllerDelegate) Delete(c *gin.Context) {
	panic("implement me")
}

func (d *ListControllerDelegate) DeleteList(c *gin.Context) {
	panic("implement me")
}
