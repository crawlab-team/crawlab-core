package models

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/emirpasic/gods/lists/arraylist"
	"reflect"
)

func assignFields(d interface{}, fieldIds ...ModelId) (res interface{}, err error) {
	doc, ok := d.(BaseModelInterface)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	if len(fieldIds) == 0 {
		return doc, nil
	}
	a, err := doc.GetArtifact()
	if err != nil {
		return nil, err
	}
	for _, fid := range fieldIds {
		switch fid {
		case ModelIdTag:
			d, ok := doc.(BaseModelWithTagsInterface)
			if !ok {
				return nil, errors.ErrorModelInvalidType
			}
			tags, err := a.GetTags()
			if err != nil {
				return nil, err
			}
			d.SetTags(tags)
			return d, nil
		}
	}
	return doc, nil
}

func _assignListFields(asPtr bool, list interface{}, fieldIds ...ModelId) (res arraylist.List, err error) {
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
		doc, ok := item.(BaseModelInterface)
		if !ok {
			return res, errors.ErrorModelInvalidType
		}
		ptr, err := assignFields(doc, fieldIds...)
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

func assignListFields(list interface{}, fieldIds ...ModelId) (res arraylist.List, err error) {
	return _assignListFields(false, list, fieldIds...)
}

func assignListFieldsAsPtr(list interface{}, fieldIds ...ModelId) (res arraylist.List, err error) {
	return _assignListFields(true, list, fieldIds...)
}

func serializeList(list arraylist.List, target interface{}) (err error) {
	bytes, err := list.ToJSON()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, target); err != nil {
		return err
	}
	return nil
}
