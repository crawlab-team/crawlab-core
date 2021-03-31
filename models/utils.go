package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"reflect"
)

func GetBaseModelInterfaceList(list interface{}) (res []BaseModelInterface, err error) {
	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(list)
		for i := 0; i < s.Len(); i++ {
			item := s.Index(i)
			obj, err := ConvertToBaseModelInterface(item.Interface())
			if err != nil {
				return nil, err
			}
			res = append(res, obj)
		}
	default:
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func ConvertToBaseModelInterface(item interface{}) (res BaseModelInterface, err error) {
	v := reflect.ValueOf(item)
	kind := v.Kind()
	switch kind {
	case reflect.Ptr:
		obj, ok := v.Interface().(BaseModelInterface)
		if !ok {
			return nil, errors.ErrorModelInvalidType
		}
		return obj, nil
	case reflect.Struct:
		return nil, errors.ErrorModelNotImplemented
	default:
		return nil, errors.ErrorModelInvalidType
	}
}
