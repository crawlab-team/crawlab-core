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
	fr := b.fr

	switch b.id {
	case ModelIdNode:
		err = fr.One(&m.Node)
		return m.Node, err
	case ModelIdProject:
		err = fr.One(&m.Project)
		return m.Project, err
	case ModelIdSpider:
		err = fr.One(&m.Spider)
		return m.Spider, err
	case ModelIdTask:
		err = fr.One(&m.Task)
		return m.Task, err
	case ModelIdJob:
		err = fr.One(&m.Job)
		return m.Job, err
	case ModelIdSchedule:
		err = fr.One(&m.Schedule)
		return m.Schedule, err
	case ModelIdUser:
		err = fr.One(&m.User)
		return m.User, err
	case ModelIdSetting:
		err = fr.One(&m.Setting)
		return m.Setting, err
	case ModelIdToken:
		err = fr.One(&m.Token)
		return m.Token, err
	case ModelIdVariable:
		err = fr.One(&m.Variable)
		return m.Variable, err
	case ModelIdTag:
		err = fr.One(&m.Tag)
		return m.Tag, err
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
