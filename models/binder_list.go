package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-db/mongo"
)

func NewListBinder(id ModelId, m *ModelListMap, fr *mongo.FindResult) (b *ListBinder) {
	return &ListBinder{
		id: id,
		m:  m,
		fr: fr,
	}
}

type ListBinder struct {
	id ModelId
	m  *ModelListMap
	fr *mongo.FindResult
}

func (b *ListBinder) Bind() (res interface{}, err error) {
	m := b.m
	fr := b.fr

	switch b.id {
	case ModelIdNode:
		err = fr.All(&m.Nodes)
		return m.Nodes, err
	case ModelIdProject:
		err = fr.All(&m.Projects)
		return m.Projects, err
	case ModelIdSpider:
		err = fr.All(&m.Spiders)
		return m.Spiders, err
	case ModelIdTask:
		err = fr.All(&m.Tasks)
		return m.Tasks, err
	case ModelIdSchedule:
		err = fr.All(&m.Schedules)
		return m.Schedules, err
	case ModelIdUser:
		err = fr.All(&m.Users)
		return m.Users, err
	case ModelIdSetting:
		err = fr.All(&m.Settings)
		return m.Settings, err
	case ModelIdToken:
		err = fr.All(&m.Tokens)
		return m.Tokens, err
	case ModelIdVariable:
		err = fr.All(&m.Variables)
		return m.Variables, err
	default:
		return nil, errors.ErrorModelInvalidModelId
	}
}

func (b *ListBinder) MustBind() (res interface{}) {
	res, err := b.Bind()
	if err != nil {
		panic(err)
	}
	return res
}
