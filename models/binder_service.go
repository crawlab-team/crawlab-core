package models

import "github.com/crawlab-team/crawlab-core/errors"

func NewServiceBinder(id ModelId) (b *ServiceBinder) {
	return &ServiceBinder{id: id}
}

type ServiceBinder struct {
	id ModelId
}

func (b *ServiceBinder) Bind() (res PublicServiceInterface, err error) {
	switch b.id {
	case ModelIdNode:
		return NodeService, nil
	case ModelIdProject:
		return ProjectService, nil
	case ModelIdSpider:
		return SpiderService, nil
	case ModelIdTask:
		return TaskService, nil
	case ModelIdJob:
		return TaskService, nil
	case ModelIdSchedule:
		return ScheduleService, nil
	case ModelIdUser:
		return UserService, nil
	case ModelIdSetting:
		return SettingService, nil
	case ModelIdToken:
		return TokenService, nil
	case ModelIdVariable:
		return VariableService, nil

	// invalid
	default:
		return res, errors.ErrorModelNotImplemented
	}
}

func (b *ServiceBinder) MustBind() (res PublicServiceInterface) {
	res, err := b.Bind()
	if err != nil {
		panic(err)
	}
	return res
}
