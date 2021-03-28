package controllers

import (
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

var ProjectController = projectController{}

type projectController struct {
}

func (ctr *projectController) Get(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	p, err := model.ProjectService.GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessData(c, p)
}

func (ctr *projectController) GetList(c *gin.Context) {
	panic("implement me")
}

func (ctr *projectController) Post(c *gin.Context) {
	panic("implement me")
}

func (ctr *projectController) PostList(c *gin.Context) {
	panic("implement me")
}

func (ctr *projectController) Put(c *gin.Context) {
	var p model.Project
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	if err := p.Add(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccessData(c, p)
}

func (ctr *projectController) PutList(c *gin.Context) {
	panic("implement me")
}

func (ctr *projectController) Delete(c *gin.Context) {
	panic("implement me")
}

func (ctr *projectController) DeleteList(c *gin.Context) {
	panic("implement me")
}
