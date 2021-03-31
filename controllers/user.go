package controllers

import (
	"github.com/crawlab-team/crawlab-core/models"
)

var UserController = NewListControllerDelegate(ControllerIdUser, models.UserService)
