package controllers

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/gin-gonic/gin"
)

type BinderInterface interface {
	Bind(c *gin.Context) (res models.BaseModelInterface, err error)
	BindList(c *gin.Context) (res []models.BaseModelInterface, err error)
	BindBatchRequestPayload(c *gin.Context) (payload entity.BatchRequestPayload, err error)
	BindBatchRequestPayloadWithStringData(c *gin.Context) (payload entity.BatchRequestPayloadWithStringData, res models.BaseModelInterface, err error)
}

func NewJsonBinder(id ControllerId) (b *JsonBinder) {
	var p models.Project
	var pl []*models.Project

	var u models.User
	var ul []models.User

	return &JsonBinder{
		id: id,
		docMap: map[ControllerId]models.BaseModelInterface{
			ControllerIdProject: &p,
			ControllerIdUser:    &u,
		},
		docListMap: map[ControllerId]interface{}{
			ControllerIdProject: &pl,
			ControllerIdUser:    &ul,
		},
	}
}

type JsonBinder struct {
	id         ControllerId
	docMap     map[ControllerId]models.BaseModelInterface
	docListMap map[ControllerId]interface{}
}

func (b *JsonBinder) Bind(c *gin.Context) (res models.BaseModelInterface, err error) {
	switch b.id {
	case ControllerIdProject:
		var doc models.Project
		err = c.ShouldBindJSON(&doc)
		return &doc, err
	case ControllerIdUser:
		var doc models.User
		err = c.ShouldBindJSON(&doc)
		return &doc, err
	default:
		return nil, errors.ErrorControllerInvalidControllerId
	}
}

func (b *JsonBinder) BindList(c *gin.Context) (res interface{}, err error) {
	switch b.id {
	case ControllerIdProject:
		var docs []models.Project
		err = c.ShouldBindJSON(&docs)
		return docs, nil
	case ControllerIdUser:
		var docs []models.User
		err = c.ShouldBindJSON(&docs)
		return docs, nil
	default:
		return nil, errors.ErrorControllerInvalidControllerId
	}
}

func (b *JsonBinder) BindBatchRequestPayload(c *gin.Context) (payload entity.BatchRequestPayload, err error) {
	if err := c.ShouldBindJSON(&payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (b *JsonBinder) BindBatchRequestPayloadWithStringData(c *gin.Context) (payload entity.BatchRequestPayloadWithStringData, res models.BaseModelInterface, err error) {
	if err := c.ShouldBindJSON(&payload); err != nil {
		return payload, nil, err
	}
	doc, err := b.getDoc()
	if err != nil {
		return payload, nil, err
	}
	if err := json.Unmarshal([]byte(payload.Data), doc); err != nil {
		return payload, nil, err
	}
	return payload, doc, err
}

func (b *JsonBinder) getDoc() (doc models.BaseModelInterface, err error) {
	doc, ok := b.docMap[b.id]
	if !ok {
		return nil, errors.ErrorControllerInvalidControllerId
	}
	return doc, nil
}

func (b *JsonBinder) getDocList() (docs []models.BaseModelInterface, err error) {
	list, ok := b.docListMap[b.id]
	docs, err = models.GetBaseModelInterfaceList(list)
	if !ok {
		return nil, errors.ErrorControllerInvalidControllerId
	}
	return docs, nil
}
