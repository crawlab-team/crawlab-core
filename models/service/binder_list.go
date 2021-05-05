package service

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/emirpasic/gods/lists/arraylist"
	"reflect"
)

func NewListBinder(id interfaces.ModelId, fr *mongo.FindResult) (b interfaces.ModelListBinder) {
	return &ListBinder{
		id: id,
		m:  models.NewModelListMap(),
		fr: fr,
		b:  NewBasicBinder(id, fr),
	}
}

type ListBinder struct {
	id    interfaces.ModelId
	m     *models.ModelListMap
	fr    *mongo.FindResult
	asPtr bool
	b     interfaces.ModelBinder
}

func (b *ListBinder) Bind() (res interface{}, err error) {
	m := b.m

	switch b.id {
	case interfaces.ModelIdArtifact:
		return b.Process(m.Artifacts)
	case interfaces.ModelIdTag:
		return b.Process(m.Tags)
	case interfaces.ModelIdNode:
		return b.Process(m.Nodes, interfaces.ModelIdTag)
	case interfaces.ModelIdProject:
		return b.Process(m.Projects, interfaces.ModelIdTag)
	case interfaces.ModelIdSpider:
		return b.Process(m.Spiders, interfaces.ModelIdTag)
	case interfaces.ModelIdTask:
		return b.Process(m.Tasks)
	case interfaces.ModelIdSchedule:
		return b.Process(m.Schedules)
	case interfaces.ModelIdUser:
		return b.Process(m.Users)
	case interfaces.ModelIdSetting:
		return b.Process(m.Settings)
	case interfaces.ModelIdToken:
		return b.Process(m.Tokens)
	case interfaces.ModelIdVariable:
		return b.Process(m.Variables)
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

func (b *ListBinder) Process(d interface{}, fieldIds ...interfaces.ModelId) (res interface{}, err error) {
	if err := b.fr.All(&d); err != nil {
		return nil, err
	}
	// TODO: implement in delegate
	if b.asPtr {
		return b.AssignListFieldsAsPtr(d, fieldIds...)
	} else {
		return b.AssignListFields(d, fieldIds...)
	}
}

func (b *ListBinder) AssignListFields(list interface{}, fieldIds ...interfaces.ModelId) (res arraylist.List, err error) {
	return b.assignListFields(list, fieldIds...)
}

func (b *ListBinder) AssignListFieldsAsPtr(list interface{}, fieldIds ...interfaces.ModelId) (res arraylist.List, err error) {
	return b.assignListFieldsAsPtr(list, fieldIds...)
}

func (b *ListBinder) _assignListFields(asPtr bool, list interface{}, fieldIds ...interfaces.ModelId) (res arraylist.List, err error) {
	vList := reflect.ValueOf(list)
	if vList.Kind() != reflect.Array &&
		vList.Kind() != reflect.Slice {
		return res, errors.ErrorModelInvalidType
	}
	for i := 0; i < vList.Len(); i++ {
		vItem := vList.Index(i)
		var item interface{}
		if vItem.CanAddr() {
			item = vItem.Addr().Interface()
		} else {
			item = vItem.Interface()
		}
		doc, ok := item.(interfaces.Model)
		if !ok {
			return res, errors.ErrorModelInvalidType
		}
		ptr, err := b.b.AssignFields(doc, fieldIds...)
		if err != nil {
			return res, err
		}
		v := reflect.ValueOf(ptr)
		if !asPtr {
			// non-pointer item
			res.Add(v.Elem().Interface())
		} else {
			// pointer item
			res.Add(v.Interface())
		}
	}
	return res, nil
}

func (b *ListBinder) assignListFields(list interface{}, fieldIds ...interfaces.ModelId) (res arraylist.List, err error) {
	return b._assignListFields(false, list, fieldIds...)
}

func (b *ListBinder) assignListFieldsAsPtr(list interface{}, fieldIds ...interfaces.ModelId) (res arraylist.List, err error) {
	return b._assignListFields(true, list, fieldIds...)
}
