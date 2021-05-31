package service

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"github.com/emirpasic/gods/lists/arraylist"
	"reflect"
)

func NewListBinder(id interfaces.ModelId, fr *mongo.FindResult) (b interfaces.ModelListBinder) {
	return &ListBinder{
		id:    id,
		m:     models.NewModelListMap(),
		fr:    fr,
		asPtr: true,
		wf:    true,
		b:     NewBasicBinder(id, fr),
	}
}

type ListBinder struct {
	id    interfaces.ModelId
	m     *models.ModelListMap
	fr    *mongo.FindResult
	asPtr bool // whether to process to pointer item list
	wf    bool // whether to process with fields
	b     interfaces.ModelBinder
}

func (b *ListBinder) Bind() (list arraylist.List, err error) {
	m := b.m

	switch b.id {
	case interfaces.ModelIdArtifact:
		return b.Process(m.Artifacts)
	case interfaces.ModelIdTag:
		return b.Process(m.Tags)
	case interfaces.ModelIdNode:
		return b.ProcessWithFieldIds(m.Nodes, interfaces.ModelIdTag)
	case interfaces.ModelIdProject:
		return b.ProcessWithFieldIds(m.Projects, interfaces.ModelIdTag)
	case interfaces.ModelIdSpider:
		return b.ProcessWithFieldIds(m.Spiders, interfaces.ModelIdTag)
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
	case interfaces.ModelIdTaskQueue:
		return b.Process(m.TaskQueueItems)
	case interfaces.ModelIdTaskStat:
		return b.Process(m.TaskStats)
	default:
		return list, errors.ErrorModelInvalidModelId
	}
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
	return b.Bind()
}

func (b *ListBinder) MustBindListWithNoFields() (res arraylist.List) {
	res, err := b.BindListWithNoFields()
	if err != nil {
		panic(err)
	}
	return res
}

func (b *ListBinder) BindListWithNoFields() (res arraylist.List, err error) {
	b.wf = false
	return b.Bind()
}

func (b *ListBinder) Process(d interface{}) (list arraylist.List, err error) {
	if err := b.fr.All(&d); err != nil {
		return list, trace.TraceError(err)
	}
	if b.asPtr {
		return b.AssignListFieldsAsPtr(d)
	} else {
		return b.AssignListFields(d)
	}
}

func (b *ListBinder) ProcessWithFieldIds(d interface{}, fieldIds ...interfaces.ModelId) (list arraylist.List, err error) {
	if !b.wf {
		return b.Process(d)
	}
	if err := b.fr.All(&d); err != nil {
		return list, trace.TraceError(err)
	}
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
