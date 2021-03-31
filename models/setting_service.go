package models

type settingService struct {
	*Service
}

func NewSettingService() (svc *settingService) {
	return &settingService{NewService(ModelIdSetting)}
}

var SettingService = NewSettingService()
