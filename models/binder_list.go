package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/emirpasic/gods/lists/arraylist"
)

func NewListBinder(id ModelId, m *ModelListMap, fr *mongo.FindResult) (b *ListBinder) {
	return &ListBinder{
		id: id,
		m:  m,
		fr: fr,
	}
}

type ListBinder struct {
	id    ModelId
	m     *ModelListMap
	fr    *mongo.FindResult
	asPtr bool
}

func (b *ListBinder) Bind() (res interface{}, err error) {
	m := b.m

	switch b.id {
	case ModelIdNode:
		return b.process(m.Nodes, ModelIdTag)
	case ModelIdProject:
		return b.process(m.Projects, ModelIdTag)
	case ModelIdSpider:
		return b.process(m.Spiders, ModelIdTag)
	case ModelIdTask:
		return b.process(m.Tasks)
	case ModelIdSchedule:
		return b.process(m.Schedules)
	case ModelIdUser:
		return b.process(m.Users)
	case ModelIdSetting:
		return b.process(m.Settings)
	case ModelIdToken:
		return b.process(m.Tokens)
	case ModelIdVariable:
		return b.process(m.Variables)
	case ModelIdTag:
		return b.process(m.Tags)
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

func (b *ListBinder) BindList() (res arraylist.List, err error) {
	r, err := b.Bind()
	if err != nil {
		return res, err
	}
	res, ok := r.(arraylist.List)
	if !ok {
		return res, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (b *ListBinder) MustBindList() (res arraylist.List) {
	res, err := b.BindList()
	if err != nil {
		panic(err)
	}
	return res
}

func (b *ListBinder) MustBindListAsPtr() (res arraylist.List) {
	res, err := b.BindListAsPtr()
	if err != nil {
		panic(err)
	}
	return res
}

func (b *ListBinder) BindListAsPtr() (res arraylist.List, err error) {
	b.asPtr = true

	r, err := b.Bind()
	if err != nil {
		return res, err
	}
	res, ok := r.(arraylist.List)
	if !ok {
		return res, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (b *ListBinder) process(d interface{}, fieldIds ...ModelId) (res interface{}, err error) {
	if err := b.fr.All(&d); err != nil {
		return nil, err
	}
	if b.asPtr {
		return assignListFieldsAsPtr(d, fieldIds...)
	} else {
		return assignListFields(d, fieldIds...)
	}
}
