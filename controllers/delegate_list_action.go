package controllers

import (
	"github.com/crawlab-team/crawlab-core/models"
)

func NewListPostActionControllerDelegate(id ControllerId, svc models.PublicServiceInterface, actions []Action) (d *ListActionControllerDelegate) {
	return &ListActionControllerDelegate{
		NewListControllerDelegate(id, svc),
		NewActionControllerDelegate(id, actions),
	}
}

type ListActionControllerDelegate struct {
	ListController
	ActionController
}
