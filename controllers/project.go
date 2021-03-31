package controllers

import (
	"github.com/crawlab-team/crawlab-core/models"
)

var ProjectController = NewListControllerDelegate(ControllerIdProject, models.ProjectService)
