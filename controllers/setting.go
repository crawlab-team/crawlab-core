package controllers

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/gin-gonic/gin"
)

var SettingController *settingController

type settingController struct {
	ListControllerDelegate
}

func (ctr *settingController) Get(c *gin.Context) {
	// key
	key := c.Param("id")

	// model service
	modelSvc, err := service.NewService()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// setting
	s, err := modelSvc.GetSettingByKey(key, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, s)
}

func (ctr *settingController) Post(c *gin.Context) {
	// key
	key := c.Param("id")

	// value
	var value string
	if err := c.ShouldBindJSON(&value); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// model service
	modelSvc, err := service.NewService()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// setting
	s, err := modelSvc.GetSettingByKey(key, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// save
	s.Value = value
	if err := delegate.NewModelDelegate(s).Save(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func newSettingController() *settingController {
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListControllerDelegate(ControllerIdSetting, modelSvc.GetBaseService(interfaces.ModelIdSetting))

	return &settingController{
		ListControllerDelegate: *ctr,
	}
}
