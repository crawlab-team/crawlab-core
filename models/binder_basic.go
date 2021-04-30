package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
)

func NewBasicBinder(id interfaces.ModelId, m *ModelMap, fr *mongo.FindResult) (b *BasicBinder) {
	return &BasicBinder{
		id: id,
		m:  m,
		fr: fr,
	}
}

type BasicBinder struct {
	id interfaces.ModelId
	m  *ModelMap
	fr *mongo.FindResult
}

func (b *BasicBinder) Bind() (res interface{}, err error) {
	m := b.m

	switch b.id {
	case interfaces.ModelIdNode:
		return b.process(&m.Node, interfaces.ModelIdTag)
	case interfaces.ModelIdProject:
		return b.process(&m.Project, interfaces.ModelIdTag)
	case interfaces.ModelIdSpider:
		return b.process(&m.Spider, interfaces.ModelIdTag)
	case interfaces.ModelIdTask:
		return b.process(&m.Task)
	case interfaces.ModelIdJob:
		return b.process(&m.Job)
	case interfaces.ModelIdSchedule:
		return b.process(&m.Schedule)
	case interfaces.ModelIdUser:
		return b.process(&m.User)
	case interfaces.ModelIdSetting:
		return b.process(&m.Setting)
	case interfaces.ModelIdToken:
		return b.process(&m.Token)
	case interfaces.ModelIdVariable:
		return b.process(&m.Variable)
	case interfaces.ModelIdTag:
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

func (b *BasicBinder) process(d interface{}, fieldIds ...interfaces.ModelId) (res interface{}, err error) {
	if err := b.fr.One(d); err != nil {
		return nil, err
	}
	return AssignFields(d, fieldIds...)
}
