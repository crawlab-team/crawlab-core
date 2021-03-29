package controllers

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
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
	if err == mongo2.ErrNoDocuments {
		HandleErrorNotFound(c, err)
		return
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessData(c, p)
}

func (ctr *projectController) GetList(c *gin.Context) {
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)
	data, err := model.ProjectService.GetList(query, &mongo.FindOptions{
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	total, err := model.ProjectService.Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessListData(c, data, total)
}

func (ctr *projectController) Post(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	var p model.Project
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if p.Id != id {
		HandleErrorBadRequest(c, errors.ErrorHttpBadRequest)
		return
	}
	_, err = model.ProjectService.GetById(id)
	if err != nil {
		HandleErrorNotFound(c, err)
		return
	}
	if err := p.Save(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessData(c, p)
}

func (ctr *projectController) PostList(c *gin.Context) {
	var payload entity.BatchRequestPayloadWithStringData
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	var p model.Project
	if err := json.Unmarshal([]byte(payload.Data), &p); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	query := bson.M{
		"_id": bson.M{
			"$in": payload.Ids,
		},
	}
	if err := model.ProjectService.UpdateList(query, p); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctr *projectController) Put(c *gin.Context) {
	var p model.Project
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := p.Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessData(c, p)
}

func (ctr *projectController) PutList(c *gin.Context) {
	var docs []model.Project
	if err := c.ShouldBindJSON(&docs); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	var resDocs []model.Project
	for _, p := range docs {
		if err := p.Add(); err != nil {
			_ = trace.TraceError(err)
			continue
		}
		resDocs = append(resDocs, p)
	}
	if len(resDocs) < len(docs) {
		HandleErrorInternalServerError(c, errors.ErrorCrudAddError)
		return
	}
	HandleSuccessData(c, resDocs)
}

func (ctr *projectController) Delete(c *gin.Context) {
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := model.ProjectService.DeleteById(oid); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctr *projectController) DeleteList(c *gin.Context) {
	var payload entity.BatchRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := model.ProjectService.DeleteList(bson.M{
		"_id": bson.M{
			"$in": payload.Ids,
		},
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}
