package controllers

import (
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/schedule"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func postScheduleEnableDisableFunc(isEnable bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
		svc, err := schedule.GetScheduleServiceV2()
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
		s, err := service.NewModelServiceV2[models.ScheduleV2]().GetById(id)
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
		u := GetUserFromContextV2(c)
		if isEnable {
			err = svc.Enable(*s, u.Id)
		} else {
			err = svc.Disable(*s, u.Id)
		}
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
		HandleSuccess(c)
	}
}

func PostScheduleEnable(c *gin.Context) {
	postScheduleEnableDisableFunc(true)(c)
}

func PostScheduleDisable(c *gin.Context) {
	postScheduleEnableDisableFunc(false)(c)
}
