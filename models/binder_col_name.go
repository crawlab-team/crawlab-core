package models

import "github.com/crawlab-team/crawlab-core/errors"

func NewColNameBinder(id ModelId) (b *ColNameBinder) {
	return &ColNameBinder{id: id}
}

type ColNameBinder struct {
	id ModelId
}

func (b *ColNameBinder) Bind() (res interface{}, err error) {
	switch b.id {
	case ModelIdNode:
		return ModelColNameNode, nil
	case ModelIdProject:
		return ModelColNameProject, nil
	case ModelIdSpider:
		return ModelColNameSpider, nil
	case ModelIdTask:
		return ModelColNameTask, nil
	case ModelIdJob:
		return ModelColNameTask, nil
	case ModelIdSchedule:
		return ModelColNameSchedule, nil
	case ModelIdUser:
		return ModelColNameUser, nil
	case ModelIdSetting:
		return ModelColNameSetting, nil
	case ModelIdToken:
		return ModelColNameToken, nil
	case ModelIdVariable:
		return ModelColNameVariable, nil
	default:
		panic(errors.ErrorModelNotImplemented)
	}
}

func (b *ColNameBinder) MustBind() (res interface{}) {
	res, _ = b.Bind()
	return res
}

func (b *ColNameBinder) BindString() (res string, err error) {
	res_, err := b.Bind()
	if err != nil {
		return "", err
	}
	res = res_.(string)
	return res, nil
}

func (b *ColNameBinder) MustBindString() (res string) {
	return b.MustBind().(string)
}
