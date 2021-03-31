package controllers

import (
	"github.com/crawlab-team/crawlab-core/models"
)

func NewListPostActionControllerDelegate(id ControllerId, svc models.PublicServiceInterface, actions []PostAction) (d *ListPostActionControllerDelegate) {
	return &ListPostActionControllerDelegate{
		NewListControllerDelegate(id, svc),
		NewPostActionControllerDelegate(id, actions),
	}
}

type ListPostActionControllerDelegate struct {
	ListController
	PostActionController
}
