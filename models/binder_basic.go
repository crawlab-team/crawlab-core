package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-db/mongo"
)

func NewBasicBinder(id ModelId, m *ModelMap, fr *mongo.FindResult) (b *BasicBinder) {
	return &BasicBinder{
		id: id,
		m:  m,
		fr: fr,
	}
}

type BasicBinder struct {
	id ModelId
	m  *ModelMap
	fr *mongo.FindResult
}

func (b *BasicBinder) Bind() (res interface{}, err error) {
	m := b.m

	switch b.id {
	case ModelIdNode:
		return b.process(&m.Node, ModelIdTag)
	case ModelIdProject:
		return b.process(&m.Project, ModelIdTag)
	case ModelIdSpider:
		return b.process(&m.Spider, ModelIdTag)
	case ModelIdTask:
		return b.process(&m.Task)
	case ModelIdJob:
		return b.process(&m.Job)
	case ModelIdSchedule:
		return b.process(&m.Schedule)
	case ModelIdUser:
		return b.process(&m.User)
	case ModelIdSetting:
		return b.process(&m.Setting)
	case ModelIdToken:
		return b.process(&m.Token)
	case ModelIdVariable:
		return b.process(&m.Variable)
	case ModelIdTag:
		return b.process(&m.Tag)
	default:
		return nil, errors.ErrorModelInvalidModelId
	}
}

func (b *BasicBinder) MustBind() (res interface{}) {
	res, err := b.Bind()
	if err != nil {
		panic(err)
	}
	return res
}

func (b *BasicBinder) process(d interface{}, fieldIds ...ModelId) (res interface{}, err error) {
	if err := b.fr.One(d); err != nil {
		return nil, err
	}
	return assignFields(d, fieldIds...)
}
