package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
)

func NewServiceBinder(id interfaces.ModelId) (b *ServiceBinder) {
	return &ServiceBinder{id: id}
}

type ServiceBinder struct {
	id interfaces.ModelId
}

func (b *ServiceBinder) Bind() (res interfaces.ModelService, err error) {
	switch b.id {
	case interfaces.ModelIdNode:
		return NodeService, nil
	case interfaces.ModelIdProject:
		return ProjectService, nil
	case interfaces.ModelIdSpider:
		return SpiderService, nil
	case interfaces.ModelIdTask:
		return TaskService, nil
	case interfaces.ModelIdJob:
		return TaskService, nil
	case interfaces.ModelIdSchedule:
		return ScheduleService, nil
	case interfaces.ModelIdUser:
		return UserService, nil
	case interfaces.ModelIdSetting:
		return SettingService, nil
	case interfaces.ModelIdToken:
		return TokenService, nil
	case interfaces.ModelIdVariable:
		return VariableService, nil

	// invalid
	default:
		return res, errors.ErrorModelNotImplemented
	}
}

func (b *ServiceBinder) MustBind() (res interfaces.ModelService) {
	res, err := b.Bind()
	if err != nil {
		panic(err)
	}
	return res
}
