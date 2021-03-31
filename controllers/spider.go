package controllers

import "github.com/crawlab-team/crawlab-core/models"

var SpiderController = NewListPostActionControllerDelegate(ControllerIdSpider, models.SpiderService, []PostAction{})
